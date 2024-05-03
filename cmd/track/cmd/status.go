package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
)

var (
	statusCmdWeekly  bool
	statusCmdMonthly bool

	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Get an overview over the timetracking",
		Long: `By default, it shows the daily overview. Check out '--weekly' or
'--monthly' if you need to check out previous activities.`,
		Run: func(cmd *cobra.Command, args []string) {
			activities, err := storage.GetAllActivities()
			if err != nil {
				log.Fatalf("get all activities from storage: %v", err)
			}

			if len(activities) == 0 {
				fmt.Println("-- no activities in range --")
				return
			}

			var start, end time.Time

			if statusCmdWeekly {
				start = time.Now().Add(
					-time.Duration(time.Now().Hour()) -
						time.Duration(time.Now().Weekday()*24)*time.Hour)
			} else if statusCmdMonthly {
				start = time.Now().Add(
					-time.Duration(time.Now().Hour()) -
						time.Duration(time.Now().Day()*24)*time.Hour)
			} else {
				start = time.Now().Add(-1 * time.Duration(time.Now().Hour()) * time.Hour)
			}

			end = time.Now().Add(1 * time.Hour)

			t := tabby.New()
			t.AddHeader("time", "duration", "id", "customer", "project", "description")

			total := time.Duration(0)

			for i, entry := range activities {
				if entry.StartTime.Before(start) || entry.StartTime.After(end) || entry.EndTime.After(end) {
					continue
				}

				if i > 0 {
					previous := activities[i-1]

					if previous.StartTime.Weekday() == entry.StartTime.Weekday() && entry.StartTime.Sub(previous.EndTime).Minutes() > 5 {
						pause := entry.StartTime.Sub(previous.EndTime)
						t.AddLine(
							fmt.Sprintf("%s - %s", previous.EndTime.Add(1*time.Second).Format(time.TimeOnly), entry.StartTime.Add(-1*time.Second).Format(time.TimeOnly)),
							fmt.Sprintf("%02d:%02d h", int(pause.Hours()), int(pause.Minutes())%60),
							"-- pause --",
							"--",
							"--",
						)
					}
				}

				total += entry.Duration()

				t.AddLine(
					fmt.Sprintf(
						"%s - %s",
						entry.StartTime.Format(time.TimeOnly),
						entry.EndTime.Format(time.TimeOnly),
					),
					fmt.Sprintf(
						"%02d:%02d h",
						int(entry.Duration().Hours()),
						int(entry.Duration().Minutes())%60,
					),
					fmt.Sprintf("%d", i),
					entry.Customer,
					entry.Project,
					entry.Description,
				)
			}

			t.AddLine("---", "---")
			t.AddLine("total", fmt.Sprintf("%02d:%02d h", int(total.Hours()), int(total.Minutes())%60))

			t.Print()
		},
	}
)

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.PersistentFlags().BoolVar(&statusCmdWeekly, "weekly", false, "Show the weekly stats")
	statusCmd.PersistentFlags().BoolVar(&statusCmdMonthly, "monthly", false, "Show the monthly stats")
}
