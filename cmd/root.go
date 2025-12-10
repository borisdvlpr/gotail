package cmd

import (
	"errors"
	"os"

	ierror "github.com/borisdvlpr/gotail/internal/error"
	"github.com/borisdvlpr/gotail/internal/file"
	"github.com/borisdvlpr/gotail/internal/system"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gotail",
	Short: "Bootstrap Tailscale into your Raspberry Pi",
	Long: `gotail is a CLI for bootstrapping Tailscale on a Raspberry Pi,
automatically adding it to a tailnet from first boot.`,
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	fs := afero.NewOsFs()

	setupDeps := SetupCommand{
		Fsys:        fs,
		RootChecker: system.DefaultRootChecker{},
		SystemSearcher: &file.SystemSearcher{
			Fsys:         fs,
			DeviceLister: &file.DefaultBlockDeviceLister{},
		},
	}

	setupCmd := NewSetupCmd(setupDeps)

	rootCmd.AddCommand(setupCmd)

	if err := rootCmd.Execute(); err != nil {
		var statusErr ierror.StatusError
		if errors.As(err, &statusErr) {
			os.Exit(statusErr.StatusCode)
		}
		os.Exit(1)
	}
}
