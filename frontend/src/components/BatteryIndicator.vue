<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    /** 0–100; null renders a muted placeholder */
    level: number | null
    /** ring = SVG dial with the number beside it; pill = compact chip */
    variant?: 'ring' | 'pill'
    /** ring diameter in px (ring variant only) */
    size?: number
  }>(),
  { variant: 'pill', size: 34 },
)

/** good >= 60, mid >= 25, low below — matches the handoff battery scale. */
const color = computed(() => {
  const l = props.level
  if (l === null) return 'var(--tm-text-dim)'
  if (l >= 60) return 'var(--tm-good)'
  if (l >= 25) return 'var(--tm-mid)'
  return 'var(--tm-low)'
})

const tint = computed(() => {
  const l = props.level
  if (l === null) return 'rgba(255,255,255,0.06)'
  if (l >= 60) return 'rgba(52,211,153,0.12)'
  if (l >= 25) return 'rgba(251,191,36,0.12)'
  return 'rgba(251,106,91,0.12)'
})

// ring geometry
const r = computed(() => props.size / 2 - 4)
const circumference = computed(() => 2 * Math.PI * r.value)
const dashoffset = computed(() => {
  const pct = props.level === null ? 0 : Math.max(0, Math.min(100, props.level))
  return circumference.value * (1 - pct / 100)
})
const display = computed(() => (props.level === null ? '—' : String(props.level)))
</script>

<template>
  <!-- RING: dial + number, used in the telemetry stat chip -->
  <div v-if="variant === 'ring'" class="ring" :style="{ height: size + 'px' }">
    <svg :width="size" :height="size" :viewBox="`0 0 ${size} ${size}`">
      <circle
        :cx="size / 2"
        :cy="size / 2"
        :r="r"
        fill="none"
        stroke="rgba(255,255,255,0.1)"
        stroke-width="4"
      />
      <circle
        :cx="size / 2"
        :cy="size / 2"
        :r="r"
        fill="none"
        :stroke="color"
        stroke-width="4"
        stroke-linecap="round"
        :stroke-dasharray="circumference"
        :stroke-dashoffset="dashoffset"
        :transform="`rotate(-90 ${size / 2} ${size / 2})`"
        style="transition: stroke-dashoffset 0.6s ease, stroke 0.3s ease"
      />
    </svg>
    <div class="ring-meta">
      <div class="ring-value">{{ display }}<span class="ring-unit">%</span></div>
      <div class="ring-label">BATTERY</div>
    </div>
  </div>

  <!-- PILL: compact chip, used on cards and the collapsed sheet -->
  <span v-else class="pill" :style="{ color, background: tint }">
    <span class="dot" :style="{ background: color }"></span>{{ display }}%
  </span>
</template>

<style scoped>
.ring {
  display: flex;
  align-items: center;
  gap: 10px;
}
.ring-value {
  font-family: var(--tm-font-display);
  font-weight: 600;
  font-size: 17px;
  line-height: 1;
}
.ring-unit {
  font-size: 11px;
  color: var(--tm-text-dim);
}
.ring-label {
  font-size: 10px;
  color: var(--tm-text-dim);
  letter-spacing: 0.06em;
  margin-top: 3px;
}

.pill {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  font-weight: 700;
  padding: 4px 9px;
  border-radius: var(--tm-r-pill);
  white-space: nowrap;
}
.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}
</style>
