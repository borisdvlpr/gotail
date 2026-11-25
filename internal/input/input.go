// Package input provides utility functions for user input. It includes
// functions to prompt the user for input and validate prompted inputs.
package input

import (
	"bufio"
	"fmt"
	"slices"
	"strings"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

// PromptUser prompts the user with the given prompt string and reads one line from reader.
// Callers asking several questions must share one reader, since it buffers input past the
// current line. If allowedReplies is provided, it ensures the user's input is one of the
// allowed replies. It returns the user's input or an error if reading fails.
func PromptUser(reader *bufio.Reader, prompt string, allowedReplies []string) (string, error) {
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
