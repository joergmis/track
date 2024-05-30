package cmd

import (
	"log"

	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	var err error

	rootCmd.AddCommand(syncCmd)

	defaultBackend = track.BackendType(viper.GetString("backend.default"))

	switch defaultBackend {
	case track.BackendLocal:
		// nothing to do
		break

	case track.BackendClockodo:
		backend, err = clockodo.NewRepository(clockodo.Config{
			EmailAddress: viper.GetString("clockodo.email"),
			ApiToken:     viper.GetString("clockodo.token"),
		})
		if err != nil {
			log.Printf("setup clockodo repository: %v\n", err)
		}
		break

	default:
		log.Fatalf("backend %v does not match any known backend!\n", defaultBackend)
		break
	}

}
