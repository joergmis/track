package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a running time entry",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stop called")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
