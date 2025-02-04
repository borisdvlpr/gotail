package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	ierror "github.com/borisdvlpr/gotail/internal/error"
	"github.com/borisdvlpr/gotail/internal/file"
	"github.com/borisdvlpr/gotail/internal/input"
)

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

func main() {
	flags := []string{"tailscale", "up", "--ssh"}
	configs := []string{
		"runcmd:\n",
		`  - [ sh, -c, curl -fsSL https://tailscale.com/install.sh | sh ]` + "\n",
		`  - [ sh, -c, echo 'net.ipv4.ip_forward = 1' | sudo tee -a /etc/sysctl.d/99-tailscale.conf && echo 'net.ipv6.conf.all.forwarding = 1' | sudo tee -a /etc/sysctl.d/99-tailscale.conf && sudo sysctl -p /etc/sysctl.d/99-tailscale.conf ]` + "\n",
	}

	err := input.CheckRoot()
	handleError(err)

	filePath, err := file.FindUserData()
	handleError(err)

	fmt.Printf("Found 'user-data' file at '%s'.\n", filePath)

	exitNode, err := input.PromptUser("Would you like this device to be an exit node?", []string{"y", "n"})
	handleError(err)

	if exitNode == "y" {
		flags = append(flags, "--advertise-exit-node")
		fmt.Println("This device will be an exit node.")
	}

	authKey, err := input.PromptUser("Please enter your Tailscale authkey:", nil)
	handleError(err)
	flags = append(flags, fmt.Sprintf("--authkey=%s", authKey))

	hostName, err := input.PromptUser("Please enter a hostname for this device:", nil)
	handleError(err)

	if hostName != "" {
		flags = append(flags, fmt.Sprintf("--hostname=%s", hostName))
		configs = append(configs, fmt.Sprintf(`  - [ sh, -c, sudo hostnamectl hostname %s ]`+"\n", hostName))
	}
	fmt.Println("Adding Tailscale to 'user-data' file.")

	configs = append(configs, fmt.Sprintf("  - [ %s ]\n", strings.Join(flags, ", ")))

	initFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	handleError(err)

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			handleError(err)
		}
	}(initFile)

	writer := bufio.NewWriter(initFile)
	for _, config := range configs {
		_, err = writer.WriteString(config)
		handleError(err)
	}

	err = writer.Flush()
	handleError(err)
	fmt.Println("Tailscale will be installed on boot. Please eject your SD card and boot your Raspberry Pi.")
}
