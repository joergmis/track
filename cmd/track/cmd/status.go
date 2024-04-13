package cmd

import (
	"fmt"
	"log"

	"github.com/joergmis/track/clockodo"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get an overview over the currently active time tracking",
	Run: func(cmd *cobra.Command, args []string) {
		activity, err := storage.GetLastActivity()
		if err != nil {
			log.Fatalf("get last activity: %v", err)
		}

		fmt.Println("currently tracking:")
		fmt.Println("")
		fmt.Printf("customer:    %s\n", activity.Customer)
		fmt.Printf("project:     %s\n", activity.Project)
		fmt.Printf("service:     %s\n", activity.Service)
		fmt.Printf("description: %s\n", activity.Description)
		fmt.Printf("start:       %s\n", activity.StartTime.Format(clockodo.TimeLayoutString))
		if activity.InProgress {
			fmt.Println("end:        -- ongoing --")
		} else {
			fmt.Printf("end:         %s\n", activity.EndTime.Format(clockodo.TimeLayoutString))
		}
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
