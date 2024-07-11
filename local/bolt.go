package local

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/joergmis/track"
)

type eventStorage struct {
	db *bolt.DB
}

func NewEventStorage() track.EventStorage {
	return &eventStorage{}
}

func (e *eventStorage) AddEvent(event track.Event) error {
	return nil
}

func (e *eventStorage) GetAllEvents() ([]track.Event, error) {
	return []track.Event{}, nil
}

func (e *eventStorage) GetAllEventsSince(date time.Time) ([]track.Event, error) {
	return []track.Event{}, nil
}
