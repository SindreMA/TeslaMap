package model

type Car struct {
	ID            int     `json:"id"`
	Name          *string `json:"name"`
	Model         string  `json:"model"`
	VIN           string  `json:"vin"`
	TrimBadging   *string `json:"trim_badging"`
	ExteriorColor *string `json:"exterior_color"`
}
