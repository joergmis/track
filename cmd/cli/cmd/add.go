package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		list := []string{}

		if len(args) == 0 {
			for _, customer := range customerData {
				list = append(list, customer.ID)
			}
		}

		if len(args) == 1 {
			for _, c := range customerData {
				if args[0] == c.ID {
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
	rootCmd.AddCommand(addCmd)
}
