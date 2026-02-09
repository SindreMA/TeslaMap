export interface Car {
  id: number
  name: string | null
  model: string
  vin: string
  trim_badging: string | null
  exterior_color: string | null
}

export interface Position {
  latitude: number
  longitude: number
  speed: number | null
  battery_level: number | null
  heading: number | null
  date: string
}

export interface CarPosition {
  car: Car
  position: Position | null
}
