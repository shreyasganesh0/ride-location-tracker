package broadcast

import (
)

type Message struct {

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	DriverID  string  `json:"driverId"`
}
