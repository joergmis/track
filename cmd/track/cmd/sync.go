package cmd

import (
	"log"

	"github.com/joergmis/track"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync all changed/new activities to the configured backend",
	Long:  `Note that this is (at least for now) a one-way process; local -> cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		activities, err := storage.GetActivities()
		if err != nil {
			log.Fatalf("get all activities: %v", err)
		}

		var backend track.ProjectRepository
		found := false
		for _, b := range backends {
			if b.Type() == track.ProjectBackendType(selectedBackend) {
				backend = b
				found = true
			}
		}

		if !found {
			log.Fatalf("no matching backend found for %v", selectedBackend)
		}

		for _, activity := range activities[track.ProjectBackendType(selectedBackend)] {
			if !activity.Synced && !activity.InProgress() {
				if err := backend.AddTimeEntry(activity); err != nil {
					log.Fatalf("sync activity: %v", err)
				}

				activity.Synced = true
				if err := storage.UpdateActivity(track.ProjectBackendType(selectedBackend), activity); err != nil {
					log.Fatalf("mark activitiy as synced: %v", err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
