package local

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/joergmis/track"
	"github.com/pkg/errors"
)

const filename = "entries.json"

type storage struct {
	location string
}

type save struct {
	Activities []track.Activity
}

func NewStorage(path string) (track.Storage, error) {
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

func (s *storage) GetActivities() ([]track.Activity, error) {
	data, err := s.getData()
	return data.Activities, err
}

func (s *storage) AddActivity(activity track.Activity) error {
	data, err := s.getData()
	if err != nil {
		return err
	}

	data.Activities = append(data.Activities, activity)

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
