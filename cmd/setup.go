package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/borisdvlpr/gotail/internal/config"
	ierror "github.com/borisdvlpr/gotail/internal/error"
	"github.com/borisdvlpr/gotail/internal/file"
	"github.com/borisdvlpr/gotail/internal/input"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var ConfigPath string

func handleError(err error) {
	if err != nil {
		fmt.Println("error:", err.Error())

		var statusErr ierror.StatusError
		if errors.As(err, &statusErr) {
			os.Exit(statusErr.StatusCode)
		}
		os.Exit(1)
	}
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup Tailscale on a new device",
	Run: func(cmd *cobra.Command, args []string) {
		flags := []string{"tailscale", "up", "--ssh"}
		initConfig := []string{
			"runcmd:\n",
			`  - [ sh, -c, curl -fsSL https://tailscale.com/install.sh | sh ]` + "\n",
		}
		config := &config.Config{}

		err := input.CheckRoot()
		handleError(err)

		filePath, err := file.FindUserData()
		handleError(err)

		fmt.Printf("Found 'user-data' file at '%s'.\n", filePath)

		confPath := viper.GetString("file")
		if confPath != "" {
			configFile, err := os.ReadFile(confPath)
			handleError(err)
			err = yaml.Unmarshal(configFile, &config)
			handleError(err)

		} else {
			config.ExitNode, err = input.PromptUser("Setup device as an exit node?", []string{"y", "n"})
			handleError(err)

			config.SubnetRouter, err = input.PromptUser("Setup device as a subnet router?", []string{"y", "n"})
			handleError(err)

			if config.SubnetRouter == "y" {
				config.Subnets, err = input.PromptUser("Please enter your subnets (comma separated):", nil)
				handleError(err)
			}

			config.Hostname, err = input.PromptUser("Please enter a hostname for this device:", nil)
			handleError(err)

			config.AuthKey, err = input.PromptUser("Please enter your Tailscale authkey:", nil)
			handleError(err)
		}

		err = config.Validate()
		handleError(err)

		if config.ExitNode == "y" {
			flags = append(flags, "--advertise-exit-node")
			fmt.Println("This device will be an exit node.")
		}

		if config.SubnetRouter == "y" && config.Subnets != "" {
			err = input.ValidateSubnets(config.Subnets)
			handleError(err)
			initConfig = append(initConfig, `  - [ sh, -c, echo 'net.ipv4.ip_forward = 1' | sudo tee -a /etc/sysctl.d/99-tailscale.conf && echo 'net.ipv6.conf.all.forwarding = 1' | sudo tee -a /etc/sysctl.d/99-tailscale.conf && sudo sysctl -p /etc/sysctl.d/99-tailscale.conf ]`+"\n")
			flags = append(flags, fmt.Sprintf("--advertise-routes=%s", config.Subnets))
		}

		if config.Hostname != "" {
			flags = append(flags, fmt.Sprintf("--hostname=%s", config.Hostname))
			initConfig = append(initConfig, fmt.Sprintf(`  - [ sh, -c, sudo hostnamectl hostname %s ]`+"\n", config.Hostname))
		}

		flags = append(flags, fmt.Sprintf("--authkey=%s", config.AuthKey))
		fmt.Println("Adding Tailscale to 'user-data' file.")

		initConfig = append(initConfig, fmt.Sprintf("  - [ %s ]\n", strings.Join(flags, ", ")))

		initFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		handleError(err)

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				handleError(err)
			}
		}(initFile)

		writer := bufio.NewWriter(initFile)
		for _, conf := range initConfig {
			_, err = writer.WriteString(conf)
			handleError(err)
		}

		err = writer.Flush()
		handleError(err)
		fmt.Println("Tailscale will be installed on boot. Please eject your SD card and boot your Raspberry Pi.")
	},
}

func init() {
	setupCmd.PersistentFlags().StringVarP(&ConfigPath, "file", "f", "", "path to the config file")
	viper.BindPFlag("file", setupCmd.PersistentFlags().Lookup("file"))
	rootCmd.AddCommand(setupCmd)
}
