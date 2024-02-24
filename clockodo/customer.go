package clockodo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/deepmap/oapi-codegen/v2/pkg/securityprovider"
	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo/api"
)

type Config struct {
	EmailAddress string
	ApiToken     string
}

type customerRepository struct {
	client *api.ClientWithResponses
}

func NewCustomerRepository(config Config) (track.CustomerRepository, error) {
	repo := &customerRepository{}

	usernameProvider, err := securityprovider.NewSecurityProviderApiKey(
		"header", "X-ClockodoApiUser", config.EmailAddress,
	)
	if err != nil {
		log.Fatal(err)
	}

	apiKeyProvider, err := securityprovider.NewSecurityProviderApiKey(
		"header", "X-ClockodoApiKey", config.ApiToken,
	)
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: this will use YOUR email address as 'technical contact':
	//
	//  Every request to our API must provide identification of the calling
	//  application, including the email address of a technical contact
	//  person. [...]
	//
	// Since you use this code, you are assumed to be technical enough. Check
	// out the API documentation for more information.
	clientIdentificationProvider := func(ctx context.Context, req *http.Request) error {
		req.Header.Add(
			"X-Clockodo-External-Application",
			fmt.Sprintf("github.com/joergmis/track;%s", config.EmailAddress),
		)
		return nil
	}

	// see https://www.clockodo.com/en/api/#c15088-headline for reference
	client, err := api.NewClientWithResponses(
		"https://my.clockodo.com/api",
		api.WithRequestEditorFn(usernameProvider.Intercept),
		api.WithRequestEditorFn(apiKeyProvider.Intercept),
		api.WithRequestEditorFn(clientIdentificationProvider),
	)
	if err != nil {
		return repo, err
	}

	repo.client = client

	return repo, nil
}

func (c *customerRepository) GetAllCustomers() ([]track.Customer, error) {
	ctx := context.Background()
	data := []track.Customer{}

	response, err := c.client.GetV2CustomersWithResponse(ctx)
	if err != nil {
		return data, err
	}

	for _, c := range response.JSON200.Customers {
		data = append(data, track.Customer{
			Name: c.Name,
		})
	}

	return data, nil
}
