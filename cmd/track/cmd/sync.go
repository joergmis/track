package cmd

import (
	"log"

	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync all changed/new activities to the configured backend",
	Long:  `Note that this is (at least for now) a one-way process; local -> cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := setupRepository(); err != nil {
			log.Fatalf("setup backend: %v", err)
		}

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

func setupRepository() error {
	defaultBackend = track.BackendType(viper.GetString("backend.default"))

	switch defaultBackend {
	case track.BackendLocal:
		// nothing to do

	case track.BackendClockodo:
		var err error
		backend, err = clockodo.NewRepository(clockodo.Config{
			EmailAddress: viper.GetString("clockodo.email"),
			ApiToken:     viper.GetString("clockodo.token"),
		})
		return errors.Wrap(err, "setup clockodo repository")

	default:
		return errors.Wrapf(track.ErrNoMatchingBackend, "backend %v does not match any known backend!\n", defaultBackend)
	}

	return nil
}
