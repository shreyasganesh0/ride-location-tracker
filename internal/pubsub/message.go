package pubsub

import (
)

type Messagetype int

const (
	Subscribe Messagetype = iota
	PublishLocation
)

type Message struct {

	Type Messagetype `json:"type"`

	Payload MessagePayload `json:"payload"`

}

type MessagePayload struct {
		Topic     string  `json:"topic"`
		Location  LatLng  `json:"location"`
} 

type LatLng struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
}
