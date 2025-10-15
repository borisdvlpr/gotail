// Package input provides utility functions for user input. It includes
// functions to prompt the user for input and validate prompted inputs.
package input

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

// PromptUser prompts the user with the given prompt string and reads the input from stdin.
// If allowedReplies is provided, it ensures the user's input is one of the allowed replies.
// It returns the user's input or an error if reading from stdin fails.
func PromptUser(prompt string, allowedReplies []string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

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
