package clockodo

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/deepmap/oapi-codegen/v2/pkg/securityprovider"
	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo/api"
	"github.com/pkg/errors"
)

// TimeLayoutString is the layout that dates are expected to have when
// submitted via API.
const TimeLayoutString = "2006-01-02T15:04:05Z"

type Config struct {
	EmailAddress string
	ApiToken     string
}

type repository struct {
	userID int
	client *api.ClientWithResponses
}

func (r *repository) Type() track.ProjectBackendType {
	return track.ProjectBackendClockodo
}

func NewRepository(config Config) (track.ProjectRepository, error) {
	repo := &repository{}

	client, err := newClockodoClient(config)
	if err != nil {
		return repo, err
	}

	repo.client = client

	// Get the userID - match it via the email address from the config. The
	// userID will be necessary to list or add new time entries.
	// This also serves as a first check if the credentials are correct.
	response, err := repo.client.GetV2UsersWithResponse(context.Background())
	if err != nil {
		return repo, err
	}

	if response.JSON200 == nil {
		return repo, errors.New("no data received")
	}

	for _, user := range response.JSON200.Users {
		if user.Email == config.EmailAddress {
			repo.userID = user.Id
		}
	}

	return repo, nil
}

// newClockodoClient sets up a client to interact with the clockodo API. It
// uses the values from the config to authenticate against it via request
// headers (basic auth would theoretically also be supported).
func newClockodoClient(config Config) (*api.ClientWithResponses, error) {
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
	return api.NewClientWithResponses(
		"https://my.clockodo.com/api",
		api.WithRequestEditorFn(usernameProvider.Intercept),
		api.WithRequestEditorFn(apiKeyProvider.Intercept),
		api.WithRequestEditorFn(clientIdentificationProvider),
	)
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
