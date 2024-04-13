package track

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Clock interface {
	Now() time.Time
}

type TimeTracking struct {
	Clock              Clock
	ProjectRepository  ProjectRepository
	ActivityRepository ActivityRepository
}

func NewTimeTracking(clock Clock, projectRepository ProjectRepository, activityRepository ActivityRepository) *TimeTracking {
	return &TimeTracking{
		Clock:              clock,
		ProjectRepository:  projectRepository,
		ActivityRepository: activityRepository,
	}
}

func (t *TimeTracking) Start(customerID, projectID, description string) error {
	var customer Customer

	customers, err := t.ProjectRepository.GetAllCustomers()
	if err != nil {
		return errors.Wrap(err, "retrieve customers")
	}

	{
		found := false
		for _, c := range customers {
			if c.ID == customerID {
				found = true
				customer = c
			}
		}

		if !found {
			return errors.Wrap(ErrCustomerNotFound, "search for matching customerID")
		}
	}

	{
		found := false
		for _, p := range customer.Projects {
			if p.ID == projectID {
				found = true
			}
		}

		if !found {
			return errors.Wrap(ErrProjectNotFound, "search for matching projectID")
		}
	}

	if err := t.ActivityRepository.Add(Activity{
		Customer:    customerID,
		Project:     projectID,
		Description: description,
		StartTime:   t.Clock.Now(),
	}); err != nil {
		return errors.Wrap(err, "add new activity")
	}

	return nil
}

func NewActivity(customer, project, service, description string) Activity {
	return Activity{
		ID:          uuid.New().String(),
		Synced:      false,
		InProgress:  true,
		Customer:    customer,
		Project:     project,
		Service:     service,
		Description: description,
	}
}

func (a *Activity) Start() {
	a.StartTime = time.Now()
}

func (a *Activity) Stop() {
	// in case the activity is stopped when a new one is started, this makes
	// sure that they don't overlap
	a.EndTime = time.Now().Add(-1 * time.Second)
	a.InProgress = false
}
