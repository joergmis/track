package clockodo

import (
	"context"

	"github.com/joergmis/track/clockodo/api"
	"github.com/pkg/errors"
)

func (r *repository) GetAllServices() ([]string, error) {
	services := []string{}

	raw, err := r.getAllServices()
	if err != nil {
		return services, err
	}

	for _, service := range raw {
		services = append(services, cleanup(service.Name))
	}

	return services, nil
}

func (r *repository) getAllServices() ([]api.Service, error) {
	response, err := r.client.GetV2ServicesWithResponse(context.Background())
	if err != nil {
		return []api.Service{}, err
	}

	if response.JSON200 == nil {
		return []api.Service{}, errors.New("no data found")
	}

	return response.JSON200.Services, nil
}
