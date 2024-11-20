package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"syscall"
)

func lsblkLinux() (map[string]interface{}, error) {
	lsblkCmd := exec.Command("lsblk", "--json")
	lsblkOut, err := lsblkCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	lsblk := make(map[string]interface{})
	if err = json.Unmarshal(lsblkOut, &lsblk); err != nil {
		return nil, fmt.Errorf("lsblk parsing: %w", err)
	}

	return lsblk, nil
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
}
