package cmd

import (
	"bytes"
	"os"
	"text/template"

	"github.com/joergmis/track"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a list of clients used for autocompletion",
	Long: `The autocompletion needs a fixed list build-in. To achieve this,
the list has to be generated in advance before building the
application.`,
	Run: func(cmd *cobra.Command, args []string) {
		customers, err := repo.GetAllCustomers()
		cobra.CheckErr(err)

		tmpl, err := template.New("data").Parse(fileTemplate)
		cobra.CheckErr(err)

		var buf bytes.Buffer
		cobra.CheckErr(tmpl.Execute(&buf, struct {
			Customers []track.Customer
		}{
			Customers: customers,
		}))

		cobra.CheckErr(os.WriteFile("./cmd/track/cmd/data.go", buf.Bytes(), 0644))
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

const fileTemplate = `
// autogenerated file - do not edit!
// see 'generate' command for more information
package cmd

import "github.com/joergmis/track"

func init() {
    customerData = []track.Customer{
    {{ range .Customers }}
        {
            ID: "{{ .ID }}",
            Name: "\"{{ .Name | html }}\"",
            Projects: []track.Project{
                {{ range .Projects }}
                    {
                        ID: "{{ .ID }}",
                        Name: "\"{{ .Name | html}}\"",
                        Active: {{ .Active }},
                        Completed: {{ .Completed }},
                        Services: []track.Service{
                            {{ range .Services }}
                                {
                                    ID: "{{ .ID }}",
                                    Name: "\"{{ .Name | html }}\"",
                                },
                            {{ end }}
                        },
                    },
                {{ end }}
            },
        },
    {{ end }}
    }
}`
