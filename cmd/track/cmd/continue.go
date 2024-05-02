package cmd

import (
	"log"
	"strconv"

	"github.com/joergmis/track"
	"github.com/spf13/cobra"
)

var continueCmd = &cobra.Command{
	Use:   "continue",
	Short: "Continue an activity",
	Long: `By default, the last activity is used as a template for the new 
activity. Append the id to the command to continue a specific activity. 
`,
	Run: func(cmd *cobra.Command, args []string) {
		previousActivity, err := storage.GetLastActivity()
		if err != nil {
			log.Fatalf("get last activity: %v", err)
		}

		if previousActivity.InProgress() {
			previousActivity.Stop()
			if err := storage.UpdateActivity(previousActivity); err != nil {
				log.Fatalf("stop last activity: %v", err)
			}
		}

		var newActivity track.Activity

		if len(args) == 0 {
			previousActivity, err := storage.GetLastActivity()
			if err != nil {
				log.Fatalf("get last activity: %v", err)
			}

			newActivity = track.NewActivity(previousActivity.Customer, previousActivity.Project, previousActivity.Service, previousActivity.Description)
		} else if len(args) == 1 {
			index, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("parse index of activity to parse: %v", err)
			}

			activities, err := storage.GetAllActivities()
			if err != nil {
				log.Fatalf("get all activities: %v", err)
			}

			if index > len(activities)+1 {
				log.Fatalf("index is greater than the number of activities: %d > %d", index, len(activities))
			}

			previousActivity := activities[index]

			newActivity = track.NewActivity(previousActivity.Customer, previousActivity.Project, previousActivity.Service, previousActivity.Description)
		} else {
			log.Fatalf("unsupported number of arguments: %d", len(args))
		}

		newActivity.Start()

		if err := storage.AddActivity(newActivity); err != nil {
			log.Fatalf("continue (restart) activity: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(continueCmd)
}
