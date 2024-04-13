package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/joergmis/track"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
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

		_ = track.Activity{
			CustomerID:  customer,
			ProjectID:   project,
			ServiceID:   service,
			Description: description,
			Start:       time.Now(),
		}

		dryrun, err := cmd.PersistentFlags().GetBool("dryrun")
		cobra.CheckErr(err)

		if dryrun {
			fmt.Printf("customer: %s\n", customer)
			fmt.Printf("project: %s\n", project)
			fmt.Printf("service: %s\n", service)
			fmt.Printf("description: %s\n", description)
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

func init() {
	rootCmd.AddCommand(startCmd)

	// TODO: switch to false / remove flag
	startCmd.PersistentFlags().Bool("dryrun", true, "dry-run - don't start actual activities")
}
