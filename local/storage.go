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

type save struct {
	Activities map[track.ProjectBackendType][]track.Activity
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

	savedata := save{}

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

func (s *storage) GetActivities() (map[track.ProjectBackendType][]track.Activity, error) {
	data, err := s.getData()
	return data.Activities, err
}

func (s *storage) AddActivity(backend track.ProjectBackendType, activity track.Activity) error {
	data, err := s.getData()
	if err != nil {
		return err
	}

	data.Activities[backend] = append(data.Activities[backend], activity)

	return s.setData(data)
}

func (s *storage) getData() (save, error) {
	savedata := save{}
	data, err := os.ReadFile(filepath.Join(s.location, filename))
	if err != nil {
		return savedata, errors.Wrap(err, "read save data")
	}

	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&savedata); err != nil {
		return savedata, errors.Wrap(err, "decode save data")
	}

	return savedata, nil
}

func (s *storage) setData(savedata save) error {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(savedata); err != nil {
		return errors.Wrap(err, "encode save data")
	}

	if err := os.WriteFile(filepath.Join(s.location, filename), buf.Bytes(), 0644); err != nil {
		return errors.Wrap(err, "update save data")
	}

	return nil
}

func (s *storage) GetLastActivity(backend track.ProjectBackendType) (track.Activity, error) {
	activity := track.Activity{}

	data, err := s.getData()
	if err != nil {
		return activity, err
	}

	if len(data.Activities) == 0 {
		return activity, track.ErrNoActivities
	}

	sort.Slice(data.Activities[backend], func(i, j int) bool {
		return data.Activities[backend][i].StartTime.Before(data.Activities[backend][j].StartTime)
	})

	activity = data.Activities[backend][len(data.Activities)-1]

	return activity, nil
}

func (s *storage) UpdateActivity(backend track.ProjectBackendType, activity track.Activity) error {
	data, err := s.getData()
	if err != nil {
		return err
	}

	found := false

	for i, act := range data.Activities[backend] {
		if act.ID == activity.ID {
			found = true
			data.Activities[backend][i] = activity
		}
	}

	if !found {
		return track.ErrNoMatchingActivity
	}

	return s.setData(data)
}

func (s *storage) DeleteActivity(backend track.ProjectBackendType, activity track.Activity) error {
	data, err := s.getData()
	if err != nil {
		return err
	}

	index := 0
	found := false

	for i, act := range data.Activities[backend] {
		if act.ID == activity.ID {
			index = i
			found = true
		}
	}

	if !found {
		return track.ErrNoMatchingActivity
	}

	data.Activities[backend] = append(data.Activities[backend][:index], data.Activities[backend][index+1:]...)

	return s.setData(data)
}
