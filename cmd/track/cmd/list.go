package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/cheynewallace/tabby"
	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo"
	"github.com/spf13/cobra"
)

var (
	listEntriesFrom string
	listEntriesTo   string

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List existing time entries",
		Run: func(cmd *cobra.Command, args []string) {
			activities, err := storage.GetActivities()
			if err != nil {
				log.Fatalf("get all activities from storage: %v", err)
			}

			start, err := time.Parse(clockodo.TimeLayoutString, listEntriesFrom)
			cobra.CheckErr(err)

			end, err := time.Parse(clockodo.TimeLayoutString, listEntriesTo)
			cobra.CheckErr(err)

			t := tabby.New()
			t.AddHeader("time", "duration", "customer", "project", "description")

			for i, entry := range activities {
				if i > 0 {
					// check for pauses between entries
					previous := activities[i-1]
					if entry.StartTime.Sub(previous.EndTime) > time.Duration(1)*time.Minute {
						t.AddLine(
							fmt.Sprintf("%s - %s", previous.EndTime.Format(time.TimeOnly), entry.StartTime.Format(time.TimeOnly)),
							fmt.Sprintf("%d min", int(entry.StartTime.Sub(previous.EndTime).Minutes())),
							"-- pause --",
							"--",
							"--",
						)
					}
				}

				if entry.StartTime.Before(start) || entry.StartTime.After(end) || entry.EndTime.After(end) {
					continue
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
		},
	}
)

func getDuration(activity track.Activity) string {
	duration := int(activity.EndTime.Sub(activity.StartTime).Minutes())

	if duration >= 0 {
		return fmt.Sprintf("%d min", duration)
	}

	return "in progress"
}

func init() {
	rootCmd.AddCommand(listCmd)

	start := time.Now().Add(-(time.Hour * 24)).Format(clockodo.TimeLayoutString)
	end := time.Now().Format(clockodo.TimeLayoutString)

	listCmd.PersistentFlags().StringVar(&listEntriesFrom, "from", start, "Start time of interval to list time entries from")
	listCmd.PersistentFlags().StringVar(&listEntriesTo, "to", end, "End time of interval to list time entries from")
}
