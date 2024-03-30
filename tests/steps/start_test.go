package steps

import (
	"context"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/joergmis/track"
	"github.com/joergmis/track/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
)

var (
	testReference *testing.T
	timeTracking  *track.TimeTracking
	customers     []track.Customer
)

func givenTheDateIs(ctx context.Context, date string) (context.Context, error) {
	start, err := time.Parse(time.DateTime, date)
	if err != nil {
		return ctx, errors.Wrap(err, "parse start date string")
	}

	timeTracking.Clock.(*mocks.MockClock).EXPECT().Now().Maybe().Return(start)

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

	timeTracking.ProjectRepository.(*mocks.MockProjectRepository).EXPECT().GetAllCustomers().Maybe().Return(customers, nil)

	return ctx, nil
}

func givenThereIsAnAcitivityRunningFor(ctx context.Context, customer, project, description, date string) (context.Context, error) {
	_, err := time.Parse(time.DateTime, date)
	if err != nil {
		return ctx, errors.Wrap(err, "parse start date string")
	}

	timeTracking.ActivityRepository.(*mocks.MockActivityRepository).EXPECT().GetAllActivities().Maybe().Return([]track.Activity{
		{
			CustomerID:  customer,
			ProjectID:   project,
			Description: description,
		},
	}, nil)
	return ctx, nil
}

func whenStartingANewActivityFor(ctx context.Context, customer, project, description string) (context.Context, error) {
	timeTracking.ActivityRepository.(*mocks.MockActivityRepository).EXPECT().Add(mock.Anything).Maybe().Return(nil)

	err := timeTracking.Start(customer, project, description)
	if err != nil {
		return ctx, errors.Wrap(err, "start activity")
	}

	return ctx, nil
}

func thenATimeentryIsAddedFor(ctx context.Context, date, customer, project, description string) (context.Context, error) {
	expectedStart, err := time.Parse(time.DateTime, date)
	if err != nil {
		return ctx, errors.Wrap(err, "parse start date string")
	}

	timeTracking.ActivityRepository.(*mocks.MockActivityRepository).AssertCalled(testReference, "Add", track.Activity{
		CustomerID:  customer,
		ProjectID:   project,
		Description: description,
		Start:       expectedStart,
	})

	return ctx, nil
}

func setup() {
	timeTracking = track.NewTimeTracking(
		mocks.NewMockClock(testReference),
		mocks.NewMockProjectRepository(testReference),
		mocks.NewMockActivityRepository(testReference),
	)

	customers = []track.Customer{}
}

func initializeScenario(ctx *godog.ScenarioContext) {
	setup()

	ctx.Given(`^the date is "([0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2})"$`, givenTheDateIs)
	ctx.Given(`^the customer "(\w+)" with project "(\w+)" exists$`, givenTheCustomerExists)
	ctx.Given(`^there is an activity running for "(\w+)" "(\w+)" "(\w+)" started on "([0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2})"$`, givenThereIsAnAcitivityRunningFor)

	ctx.When(`^starting a new activity for "(\w+)" "(\w+)" "(\w+)"$`, whenStartingANewActivityFor)

	ctx.Then(`^a time entry is added for "([0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2})" "(\w+)" "(\w+)" "(\w+)"$`, thenATimeentryIsAddedFor)
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

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
