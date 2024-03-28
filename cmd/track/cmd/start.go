package cmd

import (
	"fmt"

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
		fmt.Println("add called")
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

		return list, cobra.ShellCompDirectiveDefault
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
