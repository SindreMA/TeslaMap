<script setup lang="ts">
import { computed, ref } from 'vue'
import BatteryIndicator from './BatteryIndicator.vue'

const props = defineProps<{
  carName: string
  etaMinutes: number | null
  distanceKm: number | null
  batteryLevel: number | null
  speed: number | null
  /** ISO date of the last position fix — drives the LIVE / "updated" label */
  updatedDate: string | null
  /** 0–1 share of the journey already covered, for the approach bar */
  progress: number
}>()

const expanded = ref(true)

const timeAgo = computed(() => {
  if (!props.updatedDate) return null
  const mins = Math.floor((Date.now() - new Date(props.updatedDate).getTime()) / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hrs = Math.floor(mins / 60)
  if (hrs < 24) return `${hrs}h ago`
  return `${Math.floor(hrs / 24)}d ago`
})

const eta = computed(() => (props.etaMinutes == null ? '—' : String(props.etaMinutes)))

const distanceLabel = computed(() =>
  props.distanceKm == null ? null : `${props.distanceKm.toFixed(1)} km`,
)

const arrivalClock = computed(() => {
  if (props.etaMinutes == null) return null
  const d = new Date(Date.now() + props.etaMinutes * 60000)
  return `${d.getHours()}:${String(d.getMinutes()).padStart(2, '0')}`
})

const barWidth = computed(() => `${Math.max(0, Math.min(1, props.progress)) * 100}%`)
const speedLabel = computed(() => (props.speed == null ? '—' : String(props.speed)))
</script>

<template>
  <div class="sheet" :class="{ collapsed: !expanded }">
    <button class="grab" :aria-label="expanded ? 'Collapse' : 'Expand'" @click="expanded = !expanded">
      <span class="bar"></span>
    </button>

    <!-- title row -->
    <div class="title">
      <div class="who">
        <span class="dot-car"></span>
        <span class="name">{{ carName }}</span>
      </div>
      <span class="live">
        <span class="led"></span>LIVE<template v-if="expanded"> · updated {{ timeAgo }}</template>
      </span>
    </div>

    <!-- EXPANDED ------------------------------------------------------ -->
    <template v-if="expanded">
      <div class="hero">
        <div class="count">
          <span class="big">{{ eta }}</span>
          <span class="min">min</span>
        </div>
        <div class="caption">until your car<br />reaches you</div>
      </div>

      <div class="approach">
        <div class="track">
          <div class="fill" :style="{ width: barWidth }"></div>
          <div class="knob" :style="{ left: barWidth }"></div>
        </div>
        <div class="approach-meta">
          <span>{{ distanceLabel ?? '—' }} to go</span>
          <span v-if="arrivalClock">arrives ~{{ arrivalClock }}</span>
        </div>
      </div>

      <div class="chips">
        <div class="chip chip-batt">
          <BatteryIndicator :level="batteryLevel" variant="ring" />
        </div>
        <div class="chip">
          <div class="chip-value">{{ speedLabel }}<span class="unit"> km/h</span></div>
          <div class="chip-label">SPEED</div>
        </div>
        <div class="chip">
          <div class="chip-value">{{ distanceKm == null ? '—' : distanceKm.toFixed(1) }}<span class="unit"> km</span></div>
          <div class="chip-label">DISTANCE</div>
        </div>
      </div>
    </template>

    <!-- COLLAPSED PEEK ----------------------------------------------- -->
    <template v-else>
      <div class="peek">
        <div class="peek-count">
          <span class="big-sm">{{ eta }}</span>
          <span class="min-sm">min</span>
        </div>
        <div class="peek-text">
          <div class="peek-line">until it reaches you</div>
          <div class="peek-sub">
            {{ distanceLabel ?? '—' }}<template v-if="arrivalClock"> · arrives ~{{ arrivalClock }}</template>
          </div>
        </div>
        <BatteryIndicator :level="batteryLevel" variant="pill" />
      </div>
      <div class="track thin">
        <div class="fill" :style="{ width: barWidth }"></div>
      </div>
      <div class="hint">▲ pull up for full telemetry</div>
    </template>
  </div>
</template>

<style scoped>
.sheet {
  background: var(--tm-glass);
  backdrop-filter: blur(28px);
  -webkit-backdrop-filter: blur(28px);
  border-top: 1px solid var(--tm-line-strong);
  border-radius: var(--tm-r-sheet) var(--tm-r-sheet) 0 0;
  padding: 14px 22px calc(26px + env(safe-area-inset-bottom));
  box-shadow: 0 -20px 50px rgba(0, 0, 0, 0.5);
}

.grab {
  display: block;
  width: 100%;
  border: none;
  background: none;
  padding: 0 0 16px;
  cursor: pointer;
}
.grab .bar {
  display: block;
  width: 38px;
  height: 5px;
  border-radius: var(--tm-r-pill);
  background: rgba(255, 255, 255, 0.18);
  margin: 0 auto;
}

.title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
}
.who {
  display: flex;
  align-items: center;
  gap: 9px;
}
.dot-car {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: #f3f5f9;
}
.name {
  font-size: 15px;
  font-weight: 700;
}
.live {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: var(--tm-good);
}
.led {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--tm-good);
  box-shadow: 0 0 8px var(--tm-good);
  animation: tmBlink 1.6s ease-in-out infinite;
}

