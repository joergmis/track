package track

const (
	BackendLocal    = BackendType("local")
	BackendClockodo = BackendType("clockodo")
)

type BackendType string

func (b BackendType) Valid() bool {
	switch b {
	case BackendLocal:
		fallthrough
	case BackendClockodo:
		return true

	default:
		return false
	}
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
