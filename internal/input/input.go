// Package input provides utility functions for user input and system checks.
// It includes functions to check for root privileges and prompt the user for input.
package input

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"syscall"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

// CheckRoot checks if the current user has root privileges.
// If not, it re-executes the script with sudo.
// If the sudo command is not found or fails, an error is return.
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

// PromptUser prompts the user with the given prompt string and reads the input from stdin.
// If allowedReplies is provided, it ensures the user's input is one of the allowed replies.
// It returns the user's input or an error if reading from stdin fails.
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

		if !slices.Contains(allowedReplies, answer) && len(allowedReplies) != 0 {
			return "", ierror.StatusError{Status: "abort", StatusCode: 1}
		}

		return answer, nil
	}
}
