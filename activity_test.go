package track_test

import (
	"testing"
	"time"

	"github.com/joergmis/track"
	"github.com/stretchr/testify/assert"
)

const (
	customer    = "max muster"
	project     = "muster project"
	service     = "development"
	description = "update website"
	backend     = "local"
)

func TestActivity_New(t *testing.T) {
	activity := track.NewActivity(customer, project, service, description, backend)

	assert.Equal(t, activity.StartTime, time.Unix(0, 0))
	assert.Equal(t, activity.EndTime, time.Unix(0, 0))
}

func TestActivity_Start(t *testing.T) {
	activity := track.NewActivity(customer, project, service, description, backend)

	now := time.Now()
	activity.Start()

	assert.InDelta(t, activity.StartTime.Unix(), now.Unix(), 5)
	assert.Equal(t, activity.EndTime, time.Unix(0, 0))
}

func TestActivity_Stop(t *testing.T) {
	activity := track.NewActivity(customer, project, service, description, backend)

	now := time.Now()
	activity.Stop()

	assert.Equal(t, activity.StartTime, time.Unix(0, 0))
	assert.InDelta(t, activity.EndTime.Unix(), now.Unix(), 5)
}

func TestActivity_Duration(t *testing.T) {
	activity := track.NewActivity(customer, project, service, description, backend)

	start := time.Now().Unix()
	activity.Start()
	end := time.Now().Unix()
	activity.Stop()

	assert.InDelta(t, end-start, activity.Duration().Seconds(), 5)
}

func TestActivity_InProgress(t *testing.T) {
	activity := track.NewActivity(customer, project, service, description, backend)

	assert.False(t, activity.InProgress())

	activity.Start()

	assert.True(t, activity.InProgress())
}
