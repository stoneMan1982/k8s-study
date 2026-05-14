package event

import "time"

type Event struct {
	Type     string
	Payload  any
	From     string
	SendAt   time.Time
	HandleAt time.Time
}
