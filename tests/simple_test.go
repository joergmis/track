package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	approvals "github.com/approvals/go-approval-tests"
	"github.com/joergmis/track"
)

func activityAsString(activity track.Activity) string {
	var buf bytes.Buffer

	activity.ID = "0"

	if err := json.NewEncoder(&buf).Encode(activity); err != nil {
		panic(err)
	}

	return buf.String()
}

func Test_Simple(t *testing.T) {
	approvals.UseFolder("data")
	r := approvals.UseReporter(NewNvimReporter())
	defer r.Close()

	customers := []string{
		"roger",
		"philip",
	}
	projects := []string{
		"painting the garage",
		"assemble the table",
	}
	services := []string{
		"mix the paint",
		"bring the hardware",
	}
	descriptions := []string{
		"get it at the store",
	}
	backends := []track.BackendType{
		track.BackendLocal,
		track.BackendClockodo,
	}

	i := 0

	for _, customer := range customers {
		for _, project := range projects {
			for _, service := range services {
				for _, description := range descriptions {
					for _, backend := range backends {
						t.Run(fmt.Sprintf("simple test %d", i), func(t *testing.T) {
							act := track.NewActivity(customer, project, service, description, backend)
							approvals.VerifyString(t, activityAsString(act))
						})
					}
				}
			}
		}
	}
}
