package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/borisdvlpr/gotail/internal/config"
	ierror "github.com/borisdvlpr/gotail/internal/error"
	"github.com/borisdvlpr/gotail/internal/file"
	"github.com/borisdvlpr/gotail/internal/input"
	"github.com/borisdvlpr/gotail/internal/system"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var ConfigPath string
var rootChecker system.RootChecker = system.DefaultRootChecker{}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup Tailscale on a new device",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		flags := []string{"tailscale", "up", "--ssh"}
		initConfig := []string{
			"runcmd:\n",
			`  - [ sh, -c, curl -fsSL https://tailscale.com/install.sh | sh ]` + "\n",
		}
		config := &config.Config{}

		if err := rootChecker.CheckRoot(); err != nil {
			return err
		}

		filePath, err := file.FindUserData()
		if err != nil {
			return err
		}

		cmd.Printf("Found 'user-data' file at '%s'.\n", filePath)

		confPath := viper.GetString("file")
		if confPath != "" {
			configFile, err := os.ReadFile(confPath)
			if err != nil {
				return err
			}

			if err = yaml.Unmarshal(configFile, &config); err != nil {
				return err
			}

		} else {
			config.ExitNode, err = input.PromptUser("Setup device as an exit node?", []string{"y", "n"})
			if err != nil {
				return err
			}

			config.SubnetRouter, err = input.PromptUser("Setup device as a subnet router?", []string{"y", "n"})
			if err != nil {
				return err
			}

			if config.SubnetRouter == "y" {
				config.Subnets, err = input.PromptUser("Please enter your subnets (comma separated):", nil)
				if err != nil {
					return err
				}
			}

			config.Hostname, err = input.PromptUser("Please enter a hostname for this device:", nil)
			if err != nil {
				return err
			}

			config.AuthKey, err = input.PromptUser("Please enter your Tailscale authkey:", nil)
			if err != nil {
				return err
			}
		}

		if err = config.Validate(); err != nil {
			return err
		}

		if config.ExitNode == "y" {
			flags = append(flags, "--advertise-exit-node")
			cmd.Println("This device will be an exit node.")
		}

		if config.SubnetRouter == "y" && config.Subnets != "" {
			if err = input.ValidateSubnets(config.Subnets); err != nil {
				return err
			}

			initConfig = append(initConfig, `  - [ sh, -c, echo 'net.ipv4.ip_forward = 1' | sudo tee -a /etc/sysctl.d/99-tailscale.conf && echo 'net.ipv6.conf.all.forwarding = 1' | sudo tee -a /etc/sysctl.d/99-tailscale.conf && sudo sysctl -p /etc/sysctl.d/99-tailscale.conf ]`+"\n")
			flags = append(flags, fmt.Sprintf("--advertise-routes=%s", config.Subnets))
		}

		if config.Hostname != "" {
			flags = append(flags, fmt.Sprintf("--hostname=%s", config.Hostname))
			initConfig = append(initConfig, fmt.Sprintf(`  - [ sh, -c, sudo hostnamectl hostname %s ]`+"\n", config.Hostname))
		}

		flags = append(flags, fmt.Sprintf("--authkey=%s", config.AuthKey))
		cmd.Println("Adding Tailscale to 'user-data' file.")

		initConfig = append(initConfig, fmt.Sprintf("  - [ %s ]\n", strings.Join(flags, ", ")))

		initFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		defer func() {
			if closeErr := initFile.Close(); closeErr != nil && err == nil {
				status := fmt.Sprintf("failed to close file: %s", closeErr)
				err = ierror.StatusError{Status: status, StatusCode: 74} // EX_IOERR
			}
		}()

		writer := bufio.NewWriter(initFile)
		for _, conf := range initConfig {
			_, err = writer.WriteString(conf)
			if err != nil {
				return err
			}
		}

		if err = writer.Flush(); err != nil {
			return err
		}

		cmd.Println("Tailscale will be installed on boot. Please eject your SD card and boot your Raspberry Pi.")
		return nil
	},
}

func init() {
	setupCmd.PersistentFlags().StringVarP(&ConfigPath, "file", "f", "", "path to the config file")
	viper.BindPFlag("file", setupCmd.PersistentFlags().Lookup("file"))
	rootCmd.AddCommand(setupCmd)
}
