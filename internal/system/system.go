// Package system provides utility functions for system checks. It includes
// functions to check for root privileges and execute the program as root.
package system

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

// RootChecker defines the interface for verifying if the application is currently
// running with sufficient root or administrative privileges.
type RootChecker interface {
	CheckRoot() error
}

// DefaultRootChecker represents the primary implementation of the RootChecker interface
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
		status := fmt.Sprintf("%s", err)
		return ierror.StatusError{Status: status, StatusCode: 127}
	}

	execPath, err := os.Executable()
	if err != nil {
		status := fmt.Sprintf("failed to get executable path: %s", err)
		return ierror.StatusError{Status: status, StatusCode: 1}
	}

	args := []string{"sudo", execPath}
	args = append(args, os.Args[1:]...)

	err = syscall.Exec(sudoPath, args, os.Environ())
	if err != nil {
		status := fmt.Sprintf("failed to execute sudo: %s", err)
		return ierror.StatusError{Status: status, StatusCode: 126}
	}

	return nil
}
