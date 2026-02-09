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
