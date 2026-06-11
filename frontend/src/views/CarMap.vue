<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getCarPosition } from '@/api/client'
import { fetchRoute, type RouteData } from '@/api/route'
import type { CarPosition } from '@/api/types'
import MapView from '@/components/MapView.vue'
import MapControls from '@/components/MapControls.vue'
import ArrivalSheet from '@/components/ArrivalSheet.vue'
import StateView from '@/components/StateView.vue'
import { useGeolocation } from '@/composables/useGeolocation'

const STALE_MINUTES = 10
const RETRY_SECONDS = 8

const route = useRoute()
const data = ref<CarPosition | null>(null)
const error = ref<string | null>(null)
const loading = ref(true)
const routeData = ref<RouteData | null>(null)
const satellite = ref(false)
const retryIn = ref(RETRY_SECONDS)
const maxDistanceKm = ref(0)

const mapViewRef = ref<{ recenter: () => void } | null>(null)

const { coords: userCoords } = useGeolocation()

const carId = computed(() => Number(route.params.id))

async function fetchPosition() {
  try {
    data.value = await getCarPosition(carId.value)
    error.value = null
    retryIn.value = RETRY_SECONDS
  } catch {
    error.value = 'Failed to load car position'
    retryIn.value = retryIn.value > 1 ? retryIn.value - 1 : RETRY_SECONDS
  } finally {
    loading.value = false
  }
}

let positionInterval: ReturnType<typeof setInterval>
let routeInterval: ReturnType<typeof setInterval>

async function updateRoute() {
  if (!carCoords.value || !userCoords.value) return
  routeData.value = await fetchRoute(
    userCoords.value.lat,
    userCoords.value.lng,
    carCoords.value.lat,
    carCoords.value.lng,
  )
}

onMounted(() => {
  fetchPosition()
  positionInterval = setInterval(fetchPosition, 1000)
  // Refresh the route every 30 seconds (don't hammer OSRM on every position update)
  routeInterval = setInterval(updateRoute, 30_000)
})

onUnmounted(() => {
  clearInterval(positionInterval)
  clearInterval(routeInterval)
})

const carCoords = computed(() => {
  if (!data.value?.position) return null
  return {
    lat: data.value.position.latitude,
    lng: data.value.position.longitude,
  }
})

// Fetch route immediately when both positions become available
watch([carCoords, userCoords], () => {
  if (carCoords.value && userCoords.value && !routeData.value) {
    updateRoute()
  }
})

// Track the furthest distance seen so the approach bar can show progress
watch(
  () => routeData.value?.distance_km,
  (d) => {
    if (d != null && d > maxDistanceKm.value) maxDistanceKm.value = d
  },
)

const progress = computed(() => {
  const d = routeData.value?.distance_km
  if (d == null || maxDistanceKm.value <= 0) return 0
  return Math.max(0, Math.min(1, 1 - d / maxDistanceKm.value))
})

const displayName = computed(() => {
  if (!data.value) return ''
  const c = data.value.car
  return c.name || `${c.model} ${c.trim_badging || ''}`.trim()
})

const isStale = computed(() => {
  const d = data.value?.position?.date
  if (!d) return false
  return Date.now() - new Date(d).getTime() > STALE_MINUTES * 60_000
})

const lastSeen = computed(() => {
  const d = data.value?.position?.date
  if (!d) return null
  const mins = Math.floor((Date.now() - new Date(d).getTime()) / 60_000)
  if (mins < 60) return `${mins}m ago`
  const hrs = Math.floor(mins / 60)
  if (hrs < 24) return `${hrs}h ago`
  return `${Math.floor(hrs / 24)}d ago`
})

// which overlay (if any) sits on top of the map
type Phase = 'skeleton' | 'error' | 'stale' | 'live' | 'no-position'
const phase = computed<Phase>(() => {
  if (loading.value && !data.value) return 'skeleton'
  if (error.value && !data.value) return 'error'
  if (!data.value?.position) return 'no-position'
  if (isStale.value) return 'stale'
  return 'live'
})

const showMap = computed(() => !!carCoords.value)
</script>

<template>
  <div class="page">
    <!-- map fills the screen; states float on top -->
    <div v-if="showMap" class="map" :class="{ muted: phase === 'stale' }">
      <MapView
        ref="mapViewRef"
        :car-coords="carCoords!"
        :user-coords="userCoords"
        :car="data!.car"
        :heading="data!.position?.heading ?? null"
        :route-coords="routeData?.routeCoords ?? null"
        :satellite="satellite"
      />
    </div>

    <MapControls v-if="showMap" v-model:satellite="satellite" />

    <!-- LIVE: recenter + arrival sheet docked at the bottom -->
    <div v-if="phase === 'live'" class="dock">
      <div class="dock-row">
        <button class="recenter" aria-label="Recenter" @click="mapViewRef?.recenter()">
          <span class="ring"></span>
        </button>
      </div>
      <ArrivalSheet
        :car-name="displayName"
        :eta-minutes="routeData?.duration_min ?? null"
        :distance-km="routeData?.distance_km ?? null"
        :battery-level="data!.position?.battery_level ?? null"
        :speed="data!.position?.speed ?? null"
        :updated-date="data!.position?.date ?? null"
        :progress="progress"
      />
    </div>

    <!-- states -->
    <StateView v-if="phase === 'skeleton'" variant="skeleton" />
    <StateView
      v-else-if="phase === 'error'"
      variant="error"
      :message="error"
      :retry-seconds="retryIn"
      @retry="fetchPosition"
    />
    <StateView
      v-else-if="phase === 'stale' || phase === 'no-position'"
      variant="stale"
      :last-seen="lastSeen"
      :battery-level="data?.position?.battery_level ?? null"
    />
  </div>
</template>

<style scoped>
.page {
  position: relative;
  height: 100vh;
  height: 100dvh;
  overflow: hidden;
  background: #0e1217;
}

.map {
  position: absolute;
  inset: 0;
}
.map.muted {
  filter: grayscale(0.5) brightness(0.7);
}

/* bottom dock keeps the recenter button floating above the sheet */
.dock {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
  pointer-events: none;
}
.dock > * {
  pointer-events: auto;
}
.dock-row {
  display: flex;
  justify-content: flex-end;
  padding: 0 18px 14px;
}

.recenter {
  width: 44px;
  height: 44px;
  border-radius: var(--tm-r-md);
  background: rgba(16, 20, 26, 0.78);
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}
.recenter .ring {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 2px solid var(--tm-accent);
  position: relative;
}
.recenter .ring::after {
  content: '';
  position: absolute;
  inset: 4px;
  border-radius: 50%;
  background: var(--tm-accent);
}

/* desktop: float the sheet bottom-left instead of full-width */
@media (min-width: 900px) {
  .dock {
    left: 28px;
    right: auto;
    bottom: 28px;
    width: 380px;
  }
  .dock :deep(.sheet) {
    border: 1px solid var(--tm-line-strong);
    border-radius: var(--tm-r-lg);
  }
}
</style>
