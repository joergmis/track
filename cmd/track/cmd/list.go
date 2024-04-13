package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
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
			activities, err := storage.GetActivities()
			if err != nil {
				log.Fatalf("get all activities from storage: %v", err)
			}

			start, err := time.Parse(clockodo.TimeLayoutString, listEntriesFrom)
			cobra.CheckErr(err)

			end, err := time.Parse(clockodo.TimeLayoutString, listEntriesTo)
			cobra.CheckErr(err)

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

			headings := []string{
				"start",
				"end",
				"duration",
				"customer",
				"project",
				"service",
				"description",
			}

			fmt.Fprintf(w, "%s\n", strings.Join(headings, "\t"))

			for _, entry := range activities {
				if entry.StartTime.Before(start) || entry.StartTime.After(end) || entry.EndTime.After(end) {
					continue
				}

				duration := ""
				if entry.EndTime.Sub(entry.StartTime) < 0 {
					duration = "in progress"
				} else {
					duration = fmt.Sprintf("%v", entry.EndTime.Sub(entry.StartTime))
				}

				fmt.Fprintf(w, "%v\t%v\t%v\t%s\t%s\t%s\t%s\n",
					entry.StartTime.Format(time.TimeOnly),
					entry.EndTime.Format(time.TimeOnly),
					duration,
					entry.Customer,
					entry.Project,
					entry.Service,
					entry.Description,
				)
			}

			w.Flush()
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
