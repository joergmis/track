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

	response, err := r.client.GetV2WorkTimesWithResponse(context.Background(), &api.GetV2WorkTimesParams{
		DateSince: start.Format(time.DateOnly),
		DateUntil: end.Format(time.DateOnly),
	})
	if err != nil {
		return data, err
	}

	if response.JSON200 == nil {
		return data, errors.New("no data received")
	}

	for _, day := range response.JSON200.WorkTimeDays {
		for _, entry := range day.Intervals {
			start, err := time.Parse(TimeLayoutString, entry.TimeSince)
			if err != nil {
				return data, err
			}

			end, err := time.Parse(TimeLayoutString, entry.TimeUntil)
			if err != nil {
				return data, err
			}

			data = append(data, track.Activity{
				Start: start,
				End:   end,
			})
		}
	}

	return data, nil
}

func (r *repository) getEntriesTexts(start, end time.Time) (map[int]string, error) {
	data := map[int]string{}

	response, err := r.client.GetV2EntriesTextsWithResponse(context.Background(), &api.GetV2EntriesTextsParams{
		Text: "", // TODO: not allowed to be empty
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
