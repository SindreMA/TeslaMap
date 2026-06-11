<script setup lang="ts">
import { computed } from 'vue'
import type { Car, Position } from '@/api/types'
import BatteryIndicator from './BatteryIndicator.vue'

const props = withDefaults(
  defineProps<{
    car: Car
    /** live snapshot if we have one — drives battery + status */
    position?: Position | null
    /** featured row gets the accent border + explicit Track CTA */
    featured?: boolean
  }>(),
  { position: null, featured: false },
)

const displayName = computed(
  () => props.car.name || `${props.car.model} ${props.car.trim_badging || ''}`.trim(),
)

const subtitle = computed(() => {
  const c = props.car
  return [c.model, c.trim_badging].filter(Boolean).join(' · ')
})

const vinTail = computed(() => props.car.vin?.slice(-6) ?? '')

const moving = computed(() => (props.position?.speed ?? 0) > 0)

const lastSeen = computed(() => {
  if (!props.position?.date) return null
  const mins = Math.floor((Date.now() - new Date(props.position.date).getTime()) / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hrs = Math.floor(mins / 60)
  if (hrs < 24) return `${hrs}h ago`
  return `${Math.floor(hrs / 24)}d ago`
})

const status = computed(() => {
  if (!props.position) return 'No recent data'
  return moving.value ? 'Driving' : 'Parked'
})

// tint the little car body by the real exterior color when we know it
const PAINT: Record<string, string> = {
  white: '#f3f5f9',
  pearl: '#f3f5f9',
  black: '#23262d',
  solidblack: '#23262d',
  grey: '#5a6573',
  gray: '#5a6573',
  silver: '#8b94a3',
  midnightsilver: '#3c434e',
  blue: '#2f4a86',
  deepblue: '#1e2f56',
  red: '#c0322b',
}
const bodyColor = computed(() => {
  const raw = props.car.exterior_color?.toLowerCase().replace(/[^a-z]/g, '') ?? ''
  for (const key of Object.keys(PAINT)) if (raw.includes(key)) return PAINT[key]
  return '#cfd6e2'
})
</script>

<template>
  <router-link :to="`/car/${car.id}`" class="card" :class="{ featured }">
    <div class="top">
      <!-- thumbnail -->
      <div class="thumb">
        <div class="roof"></div>
        <div class="body" :style="{ background: `linear-gradient(${bodyColor}, ${bodyColor})` }"></div>
      </div>

      <!-- identity -->
      <div class="ident">
        <div class="name">{{ displayName }}</div>
        <div class="sub">{{ subtitle }}</div>
        <div class="meta">
          <BatteryIndicator v-if="position?.battery_level != null" :level="position.battery_level" variant="pill" />
          <span v-if="vinTail" class="vin">VIN ··{{ vinTail }}</span>
        </div>
      </div>

      <span v-if="!featured" class="chevron">›</span>
    </div>

    <!-- featured footer: status + primary action -->
    <div v-if="featured" class="footer">
      <span class="status">
        <span class="strong">{{ status }}</span>
        <template v-if="lastSeen"> · last seen {{ lastSeen }}</template>
      </span>
      <span class="track">Track <span class="arrow">→</span></span>
    </div>
  </router-link>
</template>

<style scoped>
.card {
  display: block;
  background: var(--tm-surface);
  border: 1px solid var(--tm-line);
  border-radius: 22px;
  padding: 16px;
  color: var(--tm-text);
  transition: border-color 0.2s, transform 0.15s;
}
.card:active {
  transform: scale(0.99);
}
.card.featured {
  background: linear-gradient(160deg, #161b24, #10141a);
  border-color: rgba(77, 139, 255, 0.45);
  box-shadow: 0 0 0 1px rgba(77, 139, 255, 0.12), 0 16px 40px rgba(0, 0, 0, 0.5);
}
.card:not(.featured):hover {
  border-color: rgba(77, 139, 255, 0.4);
}

.top {
  display: flex;
  gap: 14px;
  align-items: center;
}
.featured .top {
  align-items: flex-start;
}

/* striped tile with a car body, mirrors the handoff thumbnail */
.thumb {
  width: 84px;
  height: 64px;
  flex-shrink: 0;
  border-radius: var(--tm-r-md);
  background: repeating-linear-gradient(135deg, #1c222c, #1c222c 6px, #212834 6px, #212834 12px);
  border: 1px solid var(--tm-line);
  position: relative;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  overflow: hidden;
}
.thumb .body {
  width: 46px;
  height: 24px;
  border-radius: 11px 11px 6px 6px;
  margin-bottom: 8px;
  box-shadow: 0 3px 8px rgba(0, 0, 0, 0.4);
}
.thumb .roof {
  position: absolute;
  top: 6px;
  left: 7px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #f3f5f9;
  border: 2px solid var(--tm-surface);
}

.ident {
  flex: 1;
  min-width: 0;
}
.name {
  font-size: 17px;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.sub {
  font-size: 12.5px;
  color: var(--tm-text-mid);
  margin-top: 2px;
}
.meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 11px;
  flex-wrap: wrap;
}
.vin {
  font-size: 11px;
  color: var(--tm-text-dim);
  font-family: var(--tm-font-display);
}

.chevron {
  color: #5c6678;
  font-size: 20px;
  align-self: center;
}

.footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-top: 14px;
  padding-top: 14px;
  border-top: 1px solid var(--tm-line);
}
.status {
  font-size: 12.5px;
  color: var(--tm-text-mid);
}
.status .strong {
  color: var(--tm-text);
  font-weight: 600;
}
.track {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  font-size: 13px;
  font-weight: 700;
  color: var(--tm-bg);
  background: var(--tm-accent);
  padding: 9px 16px;
  border-radius: var(--tm-r-pill);
  box-shadow: 0 6px 16px var(--tm-accent-glow);
  white-space: nowrap;
}
</style>
