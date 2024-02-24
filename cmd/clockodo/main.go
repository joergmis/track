package main

import (
	"fmt"
	"log"

	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo"
	"github.com/spf13/viper"
)

func main() {
	var (
		config = clockodo.Config{
			EmailAddress: viper.GetString("clockodo.email"),
			ApiToken:     viper.GetString("clockodo.token"),
		}

		customerRepository track.CustomerRepository
		projectRepository  track.ProjectRepository

		err error
	)

	fmt.Printf("running version %s\n", track.Version)

	customerRepository, err = clockodo.NewCustomerRepository(config)
	if err != nil {
		log.Fatal(err)
	}

	projectRepository, err = clockodo.NewProjectRepository(config)
	if err != nil {
		log.Fatal(err)
	}

	customers, err := customerRepository.GetAllCustomers()
	if err != nil {
		log.Fatal(err)
	}
	for _, customer := range customers {
		log.Println(customer)
	}

	projects, err := projectRepository.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	for _, project := range projects {
		log.Println(project)
	}

	customer := track.Customer{Name: "example-customer"}

	projects, err = projectRepository.GetCustomerProjects(customer)
	if err != nil {
		log.Fatal(err)
	}
	for _, project := range projects {
		log.Printf("%s project: %s\n", customer.Name, project.Name)
	}

}

func init() {
	viper.SetConfigName("track")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/")
	viper.AddConfigPath(".")

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
}
