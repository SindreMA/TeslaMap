import { ref, onMounted } from 'vue'

export function useGeolocation() {
  const coords = ref<{ lat: number; lng: number } | null>(null)
  const error = ref<string | null>(null)
  const loading = ref(true)

  onMounted(() => {
    if (!navigator.geolocation) {
      error.value = 'Geolocation not supported'
      loading.value = false
      return
    }
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        coords.value = {
          lat: pos.coords.latitude,
          lng: pos.coords.longitude,
        }
        loading.value = false
      },
      (err) => {
        error.value = err.message
        loading.value = false
      },
      { enableHighAccuracy: true, timeout: 10000 }
    )
  })

  return { coords, error, loading }
}
