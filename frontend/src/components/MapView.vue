<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { LMap, LTileLayer, LMarker, LPopup, LPolyline } from '@vue-leaflet/vue-leaflet'
import 'leaflet/dist/leaflet.css'
import L from 'leaflet'
import type { Car } from '@/api/types'

const SIZE = 48

function buildCarIcon(heading: number): L.DivIcon {
  const rot = heading // SVG points up = north, heading 0 = north
  return L.divIcon({
    html: `<img src="/cars/car-top.svg" style="width:${SIZE}px;height:${SIZE}px;transform:rotate(${rot}deg);filter:brightness(0) invert(1) drop-shadow(0 0 4px rgba(0,0,0,0.7));" />`,
    className: '',
    iconSize: [SIZE, SIZE],
    iconAnchor: [SIZE / 2, SIZE / 2],
    popupAnchor: [0, -SIZE / 2],
  })
}

const personIcon = L.divIcon({
  html: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48" width="36" height="36">
    <circle cx="24" cy="24" r="22" fill="#166534" stroke="#22c55e" stroke-width="2"/>
    <circle cx="24" cy="16" r="4" fill="white"/>
    <path d="M18 32c0-3.3 2.7-6 6-6s6 2.7 6 6" stroke="white" stroke-width="3" stroke-linecap="round" fill="none"/>
  </svg>`,
  className: '',
  iconSize: [36, 36],
  iconAnchor: [18, 18],
  popupAnchor: [0, -18],
})

const props = defineProps<{
  carCoords: { lat: number; lng: number }
  userCoords: { lat: number; lng: number } | null
  car: Car
  heading: number | null
  routeCoords: [number, number][] | null
}>()

const carIcon = computed(() => buildCarIcon(props.heading ?? 0))

const zoom = ref(14)
const satellite = ref(false)

const tileUrl = computed(() =>
  satellite.value
    ? 'https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}'
    : 'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png'
)

const tileAttribution = computed(() =>
  satellite.value
    ? '&copy; <a href="https://www.esri.com/">Esri</a>'
    : '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> &copy; <a href="https://carto.com/">CARTO</a>'
)

const center = computed(() => [props.carCoords.lat, props.carCoords.lng] as [number, number])

const bounds = computed(() => {
  if (!props.userCoords) return undefined
  return [
    [props.carCoords.lat, props.carCoords.lng],
    [props.userCoords.lat, props.userCoords.lng],
  ] as [[number, number], [number, number]]
})

const polylinePoints = computed(() => {
  if (props.routeCoords?.length) return props.routeCoords
  if (!props.userCoords) return []
  return [
    [props.carCoords.lat, props.carCoords.lng],
    [props.userCoords.lat, props.userCoords.lng],
  ] as [number, number][]
})

const mapRef = ref<InstanceType<typeof LMap> | null>(null)

watch(() => props.userCoords, (newCoords) => {
  if (newCoords && mapRef.value) {
    const map = (mapRef.value as any).leafletObject
    if (map) {
      map.fitBounds(bounds.value, { padding: [50, 50] })
    }
  }
})
</script>

<template>
  <div class="map-container">
  Hei
    <LMap
      ref="mapRef"
      :zoom="zoom"
      :center="center"
      :bounds="bounds"
      :options="{ zoomControl: true }"
      style="height: 100%; width: 100%"
    >
      <LTileLayer
        :key="satellite ? 'sat' : 'dark'"
        :url="tileUrl"
        :attribution="tileAttribution"
      />

      <!-- Car marker -->
      <LMarker :lat-lng="[carCoords.lat, carCoords.lng]" :icon="carIcon">
        <LPopup>Your car is here</LPopup>
      </LMarker>

      <!-- User marker -->
      <LMarker v-if="userCoords" :lat-lng="[userCoords.lat, userCoords.lng]" :icon="personIcon">
        <LPopup>Your location is here</LPopup>
      </LMarker>

      <!-- Route between car and user -->
      <LPolyline
        v-if="userCoords"
        :lat-lngs="polylinePoints"
        :color="'#3b82f6'"
        :weight="3"
        :dash-array="routeCoords?.length ? undefined : '8, 8'"
      />
    </LMap>
    <button class="toggle-satellite" @click="satellite = !satellite">
      {{ satellite ? 'Map' : 'Satellite' }}
    </button>
  </div>
</template>

<style scoped>
.map-container {
  flex: 1;
  min-height: 300px;
  position: relative;
}

.toggle-satellite {
  position: absolute;
  top: 12px;
  right: 12px;
  z-index: 1000;
  padding: 6px 14px;
  background: #1a1a1a;
  color: #e0e0e0;
  border: 1px solid #3a3a3a;
  border-radius: 6px;
  font-size: 0.8rem;
  cursor: pointer;
  transition: background 0.15s;
}

.toggle-satellite:hover {
  background: #2a2a2a;
}
</style>
