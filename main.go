package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func handleError(err error) {
	fmt.Println("error:", err)
	os.Exit(1)
}

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

func main() {
	if err := checkRoot(); err != nil {
		handleError(err)
	}
}
