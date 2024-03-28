package track

import (
	"time"

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

func (t *TimeTracking) Start(customerID, projectID, serviceID, description string) error {
	var (
		customer Customer
		project  Project
	)

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
				project = p
			}
		}

		if !found {
			return errors.Wrap(ErrProjectNotFound, "search for matching projectID")
		}
	}

	{
		found := false
		for _, s := range project.Services {
			if s.ID == serviceID {
				found = true
			}
		}

		if !found {
			return errors.Wrap(ErrServiceNotFound, "search for matching serviceID")
		}
	}

	if err := t.ActivityRepository.Add(Activity{
		CustomerID:  customerID,
		ProjectID:   projectID,
		ServiceID:   serviceID,
		Description: description,
		Start:       t.Clock.Now(),
	}); err != nil {
		return errors.Wrap(err, "add new activity")
	}

	return nil
}
