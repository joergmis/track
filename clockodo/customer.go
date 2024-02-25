package clockodo

import (
	"context"

	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo/api"
	"github.com/pkg/errors"
)

var repo *repository

type repository struct {
	client    *api.ClientWithResponses
	customers []api.Customer
	projects  []api.Project
}

func NewRepository(config Config) (track.Repository, error) {
	if repo == nil {
		repo = &repository{}

		client, err := newClockodoClient(config)
		if err != nil {
			return repo, err
		}

		repo.client = client
	}

	return repo, nil
}

func (r *repository) GetAllCustomers() ([]track.Customer, error) {
	ctx := context.Background()
	data := []track.Customer{}

	response, err := r.client.GetV2CustomersWithResponse(ctx)
	if err != nil {
		return data, err
	}

	r.customers = response.JSON200.Customers

	for _, c := range r.customers {
		data = append(data, track.Customer{
			Name: c.Name,
		})
	}

	return data, nil
}

func (r *repository) GetAllProjects() ([]track.Project, error) {
	ctx := context.Background()
	data := []track.Project{}

	response, err := r.client.GetV2ProjectsWithResponse(ctx)
	if err != nil {
		return data, err
	}

	r.projects = response.JSON200.Projects

	for _, p := range r.projects {
		data = append(data, track.Project{
			Name:      p.Name,
			Active:    p.Active,
			Completed: p.Completed,
		})
	}

	return data, nil
}

func (r *repository) GetCustomerProjects(customer string) ([]track.Project, error) {
	data := []track.Project{}
	id := -1

	// TODO: this assumes that the customers and or projects have already been
	// loaded. The question is if this should be handled as error or if we
	// should just load them here
	if len(r.customers) == 0 || len(r.projects) == 0 {
		return data, errors.Wrap(track.ErrNotInitialized, "no projects and/or customers found")
	}

	for _, c := range r.customers {
		if customer == c.Name {
			id = c.Id
		}
	}

	if id == -1 {
		return data, errors.Wrapf(track.ErrCustomerNotFound, "customer id %d", id)
	}

	for _, p := range r.projects {
		if p.CustomersId == id {
			data = append(data, track.Project{
				Name: p.Name,
			})
		}
	}

	return data, nil
}
