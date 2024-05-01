package track_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/joergmis/track"
	"github.com/joergmis/track/local"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	firstactivity = track.Activity{
		ID:        "cooking",
		StartTime: time.Now().Add(-2 * time.Hour),
		EndTime:   time.Now().Add(-1 * time.Hour),
	}
	lastActivity = track.Activity{
		ID:        "coding",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
	}

	activities = []track.Activity{
		lastActivity,
		firstactivity,
	}
)

func TestStorage_GetLastActivity_withActivities(t *testing.T) {
	strg, err := local.NewStorage(filepath.Join(t.TempDir(), "entries.json"))
	if err != nil {
		t.Fatal(err)
	}

	for _, act := range activities {
		assert.Nil(t, strg.AddActivity(act))
	}

	last, err := strg.GetLastActivity()
	assert.Nil(t, err)

	assert.Equal(t, last.ID, lastActivity.ID)
}

func TestStorage_GetLastActivity_withoutActivities(t *testing.T) {
	strg, err := local.NewStorage(filepath.Join(t.TempDir(), "entries.json"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = strg.GetLastActivity()
	assert.Equal(t, track.ErrNoActivities, errors.Cause(err))
}

func TestStorage_GetAllActivities(t *testing.T) {
	strg, err := local.NewStorage(filepath.Join(t.TempDir(), "entries.json"))
	if err != nil {
		t.Fatal(err)
	}

	for _, act := range activities {
		assert.Nil(t, strg.AddActivity(act))
	}

	assert.Nil(t, strg.MarkActivityAsSynced(activities[0]))

	list, err := strg.GetAllActivities()
	assert.Nil(t, err)
	assert.Len(t, list, len(activities))
}

func TestStorage_GetUnsyncedActivities(t *testing.T) {
	strg, err := local.NewStorage(filepath.Join(t.TempDir(), "entries.json"))
	if err != nil {
		t.Fatal(err)
	}

	for _, act := range activities {
		assert.Nil(t, strg.AddActivity(act))
	}

	assert.Nil(t, strg.MarkActivityAsSynced(activities[0]))

	unsynced, err := strg.GetUnsyncedActivities()
	assert.Nil(t, err)
	assert.Len(t, unsynced, 1)
}

func TestStorage_UpdateActivity(t *testing.T) {
	strg, err := local.NewStorage(filepath.Join(t.TempDir(), "entries.json"))
	if err != nil {
		t.Fatal(err)
	}

	for _, act := range activities {
		assert.Nil(t, strg.AddActivity(act))
	}

	activity := activities[0]
	activity.Customer = "max muster ag"
	activity.Project = "house renovation"

	assert.Nil(t, strg.UpdateActivity(activity))

	list, err := strg.GetAllActivities()
	assert.Nil(t, err)

	for _, act := range list {
		if act.ID == activity.ID {
			assert.Equal(t, act.Customer, activity.Customer)
			assert.Equal(t, act.Project, activity.Project)
			assert.Equal(t, act.Service, activity.Service)
			assert.Equal(t, act.Description, activity.Description)
		}
	}
}

func TestStorage_DeleteActivity(t *testing.T) {
	strg, err := local.NewStorage(filepath.Join(t.TempDir(), "entries.json"))
	if err != nil {
		t.Fatal(err)
	}

	for _, act := range activities {
		assert.Nil(t, strg.AddActivity(act))
	}

	assert.Nil(t, strg.DeleteActivity(activities[0]))

	list, err := strg.GetAllActivities()
	assert.Nil(t, err)
	assert.Len(t, list, len(activities)-1)
}
