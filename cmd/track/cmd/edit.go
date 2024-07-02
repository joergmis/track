package cmd

import (
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the existing time entries. Uses $EDITOR and alacritty.",
	Run: func(cmd *cobra.Command, args []string) {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		home := os.Getenv("HOME")

		err := exec.Command("alacritty", "-e", editor, path.Join(home, ".config/track/entries.json")).Start()
		if err != nil {
			log.Fatal("failed to launch editor")
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
