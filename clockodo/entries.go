package clockodo

import (
	"context"
	"log"
	"time"

	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo/api"
	"github.com/pkg/errors"
)

func (r *repository) GetTimeEntries(start, end time.Time) ([]track.Activity, error) {
	data := []track.Activity{}

	response, err := r.client.GetV2EntriesWithResponse(context.Background(), &api.GetV2EntriesParams{
		TimeSince: start.Format(TimeLayoutString),
		TimeUntil: end.Format(TimeLayoutString),
		Filter: api.EntriesFilter{
			UsersId: r.userID,
		},
	})
	if err != nil {
		return data, err
	}

	if response.JSON200 == nil {
		log.Println(string(response.Body))
		return data, errors.New("no data received")
	}

	customers, err := r.getAllCustomers(context.Background())
	if err != nil {
		return data, err
	}

	projects, err := r.getAllProjects(context.Background())
	if err != nil {
		return data, err
	}

	for _, entry := range response.JSON200.Entries {
		start, err := time.Parse(TimeLayoutString, entry.TimeSince)
		if err != nil {
			return data, err
		}

		end, err := time.Parse(TimeLayoutString, entry.TimeUntil)
		if err != nil {
			return data, err
		}

		activity := track.Activity{
			Description: entry.Text,
			StartTime:   start,
			EndTime:     end,
		}

		for _, project := range projects {
			if project.Id == entry.ProjectsId {
				activity.Project = cleanup(project.Name)
			}
		}

		for _, customer := range customers {
			if customer.Id == entry.CustomersId {
				activity.Customer = cleanup(customer.Name)
			}
		}

		data = append(data, activity)
	}

	return data, nil
}

func (r *repository) AddTimeEntry(activity track.Activity) error {
	var (
		customer api.Customer
		project  api.Project
		service  api.Service
	)

	customers, err := r.getAllCustomers(context.Background())
	if err != nil {
		return err
	}

	found := false
	for _, c := range customers {
		if activity.Customer == cleanup(c.Name) {
			found = true
			customer = c
		}
	}
	if !found {
		return track.ErrCustomerNotFound
	}

	projects, err := r.getAllProjects(context.Background())
	if err != nil {
		return err
	}

	found = false
	for _, p := range projects {
		if activity.Project == cleanup(p.Name) {
			project = p
			found = true
		}
	}
	if !found {
		return track.ErrProjectNotFound
	}

	services, err := r.getAllServices()
	if err != nil {
		return err
	}

	found = false
	for _, svc := range services {
		if activity.Service == cleanup(svc.Name) {
			service = svc
			found = true
		}
	}
	if !found {
		return track.ErrServiceNotFound
	}

	// use the project default for the billing setting
	billable := 0
	if project.BillableDefault {
		billable = 1
	}

	response, err := r.client.PostV2EntriesWithResponse(context.Background(), &api.PostV2EntriesParams{
		CustomersId: customer.Id,
		ServicesId:  service.Id,
		ProjectsId:  project.Id,
		Text:        activity.Description,
		UsersId:     r.userID,
		Billable:    api.PostV2EntriesParamsBillable(billable),
		TimeSince:   activity.StartTime.UTC().Format(TimeLayoutString),
		TimeUntil:   activity.EndTime.UTC().Format(TimeLayoutString),
	})
	if err != nil {
		log.Println(string(response.Body))
		return err
	}

	return nil
}
