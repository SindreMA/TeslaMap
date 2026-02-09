<script setup lang="ts">
import { computed } from 'vue'
import type { Position, Car } from '@/api/types'

const props = defineProps<{
  position: Position
  userCoords: { lat: number; lng: number } | null
  car: Car
}>()

function haversineKm(
  lat1: number, lng1: number,
  lat2: number, lng2: number
): number {
  const R = 6371
  const dLat = toRad(lat2 - lat1)
  const dLng = toRad(lng2 - lng1)
  const a =
    Math.sin(dLat / 2) ** 2 +
    Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.sin(dLng / 2) ** 2
  return R * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))
}

function toRad(deg: number): number {
  return (deg * Math.PI) / 180
}

const distance = computed(() => {
  if (!props.userCoords) return null
  return haversineKm(
    props.position.latitude,
    props.position.longitude,
    props.userCoords.lat,
    props.userCoords.lng
  )
})

const etaMinutes = computed(() => {
  if (distance.value === null) return null
  // Rough ETA at 60 km/h average driving speed
  return Math.round((distance.value / 60) * 60)
})

const timeAgo = computed(() => {
  const now = Date.now()
  const then = new Date(props.position.date).getTime()
  const diffMs = now - then
  const mins = Math.floor(diffMs / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hrs = Math.floor(mins / 60)
  if (hrs < 24) return `${hrs}h ${mins % 60}m ago`
  const days = Math.floor(hrs / 24)
  return `${days}d ago`
})
</script>

<template>
  <div class="panel">
    <div class="row">
      <span class="label">Last updated</span>
      <span class="value">{{ timeAgo }}</span>
    </div>
    <div class="row" v-if="position.speed !== null">
      <span class="label">Speed</span>
      <span class="value">{{ position.speed }} km/h</span>
    </div>
    <div class="row" v-if="position.battery_level !== null">
      <span class="label">Battery</span>
      <span class="value">{{ position.battery_level }}%</span>
    </div>
    <div class="row" v-if="distance !== null">
      <span class="label">Distance</span>
      <span class="value">{{ distance.toFixed(1) }} km</span>
    </div>
    <div class="row" v-if="etaMinutes !== null">
      <span class="label">ETA (approx)</span>
      <span class="value">{{ etaMinutes }} min</span>
    </div>
  </div>
</template>

<style scoped>
.panel {
  background: #1a1a1a;
  border-top: 1px solid #2a2a2a;
  padding: 1rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem 1.5rem;
}

.row {
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
}

.label {
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #666;
}

.value {
  font-size: 1rem;
  font-weight: 600;
}
</style>
