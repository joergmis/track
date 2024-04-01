package clockodo

import (
	"context"
	"log"
	"time"

	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo/api"
	"github.com/pkg/errors"
)

func (r *repository) GetTimeEntries(start, end time.Time) ([]track.Activity, error) {
	data := []track.Activity{}

	response, err := r.client.GetV2EntriesWithResponse(context.Background(), &api.GetV2EntriesParams{
		TimeSince: start.Format(TimeLayoutString),
		TimeUntil: end.Format(TimeLayoutString),
		Filter: api.EntriesFilter{
			UsersId: r.userID,
		},
	})
	if err != nil {
		return data, err
	}

	if response.JSON200 == nil {
		log.Println(string(response.Body))
		return data, errors.New("no data received")
	}

	customers, err := r.getAllCustomers(context.Background())
	if err != nil {
		return data, err
	}

	projects, err := r.getAllProjects(context.Background())
	if err != nil {
		return data, err
	}

	for _, entry := range response.JSON200.Entries {
		start, err := time.Parse(TimeLayoutString, entry.TimeSince)
		if err != nil {
			return data, err
		}

		end, err := time.Parse(TimeLayoutString, entry.TimeUntil)
		if err != nil {
			return data, err
		}

		activity := track.Activity{
			Description: entry.Text,
			Start:       start,
			End:         end,
		}

		for _, project := range projects {
			if project.Id == entry.ProjectsId {
				activity.ProjectID = cleanup(project.Name)
			}
		}

		for _, customer := range customers {
			if customer.Id == entry.CustomersId {
				activity.CustomerID = cleanup(customer.Name)
			}
		}

		data = append(data, activity)
	}

	return data, nil
}
