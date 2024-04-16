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

		t := tabby.New()
		t.AddHeader("time", "duration", "customer", "project", "description")

		total := time.Duration(0)

		for i, entry := range inRange {
			if i > 0 {
				// check for pauses between activities
				previous := inRange[i-1]
				if entry.StartTime.Sub(previous.EndTime).Minutes() > 5 {
					t.AddLine(
						fmt.Sprintf("%s - %s", previous.EndTime.Add(1*time.Second).Format(time.TimeOnly), entry.StartTime.Add(-1*time.Second).Format(time.TimeOnly)),
						fmt.Sprintf("%02d:%02d h", int(entry.StartTime.Sub(previous.EndTime).Hours()), int(entry.StartTime.Sub(previous.EndTime).Minutes())%60),
						"-- pause --",
						"--",
						"--",
					)
				}
			}

			if entry.InProgress {
				total += time.Since(entry.StartTime)
			} else {
				total += entry.EndTime.Sub(entry.StartTime)
			}

			t.AddLine(
				fmt.Sprintf("%s - %s", entry.StartTime.Format(time.TimeOnly), entry.EndTime.Format(time.TimeOnly)),
				getDuration(entry),
				entry.Customer,
				entry.Project,
				entry.Description,
			)
		}

		t.Print()

		fmt.Println("")

		t = tabby.New()
		t.AddHeader("total today")
		t.AddLine(fmt.Sprintf("%02d:%02d h", int(total.Hours()), int(total.Minutes())%60))
		t.Print()
	},
}

func getDuration(activity track.Activity) string {
	if activity.InProgress {
		return fmt.Sprintf(
			"%02d:%02d h (in progress)",
			int(time.Since(activity.StartTime).Hours()),
			int(time.Since(activity.StartTime).Minutes())%60,
		)
	}

	return fmt.Sprintf(
		"%02d:%02d h",
		int(activity.EndTime.Sub(activity.StartTime).Hours()),
		int(activity.EndTime.Sub(activity.StartTime).Minutes())%60,
	)
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
