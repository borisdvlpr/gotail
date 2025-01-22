package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"syscall"
)

type LsblkOutput struct {
	Blockdevices []struct {
		Name        string   `json:"name"`
		MajMin      string   `json:"maj:min"`
		Rm          bool     `json:"rm"`
		Size        string   `json:"size"`
		Ro          bool     `json:"ro"`
		Type        string   `json:"type"`
		Mountpoints []string `json:"mountpoints"`
		Children    []struct {
			Name        string   `json:"name"`
			MajMin      string   `json:"maj:min"`
			Rm          bool     `json:"rm"`
			Size        string   `json:"size"`
			Ro          bool     `json:"ro"`
			Type        string   `json:"type"`
			Mountpoints []string `json:"mountpoints"`
		} `json:"children,omitempty"`
	} `json:"blockdevices"`
}

func lsblk() (LsblkOutput, error) {
	lsblkCmd := exec.Command("lsblk", "--json")
	lsblkOut, err := lsblkCmd.Output()
	if err != nil {
		return LsblkOutput{}, fmt.Errorf("%w", err)
	}

	var lsblk LsblkOutput
	if err = json.Unmarshal(lsblkOut, &lsblk); err != nil {
		return LsblkOutput{}, fmt.Errorf("lsblk parsing: %w", err)
	}

	return lsblk, nil
}

func getFilePath(rootDir string, fileName string) (string, error) {
	var filePath string

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		if !d.IsDir() && d.Name() == fileName {
			filePath = path
			return fs.SkipDir
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return filePath, nil
}

func searchMountpoints(mountpoints []string) (string, error) {
	ignorePaths := []string{"/boot", "/home", "/snap"}

	for _, mountpoint := range mountpoints {
		if mountpoint != "" {
			validPath := !slices.ContainsFunc(ignorePaths, func(s string) bool {
				return strings.HasPrefix(mountpoint, s)
			})

			if mountpoint != "/" && validPath {
				filePath, err := getFilePath(mountpoint, "user-data")
				if err != nil {
					return "", fmt.Errorf("%w", err)
				}

				if filePath != "" {
					return filePath, nil
				}
			}
		}
	}

	return "", nil
}

func findUserData() (string, error) {
	var filePath string
	var err error

	if runtime.GOOS == "darwin" {
		filePath, err = getFilePath("/Volumes", "user-data")
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		if filePath != "" {
			return filePath, nil
		}
	}

	if runtime.GOOS == "linux" {
		devices, err := lsblk()
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		for _, device := range devices.Blockdevices {
			if device.Mountpoints != nil {
				filePath, err = searchMountpoints(device.Mountpoints)
				if err != nil {
					return "", fmt.Errorf("%w", err)
				}

				if filePath != "" {
					return filePath, nil
				}
			}

			if device.Children != nil {
				for _, child := range device.Children {
					if child.Mountpoints != nil {
						filePath, err = searchMountpoints(child.Mountpoints)
						if err != nil {
							return "", fmt.Errorf("%w", err)
						}

						if filePath != "" {
							return filePath, nil
						}
					}
				}
			}
		}
	}

	return "", nil
}

func promptUser(prompt string, allowedReplies []string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		if allowedReplies != nil {
			fmt.Printf("%s [%s] ", prompt, strings.Join(allowedReplies, "/"))
		} else {
			fmt.Printf("%s ", prompt)
		}

		answer, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		answer = strings.TrimSpace(answer)

		if slices.Contains(allowedReplies, answer) || len(allowedReplies) == 0 {
			return answer, nil
		}

		fmt.Println("Option not available. Please try again.")
	}
}

func checkRoot() error {
	if os.Geteuid() == 0 {
		return nil
	}

	fmt.Println("This script must be run as root. Re-executing with sudo...")

	sudoPath, err := exec.LookPath("sudo")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	args := append([]string{"sudo"}, os.Args...)
	err = syscall.Exec(sudoPath, args, os.Environ())
	if err != nil {
		return fmt.Errorf("failed to execute sudo: %w", err)
	}

	return nil
}

func handleError(err error) {
	if err != nil {
		fmt.Println("error:", err)
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

	err := checkRoot()
	handleError(err)

	filePath, err := findUserData()
	handleError(err)

	if filePath == "" {
		fmt.Println("Could not find 'user-data' file. Please try removing your SD card and re-inserting it.")
		os.Exit(1)
	}

	fmt.Printf("Found 'user-data' file at '%s'.\n", filePath)

	exitNode, err := promptUser("Would you like this device to be an exit node?", []string{"y", "n"})
	handleError(err)

	if exitNode == "y" {
		flags = append(flags, "--advertise-exit-node")
		fmt.Println("This device will be an exit node.")
	}

	authKey, err := promptUser("Please enter your Tailscale authkey:", nil)
	handleError(err)
	flags = append(flags, fmt.Sprintf("--authkey=%s", authKey))

	hostName, err := promptUser("Please enter a hostname for this device:", nil)
	handleError(err)

	if hostName != "" {
		flags = append(flags, fmt.Sprintf("--hostname=%s", hostName))
		configs = append(configs, fmt.Sprintf(`  - [ sh, -c, sudo hostnamectl hostname %s ]`+"\n", hostName))
	}
	fmt.Println("Adding Tailscale to 'user-data' file.")

	configs = append(configs, fmt.Sprintf("  - [ %s ]\n", strings.Join(flags, ", ")))

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	handleError(err)

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			handleError(err)
		}
	}(file)

	writer := bufio.NewWriter(file)
	for _, config := range configs {
		_, err = writer.WriteString(config)
		handleError(err)
	}

	err = writer.Flush()
	handleError(err)
	fmt.Println("Tailscale will be installed on boot. Please eject your SD card and boot your Raspberry Pi.")
}
