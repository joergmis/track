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

	Customer    string
	Project     string
	Service     string
	Description string

	StartTime time.Time
	EndTime   time.Time
}

func NewActivity(customer, project, service, description string) Activity {
	return Activity{
		ID:          uuid.New().String(),
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

type Backend interface {
	// GetAllCustomers returns a list with all customers.
	GetAllCustomers() ([]Customer, error)

	// GetAllServices returns a list including all services.
	GetAllServices() ([]string, error)

	// AddTimeEntry creates a new timeentry from the activity. It checks if
	// there is matching data (like customer or project) and returns an error
	// if this is not the case.
	AddTimeEntry(activity Activity) error
}
