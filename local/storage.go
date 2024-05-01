package local

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"github.com/joergmis/track"
	"github.com/pkg/errors"
)

const filename = "entries.json"

type storage struct {
	location string
}

type savedata struct {
	Synced   []track.Activity
	Unsynced []track.Activity
}

func NewStorage(path string) (track.ActivityRepository, error) {
	strg := &storage{
		location: path,
	}

	data, err := os.ReadFile(filepath.Join(path, filename))
	if err != nil {
		if err := os.MkdirAll(path, 0755); err != nil {
			return strg, errors.Wrap(err, "create directory")
		}

		if err := os.WriteFile(filepath.Join(path, filename), []byte("{}"), 0644); err != nil {
			return strg, errors.Wrap(err, "initialize database")
		}

		data, err = os.ReadFile(filepath.Join(path, filename))
		if err != nil {
			return strg, errors.Wrap(err, "read/open file")
		}
	}

	savedata := savedata{}

	if len(data) < 2 {
		// this means most probably that there hasn't been any data previously
		// try to setup a basic json object so that the decode won't fail
		if err := os.WriteFile(filepath.Join(path, filename), []byte("{}"), 0644); err != nil {
			return strg, errors.Wrap(err, "initialize database")
		}
	}

	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&savedata); err != nil {
		return strg, errors.Wrap(err, "decode data")
	}

	return strg, nil
}

func (s *storage) GetAllActivities() ([]track.Activity, error) {
	data, err := s.getData()

	activities := []track.Activity{}

	activities = append(activities, data.Synced...)
	activities = append(activities, data.Unsynced...)

	return activities, err
}

func (s *storage) AddActivity(activity track.Activity) error {
	data, err := s.getData()
	if err != nil {
		return err
	}

	data.Unsynced = append(data.Unsynced, activity)

	return s.setData(data)
}

func (s *storage) getData() (savedata, error) {
	savedata := savedata{}
	data, err := os.ReadFile(filepath.Join(s.location, filename))
	if err != nil {
		return savedata, errors.Wrap(err, "read save data")
	}

	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&savedata); err != nil {
		return savedata, errors.Wrap(err, "decode save data")
	}

	return savedata, nil
}

func (s *storage) setData(savedata savedata) error {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(savedata); err != nil {
		return errors.Wrap(err, "encode save data")
	}

	if err := os.WriteFile(filepath.Join(s.location, filename), buf.Bytes(), 0644); err != nil {
		return errors.Wrap(err, "update save data")
	}

	return nil
}

func (s *storage) GetLastActivity() (track.Activity, error) {
	activity := track.Activity{}

	activities, err := s.GetAllActivities()
	if err != nil {
		return activity, errors.Wrap(err, "get all activities")
	}

	sort.Slice(activities, func(i, j int) bool {
		return activities[i].StartTime.Before(activities[j].StartTime)
	})

	if len(activities) == 0 {
		return activity, track.ErrNoActivities
	}

	activity = activities[len(activities)-1]

	return activity, nil
}

func (s *storage) UpdateActivity(activity track.Activity) error {
	data, err := s.getData()
	if err != nil {
		return err
	}

	found := false

	for i, act := range data.Unsynced {
		if act.ID == activity.ID {
			found = true
			data.Unsynced[i] = activity
		}
	}

	for i, act := range data.Synced {
		if act.ID == activity.ID {
			found = true
			data.Synced[i] = activity
		}
	}

	if !found {
		return track.ErrNoMatchingActivity
	}

	return s.setData(data)
}

func (s *storage) GetUnsyncedActivities() ([]track.Activity, error) {
	return []track.Activity{}, nil
}

func (s *storage) MarkActivityAsSynced(activity track.Activity) error {
	data, err := s.getData()
	if err != nil {
		return err
	}

	{
		// find the unsynced activity
		index := 0
		found := false

		for i, act := range data.Unsynced {
			if act.ID == activity.ID {
				index = i
				found = true
			}
		}

		if !found {
			return track.ErrNoMatchingActivity
		}

		data.Unsynced = append(data.Unsynced[:index], data.Unsynced[index+1:]...)
	}

	{
		// add the activity to the synced list
		found := false

		for _, act := range data.Synced {
			if act.ID == activity.ID {
				found = true
			}
		}

		if !found {
			data.Synced = append(data.Synced, activity)
		}

		// TODO: throw an error or just silently fail?
	}

	return s.setData(data)
}

func (s *storage) DeleteActivity(activity track.Activity) error {
	data, err := s.getData()
	if err != nil {
		return err
	}

	// TODO: fail if activity is found in neither list or silently fail?

	{
		index := 0
		found := false

		for i, act := range data.Synced {
			if act.ID == activity.ID {
				index = i
				found = true
			}
		}

		if found {
			data.Synced = append(data.Synced[:index], data.Synced[index+1:]...)
		}
	}

	{
		index := 0
		found := false

		for i, act := range data.Unsynced {
			if act.ID == activity.ID {
				index = i
				found = true
			}
		}

		if found {
			data.Unsynced = append(data.Unsynced[:index], data.Unsynced[index+1:]...)
		}
	}

	return s.setData(data)
}
