package clockodo

import (
	"context"

	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo/api"
)

var repo *repository

type repository struct {
	client    *api.ClientWithResponses
	customers []api.Customer
	projects  []api.Project
}

func NewCustomerRepository(config Config) (track.CustomerRepository, error) {
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

func (c *repository) GetAllCustomers() ([]track.Customer, error) {
	ctx := context.Background()
	data := []track.Customer{}

	response, err := c.client.GetV2CustomersWithResponse(ctx)
	if err != nil {
		return data, err
	}

	c.customers = response.JSON200.Customers

	for _, c := range c.customers {
		data = append(data, track.Customer{
			Name: c.Name,
		})
	}

	return data, nil
}