/* hero */
.hero {
  display: flex;
  align-items: flex-end;
  gap: 14px;
}
.count {
  display: flex;
  align-items: baseline;
  gap: 8px;
}
.big {
  font-family: var(--tm-font-display);
  font-weight: 600;
  font-size: 76px;
  line-height: 0.85;
  letter-spacing: -0.03em;
}
.min {
  font-family: var(--tm-font-display);
  font-weight: 500;
  font-size: 28px;
  color: var(--tm-text-mid);
}
.caption {
  flex: 1;
  padding-bottom: 8px;
  font-size: 13.5px;
  color: var(--tm-text-mid);
  line-height: 1.35;
}

/* approach bar */
.approach {
  margin-top: 16px;
}
.track {
  position: relative;
  height: 6px;
  border-radius: var(--tm-r-pill);
  background: rgba(255, 255, 255, 0.08);
  overflow: hidden;
}
.track .fill {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  border-radius: var(--tm-r-pill);
  background: linear-gradient(90deg, var(--tm-accent-deep), var(--tm-accent));
  transition: width 0.6s ease;
}
.track .knob {
  position: absolute;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: var(--tm-accent);
  box-shadow: 0 0 10px var(--tm-accent), 0 0 0 3px #0f131a;
  transition: left 0.6s ease;
}
.track.thin {
  overflow: hidden;
}
.approach-meta {
  display: flex;
  justify-content: space-between;
  margin-top: 7px;
  font-size: 11px;
  color: var(--tm-text-dim);
}

/* stat chips */
.chips {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  margin-top: 18px;
}
.chip {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--tm-line);
  border-radius: 16px;
  padding: 12px;
}
.chip-batt {
  display: flex;
  align-items: center;
}
.chip-value {
  font-family: var(--tm-font-display);
  font-weight: 600;
  font-size: 17px;
  line-height: 1;
}
.chip-value .unit {
  font-size: 11px;
  color: var(--tm-text-dim);
}
.chip-label {
  font-size: 10px;
  color: var(--tm-text-dim);
  letter-spacing: 0.06em;
  margin-top: 6px;
}

/* collapsed peek */
.peek {
  display: flex;
  align-items: center;
  gap: 16px;
}
.peek-count {
  display: flex;
  align-items: baseline;
  gap: 6px;
}
.big-sm {
  font-family: var(--tm-font-display);
  font-weight: 600;
  font-size: 44px;
  line-height: 0.9;
}
.min-sm {
  font-family: var(--tm-font-display);
  font-size: 18px;
  color: var(--tm-text-mid);
}
.peek-text {
  flex: 1;
  min-width: 0;
}
.peek-line {
  font-size: 13px;
  color: var(--tm-text-mid);
}
.peek-sub {
  font-size: 11px;
  color: var(--tm-text-dim);
  margin-top: 2px;
}
.peek + .track.thin {
  margin-top: 16px;
}
.hint {
  text-align: center;
  font-size: 11px;
  color: #5c6678;
  margin-top: 14px;
}
</style>
