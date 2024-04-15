package track

import (
	"time"

	"github.com/google/uuid"
)

type Clock interface {
	Now() time.Time
}

func NewActivity(customer, project, service, description string) Activity {
	return Activity{
		ID:          uuid.New().String(),
		Synced:      false,
		InProgress:  true,
		Customer:    customer,
		Project:     project,
		Service:     service,
		Description: description,
	}
}

func (a *Activity) Start() {
	a.StartTime = time.Now()
}

func (a *Activity) Stop() {
	// in case the activity is stopped when a new one is started, this makes
	// sure that they don't overlap
	a.EndTime = time.Now().Add(-1 * time.Second)
	a.InProgress = false
}
