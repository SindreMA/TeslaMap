<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { LMap, LTileLayer, LMarker, LPolyline } from '@vue-leaflet/vue-leaflet'
import 'leaflet/dist/leaflet.css'
import L from 'leaflet'
import type { Car } from '@/api/types'

const props = defineProps<{
  carCoords: { lat: number; lng: number }
  userCoords: { lat: number; lng: number } | null
  car: Car
  heading: number | null
  routeCoords: [number, number][] | null
  satellite: boolean
}>()

/* ---- markers: glowing rotated car tile + pulsing "you" dot ---- */
function buildCarIcon(heading: number): L.DivIcon {
  return L.divIcon({
    html: `
      <div style="position:relative;width:44px;height:44px;">
        <div style="position:absolute;inset:0;border-radius:50%;background:rgba(77,139,255,0.45);animation:tmPulse 2.4s ease-out infinite;"></div>
        <div style="position:absolute;top:50%;left:50%;transform:translate(-50%,-50%) rotate(${heading}deg);width:30px;height:30px;border-radius:10px;background:rgba(20,25,33,0.96);border:1.5px solid #4d8bff;box-shadow:0 0 16px rgba(77,139,255,0.7);display:flex;align-items:center;justify-content:center;">
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

const carIcon = computed(() => buildCarIcon(props.heading ?? 0))

const zoom = ref(14)

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

const center = computed(() => [props.carCoords.lat, props.carCoords.lng] as [number, number])

const bounds = computed(() => {
  if (!props.userCoords) return undefined
  return [
    [props.carCoords.lat, props.carCoords.lng],
    [props.userCoords.lat, props.userCoords.lng],
  ] as [[number, number], [number, number]]
})

const routeColor = computed(() => (props.satellite ? '#7fb0ff' : '#4d8bff'))

const polylinePoints = computed(() => {
  if (props.routeCoords?.length) return props.routeCoords
  if (!props.userCoords) return []
  return [
    [props.carCoords.lat, props.carCoords.lng],
    [props.userCoords.lat, props.userCoords.lng],
  ] as [number, number][]
})

const hasRoute = computed(() => !!props.routeCoords?.length)

const mapRef = ref<InstanceType<typeof LMap> | null>(null)

function leaflet(): L.Map | null {
  return (mapRef.value as unknown as { leafletObject?: L.Map })?.leafletObject ?? null
}

function recenter() {
  const map = leaflet()
  if (!map) return
  if (bounds.value) map.fitBounds(bounds.value, { padding: [60, 60] })
  else map.setView(center.value, zoom.value)
}

// auto-fit once the user's location arrives
watch(
  () => props.userCoords,
  (next) => {
    if (next) recenter()
  },
)

defineExpose({ recenter })
</script>

<template>
  <div class="map-container">
    <LMap
      ref="mapRef"
      :zoom="zoom"
      :center="center"
      :bounds="bounds"
      :options="{ zoomControl: false, attributionControl: true }"
      style="height: 100%; width: 100%"
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

      <LMarker :lat-lng="[carCoords.lat, carCoords.lng]" :icon="carIcon" />
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
