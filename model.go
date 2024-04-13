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

type Service []string

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
	ID          string
	Synced      bool
	CustomerID  string
	ProjectID   string // Project of a customer
	ServiceID   string // A service such as development
	Description string
	Start       time.Time
	End         time.Time
}

type ActivityRepository interface {
	GetAllActivities() ([]Activity, error)
	Add(activity Activity) error
}

type ProjectRepository interface {
	GetAllCustomers() ([]Customer, error)
	GetTimeEntries(start, end time.Time) ([]Activity, error)
	GetAllServices() ([]string, error)
	AddTimeEntry(activity Activity) error
}
