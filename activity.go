package track

import (
	"time"

	"github.com/google/uuid"
)

func NewActivity(customer, project, service, description string) Activity {
	return Activity{
		ID:          uuid.New().String(),
		Synced:      false,
		Customer:    customer,
		Project:     project,
		Service:     service,
		Description: description,
		EndTime:     time.Unix(0, 0),
	}
}

func (a *Activity) Start() {
	a.StartTime = time.Now()
}

func (a *Activity) Stop() {
	// in case the activity is stopped when a new one is started, this makes
	// sure that they don't overlap
	a.EndTime = time.Now().Add(-1 * time.Second)
}

func (a *Activity) Duration() time.Duration {
	// in case the activity is still in progress
	if !a.EndTime.After(a.StartTime) {
		return time.Since(a.StartTime)
	}

	return a.EndTime.Sub(a.StartTime)
}

func (a *Activity) InProgress() bool {
	return !a.EndTime.After(a.StartTime)
}
