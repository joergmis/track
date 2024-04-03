package clockodo

import (
	"context"
	"fmt"
)

func (r *repository) GetAllServices() ([]string, error) {
	response, err := r.client.GetV2ServicesWithResponse(context.Background())
	if err != nil {
		return []string{}, nil
	}

	if response.JSON200 == nil {
		return []string{}, nil
	}

	for _, service := range response.JSON200.Services {
		fmt.Println(service.Name)
	}

	return []string{}, nil
}
