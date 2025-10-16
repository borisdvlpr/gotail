// Package system provides utility functions for system checks. It includes
// functions to check for root privileges and execute the program as root.
package system

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type RootChecker interface {
	CheckRoot() error
}

type DefaultRootChecker struct{}

// CheckRoot checks if the current user has root privileges.
// If not, it re-executes the script with sudo.
// If the sudo command is not found or fails, an error is return.
func (DefaultRootChecker) CheckRoot() error {
	if os.Geteuid() == 0 {
		return nil
	}

	fmt.Println("Setup must be run as root. Re-executing with sudo...")

	sudoPath, err := exec.LookPath("sudo")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	args := []string{"sudo", execPath}
	args = append(args, os.Args[1:]...)

	err = syscall.Exec(sudoPath, args, os.Environ())
	if err != nil {
		return fmt.Errorf("failed to execute sudo: %w", err)
	}

	return nil
}
