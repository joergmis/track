package clockodo

import (
	"context"
	"strings"

	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo/api"
	"github.com/pkg/errors"
)

func (r *repository) GetAllCustomers() ([]track.Customer, error) {
	ctx := context.Background()
	data := []track.Customer{}

	customers, err := r.getAllCustomers(ctx)
	if err != nil {
		return data, err
	}

	projects, err := r.getAllProjects(ctx)
	if err != nil {
		return data, err
	}

	for _, c := range customers {
		customerProjects := []track.Project{}

		for _, p := range projects {
			if p.CustomersId == c.Id {
				customerProjects = append(customerProjects, track.Project{
					ID:   cleanup(p.Name),
					Name: p.Name,
				})
			}
		}

		data = append(data, track.Customer{
			ID:       cleanup(c.Name),
			Name:     c.Name,
			Projects: customerProjects,
		})
	}

	return data, nil
}

func (r *repository) getAllCustomers(ctx context.Context) ([]api.Customer, error) {
	response, err := r.client.GetV2CustomersWithResponse(ctx)
	if err != nil {
		return []api.Customer{}, err
	}

	if response.JSON200 == nil {
		return []api.Customer{}, errors.New("no data received")
	}

	return response.JSON200.Customers, nil
}

func (r *repository) getAllProjects(ctx context.Context) ([]api.Project, error) {
	response, err := r.client.GetV2ProjectsWithResponse(ctx)
	if err != nil {
		return []api.Project{}, err
	}

	if response.JSON200 == nil {
		return []api.Project{}, errors.New("no data received")
	}

	return response.JSON200.Projects, nil
}

// cleanup names in order for them being consistent and not having whitespaces
// in the name since this creates (probably) some issues with the
// autocompletion on the commandline..
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
