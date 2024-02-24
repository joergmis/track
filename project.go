package track

import "github.com/pkg/errors"

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrNotInitialized   = errors.New("repository has not been initialized")
)

type Customer struct {
	Name string
}

type CustomerRepository interface {
	GetAllCustomers() ([]Customer, error)
}

type Project struct {
	Name string
}

type ProjectRepository interface {
	GetAllProjects() ([]Project, error)
	GetCustomerProjects(customer Customer) ([]Project, error)
}
