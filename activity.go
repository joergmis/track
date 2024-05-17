package track

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

//go:generate go run scripts/main.go
var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrProjectNotFound  = errors.New("project not found")
	ErrServiceNotFound  = errors.New("service not found")
	ErrNotInitialized   = errors.New("repository has not been initialized")
)

type Customer struct {
	ID       string
	Name     string
	Projects []Project
}

type Project struct {
	ID   string
	Name string
}

type Activity struct {
	// ID is a uuid mainly for storage; this does not correspond to any
	// supported backend.
	ID string

	Backend BackendType

	Customer    string
	Project     string
	Service     string
	Description string

	StartTime time.Time
	EndTime   time.Time
}

func NewActivity(customer, project, service, description string, backend BackendType) Activity {
	return Activity{
		ID:          uuid.New().String(),
		Backend:     backend,
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
	a.EndTime = time.Now()
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
