package steps

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/joergmis/track"
	"github.com/joergmis/track/local"
	"github.com/joergmis/track/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	testReference *testing.T
	activityRepo  track.ActivityRepository
	projectRepo   track.ProjectRepository
	customers     []track.Customer
	now           time.Time
	testfile      string
)

func givenTheDateIs(ctx context.Context, date string) (context.Context, error) {
	start, err := time.Parse(time.DateTime, date)
	if err != nil {
		return ctx, errors.Wrap(err, "parse start date string")
	}

	now = start

	return ctx, nil
}

func givenTheCustomerExists(ctx context.Context, customer, project string) (context.Context, error) {
	customers = append(customers, track.Customer{
		ID:   customer,
		Name: customer,
		Projects: []track.Project{
			{
				ID:   project,
				Name: project,
			},
		},
	})

	projectRepo.(*mocks.MockProjectRepository).EXPECT().GetAllCustomers().Maybe().Return(customers, nil)

	return ctx, nil
}

func givenThereIsAnAcitivityRunningFor(ctx context.Context, customer, project, service, description, date string) (context.Context, error) {
	_, err := time.Parse(time.DateTime, date)
	assert.Nil(testReference, err)

	err = activityRepo.AddActivity(track.Activity{
		ID:          uuid.New().String(),
		Synced:      false,
		InProgress:  true,
		Customer:    customer,
		Project:     project,
		Service:     service,
		Description: description,
		StartTime:   now.Add(-5 * time.Hour),
	})
	assert.Nil(testReference, err)

	return ctx, nil
}

func givenThereAreNoTimeEntriesPresent(ctx context.Context) (context.Context, error) {
	activities, err := activityRepo.GetActivities()
	assert.Nil(testReference, err)

	for _, activity := range activities {
		err := activityRepo.DeleteActivity(activity)
		assert.Nil(testReference, err)
	}

	return ctx, nil
}

func whenStartingANewActivityFor(ctx context.Context, customer, project, service, description string) (context.Context, error) {
	err := activityRepo.AddActivity(track.Activity{
		ID:          uuid.New().String(),
		Synced:      false,
		InProgress:  true,
		Customer:    customer,
		Project:     project,
		Service:     service,
		Description: description,
		StartTime:   now,
	})
	assert.Nil(testReference, err)

	return ctx, nil
}

func thenATimeentryIsAddedFor(ctx context.Context, date, customer, project, service, description string) (context.Context, error) {
	expectedStart, err := time.Parse(time.DateTime, date)
	assert.Nil(testReference, err)

	activity, err := activityRepo.GetLastActivity()
	assert.Nil(testReference, err)

	assert.Equal(testReference, customer, activity.Customer)
	assert.Equal(testReference, project, activity.Project)
	assert.Equal(testReference, service, activity.Service)
	assert.Equal(testReference, description, activity.Description)
	assert.Equal(testReference, expectedStart, activity.StartTime)

	return ctx, nil
}

func thenThereAreXEntries(ctx context.Context, count int) (context.Context, error) {
	activities, err := activityRepo.GetActivities()
	assert.Nil(testReference, err)

	assert.Equal(testReference, count, len(activities))

	return ctx, nil
}

func setup() {
	testfile = filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().Unix()))
	storage, err := local.NewStorage(testfile)
	assert.Nil(testReference, err)

	projectRepo = mocks.NewMockProjectRepository(testReference)
	activityRepo = storage

	customers = []track.Customer{}
}

func initializeScenario(ctx *godog.ScenarioContext) {
	setup()

	ctx.Given(`^the date is "([0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2})"$`, givenTheDateIs)
	ctx.Given(`^the customer "(\w+)" with project "(\w+)" exists$`, givenTheCustomerExists)
	ctx.Given(`^there is an activity running for "(\w+)" "(\w+)" "(\w+)" "(\w+)" started on "([0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2})"$`, givenThereIsAnAcitivityRunningFor)
	ctx.Given(`^no time entries are present$`, givenThereAreNoTimeEntriesPresent)

	ctx.When(`^starting a new activity for "(\w+)" "(\w+)" "(\w+)" "(\w+)"$`, whenStartingANewActivityFor)

	ctx.Then(`^a time entry is added for "([0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2})" "(\w+)" "(\w+)" "(\w+)" "(\w+)"$`, thenATimeentryIsAddedFor)
	ctx.Then(`^there are (\d+) time entries in the database$`, thenThereAreXEntries)
}

func TestFeatures(t *testing.T) {
	testReference = t

	suite := godog.TestSuite{
		ScenarioInitializer: initializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features"},
			TestingT: t,
		},
	}

	defer os.Remove(testfile)

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
