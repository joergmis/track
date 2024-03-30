package clockodo

import (
	"context"
	"strconv"
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

	texts, err := r.getEntriesTexts(start, end)
	if err != nil {
		return data, err
	}

	if response.JSON200 == nil {
		return data, errors.New("no data received")
	}

	for _, entry := range response.JSON200.Entries {
		text, ok := texts[entry.TextsId]
		if !ok {
			return data, errors.New("no matching description found")
		}

		data = append(data, track.Activity{
			Description: text,
		})
	}

	return data, nil
}

func (r *repository) getEntriesTexts(start, end time.Time) (map[int]string, error) {
	data := map[int]string{}

	response, err := r.client.GetV2EntriesTextsWithResponse(context.Background(), &api.GetV2EntriesTextsParams{
		Text: "",
		Filter: api.EntriesTextsFilter{
			TimeSince: start.Format(TimeLayoutString),
			TimeUntil: end.Format(TimeLayoutString),
			UsersId:   r.userID,
		},
	})
	if err != nil {
		return data, err
	}

	if response.JSON200 == nil {
		return data, errors.New("no data received")
	}

	for rawID, rawText := range response.JSON200.Texts {
		id, err := strconv.Atoi(rawID)
		if err != nil {
			return data, err
		}

		text, ok := rawText.(string)
		if !ok {
			return data, errors.New("data is not a string")
		}

		data[id] = text
	}

	return data, nil
}
