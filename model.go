package track

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrNoCurrentActivity = errors.New("no active activity")
	ErrCustomerNotFound  = errors.New("customer not found")
	ErrProjectNotFound   = errors.New("project not found")
	ErrServiceNotFound   = errors.New("service not found")
	ErrNotInitialized    = errors.New("repository has not been initialized")
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

	// TODO: sync / in progress seem like implementation details? Maybe even
	// the ID...
	// Synced keeps track if the activity has been synced.
	Synced bool

	Customer    string
	Project     string
	Service     string
	Description string

	StartTime time.Time
	EndTime   time.Time
}

type ProjectRepository interface {
	GetAllCustomers() ([]Customer, error)
	GetTimeEntries(start, end time.Time) ([]Activity, error)
	GetAllServices() ([]string, error)
	AddTimeEntry(activity Activity) error
}
