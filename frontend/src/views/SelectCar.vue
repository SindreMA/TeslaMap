<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { listCars } from '@/api/client'
import type { Car } from '@/api/types'
import CarCard from '@/components/CarCard.vue'

const cars = ref<Car[]>([])
const error = ref<string | null>(null)
const loading = ref(true)

onMounted(async () => {
  try {
    cars.value = await listCars()
  } catch (e) {
    error.value = 'Failed to load cars'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="page">
    <h1>Select a car</h1>
    <p v-if="loading" class="status">Loading...</p>
    <p v-else-if="error" class="status error">{{ error }}</p>
    <div v-else class="grid">
      <CarCard v-for="car in cars" :key="car.id" :car="car" />
    </div>
  </div>
</template>

<style scoped>
.page {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

h1 {
  font-size: 1.5rem;
  margin-bottom: 1.5rem;
}

.status {
  color: #888;
}

.error {
  color: #ef4444;
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 1rem;
}
</style>
