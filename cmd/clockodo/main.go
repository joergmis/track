package main

import (
	"log"

	"github.com/joergmis/track/clockodo"
	"github.com/joergmis/track/ui"
	"github.com/spf13/viper"
)

func main() {
	repo, err := clockodo.NewRepository(clockodo.Config{
		EmailAddress: viper.GetString("clockodo.email"),
		ApiToken:     viper.GetString("clockodo.token"),
	})

	if err != nil {
		log.Fatal(err)
	}

	ui.New(repo).Run()
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
