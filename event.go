package track

import (
	"time"
)

type EventType int

const (
	EventTypeAdd = EventType(iota)
	EventTypeUpdate
)

type Event struct {
	// Date acts as an ID for the event
	Date    time.Time
	Type    EventType
	Changes Activity
}

type EventStorage interface {
	// AddEvent adds a new event to the storage
	AddEvent(Event) error

	// GetAllEvents returns all events which are stored
	GetAllEvents() ([]Event, error)

	// GetAllEventsSince limits the events that are returned based on the
	// given time.
	GetAllEventsSince(time.Time) ([]Event, error)
}
