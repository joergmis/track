package cmd

import (
	"log"
	"os"

	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	repo         track.Repository
	customerData []track.Customer

	rootCmd = &cobra.Command{
		Use:   "track",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var err error

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

	repo, err = clockodo.NewRepository(clockodo.Config{
		EmailAddress: viper.GetString("clockodo.email"),
		ApiToken:     viper.GetString("clockodo.token"),
	})
	if err != nil {
		log.Fatalf("setup clockodo repository: %v", err)
	}
}