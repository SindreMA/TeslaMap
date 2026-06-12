<script setup lang="ts">
import { computed, ref, watch, onUnmounted } from 'vue'
import { LMap, LTileLayer, LMarker, LPolyline } from '@vue-leaflet/vue-leaflet'
import 'leaflet/dist/leaflet.css'
import L from 'leaflet'
import type { Car } from '@/api/types'
import { useSmoothPosition, type Fix } from '@/composables/useSmoothPosition'

const props = defineProps<{
  carFix: Fix
  userCoords: { lat: number; lng: number } | null
  car: Car
  routeCoords: [number, number][] | null
  satellite: boolean
}>()

const zoom = ref(14)
// Captured once: the map frames here on first render and then stays put — it
// must NOT follow every position update, or the whole map jumps. The car
// marker glides within it instead; use the recenter button to re-frame.
const initialCenter: [number, number] = [props.carFix.lat, props.carFix.lng]

/* ---- imperative, gliding car marker (so the pulse never resets and we can
       rotate without rebuilding the icon every frame) ---- */
function buildCarIcon(): L.DivIcon {
  return L.divIcon({
    html: `
      <div style="position:relative;width:44px;height:44px;">
        <div style="position:absolute;inset:0;border-radius:50%;background:rgba(77,139,255,0.45);animation:tmPulse 2.4s ease-out infinite;"></div>
        <div class="tm-rotor" style="position:absolute;top:50%;left:50%;transform:translate(-50%,-50%) rotate(0deg);width:30px;height:30px;border-radius:10px;background:rgba(20,25,33,0.96);border:1.5px solid #4d8bff;box-shadow:0 0 16px rgba(77,139,255,0.7);display:flex;align-items:center;justify-content:center;">
          <div style="width:11px;height:17px;border-radius:5px 5px 4px 4px;background:linear-gradient(#eef1f6,#c7d0df);"></div>
        </div>
      </div>`,
    className: '',
    iconSize: [44, 44],
    iconAnchor: [22, 22],
  })
}

const userIcon = L.divIcon({
  html: `
    <div style="position:relative;width:30px;height:30px;">
      <div style="position:absolute;inset:0;border-radius:50%;background:rgba(255,255,255,0.5);animation:tmPulse 2.4s ease-out infinite;"></div>
      <div style="position:absolute;top:50%;left:50%;transform:translate(-50%,-50%);width:13px;height:13px;border-radius:50%;background:#fff;border:2px solid #0d1016;box-shadow:0 2px 6px rgba(0,0,0,0.5);"></div>
    </div>`,
  className: '',
  iconSize: [30, 30],
  iconAnchor: [15, 15],
})

let map: L.Map | null = null
let carMarker: L.Marker | null = null
let rotor: HTMLElement | null = null

const smooth = useSmoothPosition((lat, lng, heading) => {
  carMarker?.setLatLng([lat, lng])
  if (rotor) rotor.style.transform = `translate(-50%,-50%) rotate(${heading}deg)`
})

function onMapReady(m: L.Map) {
  map = m
  carMarker = L.marker(initialCenter, {
    icon: buildCarIcon(),
    interactive: false,
    keyboard: false,
    zIndexOffset: 1000,
  }).addTo(m)
  rotor = (carMarker.getElement()?.querySelector('.tm-rotor') as HTMLElement) ?? null
  smooth.push({ ...props.carFix })
  recenter()
}

// feed each new fix into the smoother (it glides the marker)
watch(
  () => props.carFix,
  (f) => smooth.push({ lat: f.lat, lng: f.lng, heading: f.heading, t: f.t }),
)

// switching cars: drop buffered motion and snap to the new car
watch(
  () => props.car?.id,
  () => {
    smooth.reset()
    smooth.push({ ...props.carFix })
  },
)

/* ---- tiles ---- */
const tileUrl = computed(() =>
  props.satellite
    ? 'https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}'
    : 'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png',
)
const tileAttribution = computed(() =>
  props.satellite
    ? '&copy; <a href="https://www.esri.com/">Esri</a>'
    : '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> &copy; <a href="https://carto.com/">CARTO</a>',
)

/* ---- route ---- */
const routeColor = computed(() => (props.satellite ? '#7fb0ff' : '#4d8bff'))
const hasRoute = computed(() => !!props.routeCoords?.length)
const polylinePoints = computed(() => {
  if (props.routeCoords?.length) return props.routeCoords
  if (!props.userCoords) return []
  return [
    [props.carFix.lat, props.carFix.lng],
    [props.userCoords.lat, props.userCoords.lng],
  ] as [number, number][]
})

/* ---- framing (only on demand, never per update) ---- */
function boundsNow(): [[number, number], [number, number]] | null {
  if (!props.userCoords) return null
  return [
    [props.carFix.lat, props.carFix.lng],
    [props.userCoords.lat, props.userCoords.lng],
  ]
}

function recenter() {
  if (!map) return
  const b = boundsNow()
  if (b) map.fitBounds(b, { padding: [60, 60] })
  else map.setView([props.carFix.lat, props.carFix.lng], zoom.value)
}

// frame both pins once the user's location arrives (fires ~once)
watch(
  () => props.userCoords,
  (next) => {
    if (next) recenter()
  },
)

onUnmounted(() => {
  carMarker?.remove()
})

defineExpose({ recenter })
</script>

<template>
  <div class="map-container">
    <LMap
      :zoom="zoom"
      :center="initialCenter"
      :options="{ zoomControl: false, attributionControl: true }"
      style="height: 100%; width: 100%"
      @ready="onMapReady"
    >
      <LTileLayer
        :key="satellite ? 'sat' : 'dark'"
        :url="tileUrl"
        :attribution="tileAttribution"
      />

      <!-- route: soft glow underlay + crisp line on top -->
      <LPolyline
        v-if="userCoords && hasRoute"
        :lat-lngs="polylinePoints"
        :color="routeColor"
        :weight="13"
        :opacity="0.22"
      />
      <LPolyline
        v-if="userCoords"
        :lat-lngs="polylinePoints"
        :color="routeColor"
        :weight="hasRoute ? 4.5 : 3"
        :opacity="hasRoute ? 0.9 : 0.7"
        :dash-array="hasRoute ? undefined : '8, 10'"
      />

      <!-- car marker is added imperatively in onMapReady so it can glide -->
      <LMarker v-if="userCoords" :lat-lng="[userCoords.lat, userCoords.lng]" :icon="userIcon" />
    </LMap>
  </div>
</template>

<style scoped>
.map-container {
  position: absolute;
  inset: 0;
  background: #0e1217;
}

/* tidy up leaflet's default chrome to match the dark theme */
.map-container :deep(.leaflet-control-attribution) {
  background: rgba(10, 12, 15, 0.6);
  color: var(--tm-text-dim);
  font-size: 10px;
}
.map-container :deep(.leaflet-control-attribution a) {
  color: var(--tm-text-mid);
}
</style>
