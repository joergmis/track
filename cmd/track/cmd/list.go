package cmd

import (
	"fmt"
	"time"

	"github.com/joergmis/track/clockodo"
	"github.com/spf13/cobra"
)

var (
	listEntriesFrom string
	listEntriesTo   string

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List existing time entries",
		Run: func(cmd *cobra.Command, args []string) {
			start, err := time.Parse(clockodo.TimeLayoutString, listEntriesFrom)
			cobra.CheckErr(err)

			end, err := time.Parse(clockodo.TimeLayoutString, listEntriesTo)
			cobra.CheckErr(err)

			entries, err := repo.GetTimeEntries(start, end)
			cobra.CheckErr(err)

			for _, entry := range entries {
				fmt.Println(entry.Description)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)

	start := time.Now().Add(-(time.Hour * 24)).Format(clockodo.TimeLayoutString)
	end := time.Now().Format(clockodo.TimeLayoutString)

	listCmd.PersistentFlags().StringVar(&listEntriesFrom, "from", start, "Start time of interval to list time entries from")
	listCmd.PersistentFlags().StringVar(&listEntriesTo, "to", end, "End time of interval to list time entries from")
}
