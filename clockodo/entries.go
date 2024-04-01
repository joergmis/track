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

	for _, entry := range response.JSON200.Entries {
		start, err := time.Parse(TimeLayoutString, entry.TimeSince)
		if err != nil {
			return data, err
		}

		end, err := time.Parse(TimeLayoutString, entry.TimeUntil)
		if err != nil {
			return data, err
		}

		data = append(data, track.Activity{
			Description: entry.Text,
			Start:       start,
			End:         end,
		})
	}

	return data, nil
}
