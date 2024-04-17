package main

import (
	"bytes"
	"log"
	"os"
	"text/template"

	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const fileTemplate = `
// autogenerated file - do not edit!
// see 'track generate' command for more information
package cmd

import "github.com/joergmis/track"

func init() {
    services = []string{
        {{ range .Services }}"{{ . }}",{{ end }}
    }
    customerData = []track.Customer{
        {{ range .Customers }}{
            ID: "{{ .ID }}",
            Name: "\"{{ .Name | html }}\"",
            Projects: []track.Project{
                {{ range .Projects }}{
                        ID: "{{ .ID }}",
                        Name: "\"{{ .Name | html}}\"",
                    },{{ end }}
            },
        },{{ end }}
    }
}`

func main() {
	viper.SetConfigName("track")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("no config file found")
		}
	}

	if !viper.IsSet("clockodo.email") {
		log.Fatal("no email address set in configuration")
	}

	if !viper.IsSet("clockodo.token") {
		log.Fatal("no apikey set in configuration")
	}

	if !viper.IsSet("storage.dir") {
		log.Fatal("no storage filepath set in configuration")
	}

	backend, err := clockodo.NewRepository(clockodo.Config{
		EmailAddress: viper.GetString("clockodo.email"),
		ApiToken:     viper.GetString("clockodo.token"),
	})
	if err != nil {
		log.Fatalf("setup clockodo repository: %v", err)
	}

	customers, err := backend.GetAllCustomers()
	if err != nil {
		log.Fatal(err)
	}

	services, err := backend.GetAllServices()
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("data").Parse(fileTemplate)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	cobra.CheckErr(tmpl.Execute(&buf, struct {
		Customers []track.Customer
		Services  []string
	}{
		Customers: customers,
		Services:  services,
	}))

	if err := os.WriteFile("./cmd/track/cmd/data.go", buf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
}
