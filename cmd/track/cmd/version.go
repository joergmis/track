package cmd

import (
	"fmt"

	"github.com/joergmis/track"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the build version of the app",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(track.VersionString)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
