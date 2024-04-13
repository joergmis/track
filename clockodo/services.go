package clockodo

import (
	"context"

	"github.com/pkg/errors"
)

func (r *repository) GetAllServices() ([]string, error) {
	services := []string{}

	response, err := r.client.GetV2ServicesWithResponse(context.Background())
	if err != nil {
		return services, err
	}

	if response.JSON200 == nil {
		return services, errors.New("no data found")
	}

	for _, service := range response.JSON200.Services {
		services = append(services, service.Name)
	}

	return services, nil
}
