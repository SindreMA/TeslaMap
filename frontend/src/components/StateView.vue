<script setup lang="ts">
defineProps<{
  variant: 'skeleton' | 'error' | 'stale'
  /** error: message + optional retry countdown */
  message?: string | null
  retrySeconds?: number | null
  /** stale: how long since the last fix + last known battery */
  lastSeen?: string | null
  batteryLevel?: number | null
}>()

const emit = defineEmits<{ (e: 'retry'): void }>()
</script>

<template>
  <!-- SKELETON -->
  <div v-if="variant === 'skeleton'" class="overlay">
    <div class="sheet skel-sheet">
      <span class="grab"></span>
      <div class="skel-row">
        <span class="sk" style="width: 140px; height: 16px"></span>
        <span class="sk" style="width: 70px; height: 16px"></span>
      </div>
      <span class="sk" style="width: 160px; height: 56px; margin: 4px 0 18px"></span>
      <span class="sk-static" style="width: 100%; height: 6px; margin-bottom: 22px"></span>
      <div class="skel-chips">
        <span class="sk" style="height: 64px"></span>
        <span class="sk" style="height: 64px"></span>
        <span class="sk" style="height: 64px"></span>
      </div>
    </div>
  </div>

  <!-- ERROR -->
  <div v-else-if="variant === 'error'" class="overlay center">
    <div class="dim"></div>
    <div class="card">
      <div class="badge err"><span>!</span></div>
      <h3>Couldn't reach your car</h3>
      <p>{{ message || "The car may be asleep or out of signal. We'll keep trying in the background." }}</p>
      <button class="cta" @click="emit('retry')">Try again</button>
      <div v-if="retrySeconds != null" class="retry">Retrying in {{ retrySeconds }}s…</div>
    </div>
  </div>

  <!-- STALE / NO LIVE POSITION -->
  <div v-else class="overlay bottom">
    <div class="sheet stale-sheet">
      <span class="not-live"><span class="ledg"></span>NOT LIVE</span>
      <h3>No live position</h3>
      <p>Your car hasn't reported in a while, so we can't time its arrival. Showing where it was last seen.</p>
      <div class="stale-meta">
        <div>
          <div class="ml">LAST SEEN</div>
          <div class="mv">{{ lastSeen || 'unknown' }}</div>
        </div>
        <div class="right">
          <div class="ml">BATTERY</div>
          <div class="mv batt">{{ batteryLevel == null ? '—' : batteryLevel + '%' }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.overlay {
  position: absolute;
  inset: 0;
  z-index: 1100;
}
.overlay.center {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}
.overlay.bottom,
.overlay:not(.center) {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
}
.dim {
  position: absolute;
  inset: 0;
  background: rgba(10, 12, 15, 0.55);
}

/* sheets */
.sheet {
  position: relative;
  background: var(--tm-glass);
  backdrop-filter: blur(28px);
  -webkit-backdrop-filter: blur(28px);
  border-top: 1px solid var(--tm-line-strong);
  border-radius: var(--tm-r-sheet) var(--tm-r-sheet) 0 0;
  padding: 18px 22px calc(28px + env(safe-area-inset-bottom));
  box-shadow: 0 -20px 50px rgba(0, 0, 0, 0.5);
}
.grab {
  display: block;
  width: 36px;
  height: 5px;
  border-radius: var(--tm-r-pill);
  background: rgba(255, 255, 255, 0.14);
  margin: 0 auto 18px;
}

/* skeleton */
.skel-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}
.skel-chips {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}
.sk {
  display: block;
  border-radius: 8px;
  background: linear-gradient(90deg, #181d25 25%, #222936 50%, #181d25 75%);
  background-size: 520px 100%;
  animation: tmShimmer 1.3s linear infinite;
}
.sk-static {
  display: block;
  border-radius: var(--tm-r-pill);
  background: #181d25;
}

/* error */
.card {
  position: relative;
  text-align: center;
  max-width: 320px;
}
.badge {
  width: 60px;
  height: 60px;
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
  font-size: 26px;
}
.badge.err {
  background: rgba(251, 106, 91, 0.12);
  border: 1px solid rgba(251, 106, 91, 0.35);
  color: var(--tm-low);
}
.card h3 {
  font-family: var(--tm-font-display);
  font-weight: 600;
  font-size: 21px;
  margin-bottom: 8px;
}
.card p {
  color: var(--tm-text-mid);
  font-size: 13.5px;
  line-height: 1.6;
  margin-bottom: 22px;
}
.cta {
  border: none;
  cursor: pointer;
  font-size: 14px;
  font-weight: 700;
  color: var(--tm-bg);
  background: var(--tm-accent);
  padding: 13px 26px;
  border-radius: var(--tm-r-pill);
  box-shadow: 0 8px 20px rgba(77, 139, 255, 0.4);
}
.retry {
  margin-top: 14px;
  font-size: 12px;
  color: var(--tm-text-dim);
}

/* stale */
.not-live {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: var(--tm-text-mid);
  background: rgba(255, 255, 255, 0.06);
  padding: 6px 12px;
  border-radius: var(--tm-r-pill);
  margin-bottom: 16px;
}
.ledg {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--tm-text-dim);
}
.stale-sheet h3 {
  font-family: var(--tm-font-display);
  font-weight: 600;
  font-size: 20px;
  margin-bottom: 8px;
}
.stale-sheet p {
  color: var(--tm-text-mid);
  font-size: 13.5px;
  line-height: 1.6;
  margin-bottom: 18px;
}
.stale-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 16px;
  border-top: 1px solid var(--tm-line);
}
.stale-meta .right {
  text-align: right;
}
.ml {
  font-size: 10px;
  color: var(--tm-text-dim);
  letter-spacing: 0.08em;
}
.mv {
  font-family: var(--tm-font-display);
  font-size: 16px;
  font-weight: 600;
  margin-top: 3px;
}
.mv.batt {
  color: var(--tm-mid);
}
</style>
