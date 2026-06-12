import { onUnmounted } from 'vue'

export interface Fix {
  lat: number
  lng: number
  heading: number | null
  /** event time of the fix (ms) — used only for relative pacing between fixes */
  t: number
}

// Segment timing bounds (ms).
const MIN_SEG = 350 // floor so tiny GPS jitter doesn't produce frantic micro-moves
const MAX_SEG = 5000 // cap so a long reporting gap doesn't make the car crawl
// Once more than this many fixes are queued (bursty/chunked delivery), start
// compressing segment durations so we catch up instead of lagging further behind.
const CATCHUP_AFTER = 3

function toRad(d: number) {
  return (d * Math.PI) / 180
}
function toDeg(r: number) {
  return (r * 180) / Math.PI
}

/** Compass bearing a→b in degrees, or null if the points are effectively equal. */
function bearing(aLat: number, aLng: number, bLat: number, bLng: number): number | null {
  const dLng = toRad(bLng - aLng)
  const la1 = toRad(aLat)
  const la2 = toRad(bLat)
  const y = Math.sin(dLng) * Math.cos(la2)
  const x = Math.cos(la1) * Math.sin(la2) - Math.sin(la1) * Math.cos(la2) * Math.cos(dLng)
  if (Math.abs(x) < 1e-12 && Math.abs(y) < 1e-12) return null
  return (toDeg(Math.atan2(y, x)) + 360) % 360
}

/** Interpolate an angle along the shortest arc. */
function lerpAngle(from: number, to: number, t: number): number {
  const diff = ((to - from + 540) % 360) - 180
  return (from + diff * t + 360) % 360
}

/**
 * Smoothly animates a vehicle marker between discrete position fixes.
 *
 * Fixes arrive irregularly (and sometimes in bursts) from the pipeline, so
 * instead of snapping the marker to each one we queue them and glide between
 * them, pacing each segment by the real time-gap between fixes. The result is a
 * marker that "drives" continuously, always trailing realtime by about one
 * update — exactly the buffered playback you'd want for live tracking.
 *
 * `onFrame` is called on every animation frame with the interpolated position;
 * apply it imperatively (e.g. marker.setLatLng) to avoid per-frame reactivity.
 */
export function useSmoothPosition(onFrame: (lat: number, lng: number, heading: number) => void) {
  let queue: Fix[] = []
  let curLat: number | null = null
  let curLng = 0
  let curHeading = 0

  let segFromLat = 0
  let segFromLng = 0
  let segTo: Fix | null = null
  let segStart = 0
  let segDuration = 0
  let lastT = 0
  let raf = 0

  function ensureLoop() {
    if (!raf) raf = requestAnimationFrame(tick)
  }

  function push(fix: Fix) {
    const lastKnownT = queue.length ? queue[queue.length - 1].t : (segTo?.t ?? lastT)
    if (curLat !== null && fix.t <= lastKnownT) return // duplicate or out-of-order

    if (curLat === null) {
      // first fix — snap into place, no animation
      curLat = fix.lat
      curLng = fix.lng
      curHeading = fix.heading ?? 0
      lastT = fix.t
      onFrame(curLat, curLng, curHeading)
      return
    }

    queue.push(fix)
    ensureLoop()
  }

  function beginSegment(now: number) {
    const next = queue.shift() as Fix
    segFromLat = curLat as number
    segFromLng = curLng
    segTo = next
    let dur = next.t - lastT
    if (queue.length >= CATCHUP_AFTER) dur = dur / (queue.length - CATCHUP_AFTER + 2)
    segDuration = Math.min(MAX_SEG, Math.max(MIN_SEG, dur))
    segStart = now
    lastT = next.t
  }

  function tick(now: number) {
    if (!segTo) {
      if (!queue.length) {
        raf = 0
        return
      }
      beginSegment(now)
    }
    const to = segTo as Fix
    const p = segDuration > 0 ? Math.min(1, (now - segStart) / segDuration) : 1

    curLat = segFromLat + (to.lat - segFromLat) * p
    curLng = segFromLng + (to.lng - segFromLng) * p

    const targetH = to.heading ?? bearing(segFromLat, segFromLng, to.lat, to.lng) ?? curHeading
    curHeading = lerpAngle(curHeading, targetH, p)
    onFrame(curLat, curLng, curHeading)

    if (p >= 1) {
      curHeading = targetH
      segTo = null
    }
    raf = requestAnimationFrame(tick)
  }

  /** Drop all buffered motion (e.g. when switching cars) so the next fix snaps. */
  function reset() {
    queue = []
    segTo = null
    curLat = null
    curLng = 0
    curHeading = 0
    lastT = 0
  }

  onUnmounted(() => {
    if (raf) cancelAnimationFrame(raf)
  })

  return { push, reset }
}
