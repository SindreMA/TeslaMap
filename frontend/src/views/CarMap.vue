<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { getCarPosition } from '@/api/client'
import type { CarPosition } from '@/api/types'
import MapView from '@/components/MapView.vue'
import DistanceInfo from '@/components/DistanceInfo.vue'
import { useGeolocation } from '@/composables/useGeolocation'

const route = useRoute()
const data = ref<CarPosition | null>(null)
const error = ref<string | null>(null)
const loading = ref(true)

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

let interval: ReturnType<typeof setInterval>

onMounted(() => {
  fetchPosition()
  interval = setInterval(fetchPosition, 1000)
})

onUnmounted(() => {
  clearInterval(interval)
})

const carCoords = computed(() => {
  if (!data.value?.position) return null
  return {
    lat: data.value.position.latitude,
    lng: data.value.position.longitude,
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
        />
        <DistanceInfo
          :position="data.position!"
          :user-coords="userCoords"
          :car="data.car"
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
