<script setup lang="ts">
import type { Car } from '@/api/types'

defineProps<{ car: Car }>()

function displayName(car: Car): string {
  return car.name || `${car.model} ${car.trim_badging || ''}`.trim()
}
</script>

<template>
  <router-link :to="`/car/${car.id}`" class="card">
    <div class="name">{{ displayName(car) }}</div>
    <div class="details">
      <span v-if="car.exterior_color" class="tag">{{ car.exterior_color }}</span>
      <span class="tag vin">{{ car.vin.slice(-6) }}</span>
    </div>
  </router-link>
</template>

<style scoped>
.card {
  display: block;
  padding: 1.25rem;
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  border-radius: 12px;
  transition: border-color 0.2s, background 0.2s;
  color: #e0e0e0;
}

.card:hover {
  border-color: #3b82f6;
  background: #1e1e2e;
}

.name {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.details {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.tag {
  font-size: 0.75rem;
  padding: 0.2rem 0.5rem;
  background: #2a2a2a;
  border-radius: 4px;
  color: #aaa;
}

.vin {
  font-family: monospace;
}
</style>
