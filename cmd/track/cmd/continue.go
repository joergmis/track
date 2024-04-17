package cmd

import (
	"log"

	"github.com/joergmis/track"
	"github.com/spf13/cobra"
)

var continueCmd = &cobra.Command{
	Use:   "continue",
	Short: "Continue the last activity",
	Run: func(cmd *cobra.Command, args []string) {
		previousActivity, err := storage.GetLastActivity(track.ProjectBackendType(selectedBackend))
		if err != nil {
			log.Fatalf("get last activity: %v", err)
		}

		// clone the activity since we want an additional entry in the storage
		newActivity := track.NewActivity(previousActivity.Customer, previousActivity.Project, previousActivity.Service, previousActivity.Description)
		newActivity.Start()

		if err := storage.AddActivity(track.ProjectBackendType(selectedBackend), newActivity); err != nil {
			log.Fatalf("continue (restart) activity: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(continueCmd)
}
