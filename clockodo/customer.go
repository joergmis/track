package clockodo

import (
	"context"
	"strings"

	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo/api"
)

var repo *repository

type repository struct {
	client *api.ClientWithResponses
}

func NewRepository(config Config) (track.ProjectRepository, error) {
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
	var projects []api.Project

	{
		response, err := r.client.GetV2ProjectsWithResponse(ctx)
		if err != nil {
			return data, err
		}

		projects = response.JSON200.Projects
	}

	{
		response, err := r.client.GetV2CustomersWithResponse(ctx)
		if err != nil {
			return data, err
		}

		customers := response.JSON200.Customers

		for _, c := range customers {
			customerProjects := []track.Project{}

			for _, p := range projects {
				if p.CustomersId == c.Id {
					customerProjects = append(customerProjects, track.Project{
						ID:       cleanup(p.Name),
						Name:     p.Name,
						Services: []track.Service{},
					})
				}
			}

			data = append(data, track.Customer{
				ID:       cleanup(c.Name),
				Name:     c.Name,
				Projects: customerProjects,
			})
		}
	}

	return data, nil
}

func cleanup(in string) string {
	replacer := strings.NewReplacer(
		" ", "_",
		"ä", "ae",
		"ö", "oe",
		"ü", "ue",
		"\"", "'",
	)

	return replacer.Replace(strings.ToLower(in))
}
