package track

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrNoCurrentActivity = errors.New("no active activity")
	ErrCustomerNotFound  = errors.New("customer not found")
	ErrNotInitialized    = errors.New("repository has not been initialized")
)

type Customer struct {
	ID       string
	Name     string
	Projects []Project
}

type Project struct {
	ID        string
	Name      string
	Active    bool
	Completed bool
	Services  []Service
}

type Service struct {
	ID   string
	Name string
}

type Activity struct {
	ID          string
	CustomerID  string
	ProjectID   string
	ServiceID   string
	Description string
	Start       time.Time
	End         time.Time
}

type ActivityRepository interface {
	// Start starts a new activity, stopping any previous activities if there
	// is still one running.
	Start(Activity) error
	// GetCurrentActivity returns the currently running activity or an error if
	// there is none running.
	GetCurrentActivity() (Activity, error)
	// Stop stops the currently active activity.
	Stop() error
}
type ProjectRepository interface {
	// GetAllCustomers returns all customers, including their projects and
	// services.
	GetAllCustomers() ([]Customer, error)
}
