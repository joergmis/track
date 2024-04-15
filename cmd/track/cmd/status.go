package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/cheynewallace/tabby"
	"github.com/joergmis/track"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get an overview over the timetracking",
	Run: func(cmd *cobra.Command, args []string) {
		activities, err := storage.GetActivities()
		if err != nil {
			log.Fatalf("get all activities from storage: %v", err)
		}

		inRange := []track.Activity{}

		start := time.Now().Add(-24 * time.Hour)
		end := time.Now().Add(1 * time.Hour)

		t := tabby.New()
		t.AddHeader("time", "duration", "customer", "project", "description")

		for _, entry := range activities {
			if entry.StartTime.Before(start) || entry.StartTime.After(end) || entry.EndTime.After(end) {
				continue
			}

			inRange = append(inRange, entry)
		}

		if len(inRange) == 0 {
			fmt.Println("-- no activities in range --")
			return
		}

		for _, entry := range inRange {
			t.AddLine(
				fmt.Sprintf("%s - %s", entry.StartTime.Format(time.TimeOnly), entry.EndTime.Format(time.TimeOnly)),
				getDuration(entry),
				entry.Customer,
				entry.Project,
				entry.Description,
			)
		}

		t.Print()
	},
}

func getDuration(activity track.Activity) string {
	duration := int(activity.EndTime.Sub(activity.StartTime).Minutes())

	if duration >= 0 {
		return fmt.Sprintf("%d min", duration)
	}

	return "in progress"
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
