import type { Car, CarPosition } from './types'

const BASE = '/api'

export async function listCars(): Promise<Car[]> {
  const res = await fetch(`${BASE}/cars`)
  if (!res.ok) throw new Error('Failed to fetch cars')
  return res.json()
}

export async function getCarPosition(id: number): Promise<CarPosition> {
  const res = await fetch(`${BASE}/cars/${id}`)
  if (!res.ok) throw new Error('Car not found')
  return res.json()
}

export interface StreamHandlers {
  onData: (cp: CarPosition) => void
  /** fired on connection drop; EventSource then auto-reconnects */
  onError?: () => void
}

/**
 * Subscribe to live position updates over Server-Sent Events.
 * Returns an unsubscribe function. EventSource handles reconnection itself.
 * Falls back to 1s polling where EventSource is unavailable.
 */
export function streamCarPosition(id: number, handlers: StreamHandlers): () => void {
  if (typeof EventSource === 'undefined') {
    return pollCarPosition(id, handlers)
  }

  const es = new EventSource(`${BASE}/cars/${id}/stream`)
  es.onmessage = (e) => {
    try {
      handlers.onData(JSON.parse(e.data) as CarPosition)
    } catch {
      /* ignore malformed frame */
    }
  }
  es.onerror = () => handlers.onError?.()
  return () => es.close()
}

/** Legacy fallback: poll the REST endpoint once a second. */
function pollCarPosition(id: number, handlers: StreamHandlers): () => void {
  const tick = async () => {
    try {
      handlers.onData(await getCarPosition(id))
    } catch {
      handlers.onError?.()
    }
  }
  tick()
  const interval = setInterval(tick, 1000)
  return () => clearInterval(interval)
}
