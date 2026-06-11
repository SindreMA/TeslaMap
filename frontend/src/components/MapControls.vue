<script setup lang="ts">
defineProps<{
  /** true = satellite tiles, false = dark map */
  satellite: boolean
}>()

const emit = defineEmits<{
  (e: 'update:satellite', value: boolean): void
}>()
</script>

<template>
  <div class="chrome">
    <router-link to="/select" class="back glass">
      <span class="chev">‹</span> Garage
    </router-link>

    <div class="toggle glass" role="group" aria-label="Map style">
      <button
        type="button"
        :class="{ on: !satellite }"
        @click="emit('update:satellite', false)"
      >
        Map
      </button>
      <button
        type="button"
        :class="{ on: satellite }"
        @click="emit('update:satellite', true)"
      >
        Satellite
      </button>
    </div>
  </div>
</template>

<style scoped>
.chrome {
  position: absolute;
  top: calc(12px + env(safe-area-inset-top));
  left: 16px;
  right: 16px;
  z-index: 1000;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  pointer-events: none;
}
.chrome > * {
  pointer-events: auto;
}

.glass {
  background: var(--tm-glass-soft);
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.back {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  border-radius: var(--tm-r-pill);
  padding: 9px 15px 9px 12px;
  font-size: 13.5px;
  font-weight: 600;
  color: var(--tm-text);
}
.chev {
  font-size: 15px;
}

.toggle {
  display: inline-flex;
  flex-direction: column;
  gap: 2px;
  border-radius: var(--tm-r-md);
  padding: 5px;
}
.toggle button {
  border: none;
  background: none;
  cursor: pointer;
  font-size: 11px;
  font-weight: 700;
  color: var(--tm-text-mid);
  border-radius: 10px;
  padding: 5px 11px;
  transition: background 0.15s, color 0.15s;
}
.toggle button.on {
  color: var(--tm-bg);
  background: var(--tm-text);
}
</style>
