package local

import "github.com/joergmis/track"

type repo struct{}

func NewRepository() track.Repository {
	return &repo{}
}

func (r *repo) GetAllCustomers() ([]track.Customer, error) {
	return []track.Customer{
		{Name: "Addidas"},
		{Name: "Puma"},
		{Name: "Nike"},
	}, nil
}

func (r *repo) GetAllProjects() ([]track.Project, error) {
	return []track.Project{
		{
			Name:      "Test new shoes",
			Active:    true,
			Completed: false,
		},
		{
			Name:      "Code refactoring",
			Active:    true,
			Completed: false,
		},
		{
			Name:      "Make some burritos!",
			Active:    true,
			Completed: false,
		},
	}, nil
}

func (r *repo) GetCustomerProjects(customer string) ([]track.Project, error) {
	switch customer {
	case "Addidas":
		return []track.Project{
			{
				Name:      "Test new shoes",
				Active:    true,
				Completed: false,
			},
		}, nil

	case "Puma":
		return []track.Project{
			{
				Name:      "Code refactoring",
				Active:    true,
				Completed: false,
			},
		}, nil

	case "Nike":
		return []track.Project{
			{
				Name:      "Make some burritos!",
				Active:    true,
				Completed: false,
			},
		}, nil
	}

	return []track.Project{}, nil
}
