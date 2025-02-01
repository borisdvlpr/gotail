package input

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"syscall"
)

func CheckRoot() error {
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

func PromptUser(prompt string, allowedReplies []string) (string, error) {
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
