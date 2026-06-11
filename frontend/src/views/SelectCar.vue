<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { listCars, getCarPosition } from '@/api/client'
import type { Car, Position } from '@/api/types'
import CarCard from '@/components/CarCard.vue'

const cars = ref<Car[]>([])
const positions = ref<Record<number, Position | null>>({})
const error = ref<string | null>(null)
const loading = ref(true)

const count = computed(() => cars.value.length)

onMounted(async () => {
  try {
    cars.value = await listCars()
    // hydrate each card's battery/status in the background; tolerate failures
    await Promise.all(
      cars.value.map(async (c) => {
        try {
          positions.value[c.id] = (await getCarPosition(c.id)).position
        } catch {
          positions.value[c.id] = null
        }
      }),
    )
  } catch {
    error.value = 'Failed to load cars'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="page">
    <div class="shell">
      <header class="head">
        <div class="tm-label">Your garage</div>
        <div class="head-row">
          <h1>{{ loading ? 'Loading…' : `${count} ${count === 1 ? 'vehicle' : 'vehicles'}` }}</h1>
        </div>
      </header>

      <!-- skeleton -->
      <div v-if="loading" class="list">
        <div v-for="n in 3" :key="n" class="skel" :class="{ tall: n === 1 }"></div>
      </div>

      <p v-else-if="error" class="status err">{{ error }}</p>
      <p v-else-if="!cars.length" class="status">No vehicles on this account.</p>

      <div v-else class="list">
        <CarCard
          v-for="(car, i) in cars"
          :key="car.id"
          :car="car"
          :position="positions[car.id] ?? null"
          :featured="i === 0"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.page {
  min-height: 100vh;
  background:
    radial-gradient(900px 500px at 20% -8%, rgba(77, 139, 255, 0.1), transparent 60%),
    var(--tm-bg);
  padding: 32px 16px calc(48px + env(safe-area-inset-bottom));
}
.shell {
  max-width: 560px;
  margin: 0 auto;
}

.head {
  padding: 8px 6px 20px;
}
.head-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 6px;
}
h1 {
  font-family: var(--tm-font-display);
  font-weight: 600;
  font-size: 30px;
  letter-spacing: -0.01em;
}

.list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.status {
  color: var(--tm-text-mid);
  padding: 24px 6px;
}
.status.err {
  color: var(--tm-low);
}

/* loading skeletons */
.skel {
  height: 96px;
  border-radius: 22px;
  background: linear-gradient(90deg, #14181f 25%, #1d2330 50%, #14181f 75%);
  background-size: 520px 100%;
  animation: tmShimmer 1.3s linear infinite;
}
.skel.tall {
  height: 150px;
}
</style>
