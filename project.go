package track

import (
	"github.com/pkg/errors"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrNotInitialized   = errors.New("repository has not been initialized")
)

type Customer struct {
	Name string
}

type Project struct {
	Name      string
	Active    bool
	Completed bool
}

type Repository interface {
	GetAllCustomers() ([]Customer, error)
	GetAllProjects() ([]Project, error)
	GetCustomerProjects(customer string) ([]Project, error)
}
