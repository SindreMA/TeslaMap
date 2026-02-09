<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getCarPosition } from '@/api/client'
import { fetchRoute, type RouteData } from '@/api/route'
import type { CarPosition } from '@/api/types'
import MapView from '@/components/MapView.vue'
import DistanceInfo from '@/components/DistanceInfo.vue'
import { useGeolocation } from '@/composables/useGeolocation'

const route = useRoute()
const data = ref<CarPosition | null>(null)
const error = ref<string | null>(null)
const loading = ref(true)
const routeData = ref<RouteData | null>(null)

const { coords: userCoords } = useGeolocation()

const carId = computed(() => Number(route.params.id))

async function fetchPosition() {
  try {
    data.value = await getCarPosition(carId.value)
    error.value = null
  } catch {
    error.value = 'Failed to load car position'
  } finally {
    loading.value = false
  }
}

let positionInterval: ReturnType<typeof setInterval>
let routeInterval: ReturnType<typeof setInterval>

async function updateRoute() {
  if (!carCoords.value || !userCoords.value) return
  routeData.value = await fetchRoute(
    userCoords.value.lat, userCoords.value.lng,
    carCoords.value.lat, carCoords.value.lng,
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

const displayName = computed(() => {
  if (!data.value) return ''
  const c = data.value.car
  return c.name || `${c.model} ${c.trim_badging || ''}`.trim()
})
</script>

<template>
  <div class="page">
    <header class="header">
      <router-link to="/select" class="back">&larr; Cars</router-link>
      <h1 v-if="data">{{ displayName }}</h1>
    </header>

    <p v-if="loading" class="status">Loading...</p>
    <p v-else-if="error" class="status error">{{ error }}</p>
    <p v-else-if="!data?.position" class="status">No position data available</p>

    <template v-else>
      <div class="content">
        <MapView
          :car-coords="carCoords!"
          :user-coords="userCoords"
          :car="data.car"
          :heading="data.position?.heading ?? null"
          :route-coords="routeData?.routeCoords ?? null"
        />
        <DistanceInfo
          :position="data.position!"
          :user-coords="userCoords"
          :car="data.car"
          :route-distance="routeData?.distance_km ?? null"
          :route-duration="routeData?.duration_min ?? null"
        />
      </div>
    </template>
  </div>
</template>

<style scoped>
.page {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem 1rem;
  background: #1a1a1a;
  border-bottom: 1px solid #2a2a2a;
}

.back {
  font-size: 0.9rem;
  color: #3b82f6;
}

h1 {
  font-size: 1.1rem;
  font-weight: 600;
}

.status {
  padding: 2rem;
  text-align: center;
  color: #888;
}

.error {
  color: #ef4444;
}

.content {
  flex: 1;
  display: flex;
  flex-direction: column;
  position: relative;
}
</style>
