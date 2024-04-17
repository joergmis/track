package cmd

import (
	"fmt"
	"log"

	"github.com/joergmis/track"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a running time entry",
	Run: func(cmd *cobra.Command, args []string) {
		activity, err := storage.GetLastActivity(track.ProjectBackendType(selectedBackend))
		if err != nil {
			log.Fatalf("get last activity: %v", err)
		}

		if activity.InProgress() {
			activity.Stop()
			if err := storage.UpdateActivity(track.ProjectBackendType(selectedBackend), activity); err != nil {
				log.Fatalf("stop activity: %v", err)
			}
		} else {
			fmt.Println("activity was already stopped... (?)")
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
