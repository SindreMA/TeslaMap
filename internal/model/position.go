package model

import "time"

type Position struct {
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	Speed        *int       `json:"speed"`
	BatteryLevel *int       `json:"battery_level"`
	Heading      *float64   `json:"heading"`
	Date         time.Time  `json:"date"`
}

type CarPosition struct {
	Car      Car       `json:"car"`
	Position *Position `json:"position"`
}
