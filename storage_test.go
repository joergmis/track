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
