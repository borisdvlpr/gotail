package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func lsblkLinux() (map[string]interface{}, error) {
	lsblkCmd := exec.Command("lsblk", "--json")
	lsblkOut, err := lsblkCmd.Output()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	lsblk := make(map[string]interface{})
	err = json.Unmarshal(lsblkOut, &lsblk)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	return lsblk, nil
}

func checkRoot() {
	if os.Geteuid() != 0 {
		fmt.Println("This script must be run as root. Re-executing with sudo...")

		sudoPath, err := exec.LookPath("sudo")
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		args := append([]string{"sudo"}, os.Args...)
		err = syscall.Exec(sudoPath, args, os.Environ())
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}
}

func main() {
	checkRoot()
}
