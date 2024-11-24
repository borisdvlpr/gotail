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
	for _, mountpoint := range mountpoints {
		if mountpoint != "" && !strings.Contains(mountpoint, "/snap") {
			filePath, err := getFilePath(mountpoint, "user-data")
			if err != nil {
				return "", fmt.Errorf("%w", err)
			}

			if filePath != "" {
				return filePath, nil
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
		fmt.Printf("%s [%s] ", prompt, strings.Join(allowedReplies, "/"))

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
	fmt.Println("error:", err)
	os.Exit(1)
}

func main() {
	if err := checkRoot(); err != nil {
		handleError(err)
	}

	filePath, err := findUserData()
	if err != nil {
		handleError(err)
	}

	if filePath == "" {
		fmt.Println("Could not find user-data file. Please try removing your SD card and re-inserting it.")
		os.Exit(1)
	}

	fmt.Printf("Found user-data file at %s.\n", filePath)
}
