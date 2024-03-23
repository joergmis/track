package track

import (
	"github.com/pkg/errors"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrNotInitialized   = errors.New("repository has not been initialized")
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

type Repository interface {
	GetAllCustomers() ([]Customer, error)
}
