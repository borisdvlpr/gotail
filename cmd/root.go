package cmd

import (
	"errors"
	"os"

	ierror "github.com/borisdvlpr/gotail/internal/error"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gotail",
	Short: "Bootstrap Tailscale into your Raspberry Pi",
	Long: `gotail is a CLI for bootstrapping Tailscale on a Raspberry Pi, 
automatically adding it to a tailnet from first boot.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		rootCmd.Printf("error: %s\n", err.Error())

		var statusErr ierror.StatusError
		if errors.As(err, &statusErr) {
			os.Exit(statusErr.StatusCode)
		}
		os.Exit(1)
	}
}
