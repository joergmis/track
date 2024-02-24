package steps

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
	"github.com/joergmis/track"
	"github.com/joergmis/track/mocks"
	"github.com/pkg/errors"
)

var testReference *testing.T

type godogsCtxKey struct{}

func setupCustomerRepository(ctx context.Context) (context.Context, error) {
	repo := mocks.NewMockCustomerRepository(testReference)

	repo.On("GetAllCustomers").Return([]track.Customer{}, nil)

	return context.WithValue(ctx, godogsCtxKey{}, repo), nil
}

func addNewAcitvity(ctx context.Context) (context.Context, error) {
	repo, ok := ctx.Value(godogsCtxKey{}).(*mocks.MockCustomerRepository)
	if !ok {
		return ctx, errors.New("unexpected type of customer repository")
	}

	// TODO: this is just a dummy testcase

	_, _ = repo.GetAllCustomers()

	return ctx, nil
}

func fetchCustomers(ctx context.Context) (context.Context, error) {
	repo, ok := ctx.Value(godogsCtxKey{}).(*mocks.MockCustomerRepository)
	if !ok {
		return ctx, errors.New("unexpected type of customer repository")
	}

	repo.AssertCalled(testReference, "GetAllCustomers")

	return ctx, nil
}

func initializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given("^the customer repository is set up$", setupCustomerRepository)

	ctx.When("^adding a new activity$", addNewAcitvity)

	ctx.Then("^it should fetch the customers from the repository$", fetchCustomers)
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
