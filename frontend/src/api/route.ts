export interface RouteData {
  distance_km: number
  duration_min: number
  routeCoords: [number, number][]
}

/** Decode a Google-encoded polyline string into [lat, lng] pairs. */
function decodePolyline(encoded: string): [number, number][] {
  const points: [number, number][] = []
  let index = 0
  let lat = 0
  let lng = 0

  while (index < encoded.length) {
    let shift = 0
    let result = 0
    let byte: number
    do {
      byte = encoded.charCodeAt(index++) - 63
      result |= (byte & 0x1f) << shift
      shift += 5
    } while (byte >= 0x20)
    lat += result & 1 ? ~(result >> 1) : result >> 1

    shift = 0
    result = 0
    do {
      byte = encoded.charCodeAt(index++) - 63
      result |= (byte & 0x1f) << shift
      shift += 5
    } while (byte >= 0x20)
    lng += result & 1 ? ~(result >> 1) : result >> 1

    points.push([lat / 1e5, lng / 1e5])
  }

  return points
}

/**
 * Fetch a driving route between two points using the public OSRM API.
 * Returns null if the request fails so callers can fall back gracefully.
 */
export async function fetchRoute(
  lat1: number, lng1: number,
  lat2: number, lng2: number,
): Promise<RouteData | null> {
  try {
    const url =
      `https://router.project-osrm.org/route/v1/driving/` +
      `${lng1},${lat1};${lng2},${lat2}?overview=full&geometries=polyline`

    const res = await fetch(url)
    if (!res.ok) return null

    const json = await res.json()
    if (json.code !== 'Ok' || !json.routes?.length) return null

    const route = json.routes[0]
    return {
      distance_km: route.distance / 1000,
      duration_min: Math.round(route.duration / 60),
      routeCoords: decodePolyline(route.geometry),
    }
  } catch {
    return null
  }
}
