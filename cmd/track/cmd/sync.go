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
		activities, err := storage.GetUnsyncedActivities()
		if err != nil {
			log.Fatalf("get all unsynced activities: %v", err)
		}

		for _, activity := range activities {
			// only sync clockodo entries
			// TODO: what about other backends?
			if !activity.InProgress() && activity.Backend == track.BackendClockodo {
				if err := backend.AddTimeEntry(activity); err != nil {
					log.Fatalf("sync activity: %v", err)
				}

				if err := storage.MarkActivityAsSynced(activity); err != nil {
					log.Fatalf("mark activity as synced: %v", err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
