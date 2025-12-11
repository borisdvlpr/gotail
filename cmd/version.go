package cmd

import (
	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of your CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("gotail %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
