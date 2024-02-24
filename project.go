package track

type Customer struct {
	Name string
}

type CustomerRepository interface {
	GetAllCustomers() ([]Customer, error)
}

type Project struct{}
