package cmd

import (
	"log"
	"strings"

	"github.com/joergmis/track"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	startBackendType string

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start (and stop a previous) time entry",
		Long: `Use the format:

./track start <client> <project> <service> <description>

'generate' and 'completion' can be used to generate the 
autocompletion for your prefered shell.`,
		Run: func(cmd *cobra.Command, args []string) {
			customer := args[0]
			project := args[1]
			service := args[2]
			description := strings.Join(args[3:], " ")

			previousActivity, err := storage.GetLastActivity()
			if err != nil {
				if errors.Cause(err) != track.ErrNoActivities {
					log.Fatalf("get last activity: %v", err)
				}
			} else {
				previousActivity.Stop()
				if err := storage.UpdateActivity(previousActivity); err != nil {
					log.Fatalf("stop the previous activity: %v", err)
				}
			}

			if !track.BackendType(startBackendType).Valid() {
				log.Fatalf("invalid backend selected: %v", startBackendType)
			}

			newActivity := track.NewActivity(
				customer,
				project,
				service,
				description,
				track.BackendType(startBackendType),
			)
			newActivity.Start()

			if err := storage.AddActivity(newActivity); err != nil {
				log.Fatalf("start new activity: %v", err)
			}
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// used for the autocompletion feature - this only works if the data
			// was generated in advance with the `generate` command
			list := []string{}

			if len(args) == 0 {
				for _, customer := range customerData {
					list = append(list, customer.ID)
				}
			}

			if len(args) == 1 {
				customer := args[0]

				for _, c := range customerData {
					if customer == c.ID {
						for _, project := range c.Projects {
							list = append(list, project.ID)
						}
					}
				}
			}

			if len(args) == 2 {
				list = services
			}

			return list, cobra.ShellCompDirectiveDefault
		},
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringVarP(&startBackendType, "backend", "b", string(defaultBackend), "Select backend to assign to the activity")
}
