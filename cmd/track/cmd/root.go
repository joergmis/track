package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joergmis/track"
	"github.com/joergmis/track/clockodo"
	"github.com/joergmis/track/local"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	selectedBackend string

	backends = []track.ProjectRepository{}
	storage  track.ActivityRepository

	// only used for autocompletion!
	services     []string
	customerData []track.Customer

	rootCmd = &cobra.Command{
		Use:   "track",
		Short: "A timetracking application",
	}
)

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&selectedBackend, "backend", "b", string(backends[0].Type()), "Backend sets the what the activity should (eventually) be synced to")

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

	backends = append(backends, backend)

	storage, err = local.NewStorage(filepath.Join(os.Getenv("HOME"), ".config", viper.GetString("storage.dir")))
	if err != nil {
		log.Fatalf("setup storage repository: %v", err)
	}
}
