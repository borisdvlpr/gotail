package input

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

type PromptUserTestCase struct {
	id             string
	input          string
	allowedInputs  []string
	expectedAnswer string
	expectedError  error
}

func TestPromptUser(t *testing.T) {
	testCases := []PromptUserTestCase{
		{
			id:             "valid_allowed_reply",
			input:          "y\n",
			allowedInputs:  []string{"y", "n"},
			expectedAnswer: "y",
			expectedError:  nil,
		},
		{
			id:             "freeform_input_accepted",
			input:          "tskey_test_1234_5678\n",
			allowedInputs:  nil,
			expectedAnswer: "tskey_test_1234_5678",
			expectedError:  nil,
		},
		{
			id:            "disallowed_reply_aborts",
			input:         "asdf\n",
			allowedInputs: []string{"y", "n"},
			expectedError: ierror.StatusError{Status: "abort", StatusCode: 1},
		},
		{
			id:            "read_error_on_eof",
			input:         "",
			allowedInputs: nil,
			expectedError: io.EOF,
		},
	}

	for _, tc := range testCases {
		reader := bufio.NewReader(strings.NewReader(tc.input))

		answer, err := PromptUser(reader, tc.id, tc.allowedInputs)
		fmt.Println("")

		if !errors.Is(err, tc.expectedError) {
			t.Errorf("%v: PromptUser() error = %q, wantErr %q", tc.id, err, tc.expectedError)
		}

		if tc.expectedError == nil && answer != tc.expectedAnswer {
			t.Errorf("%v: PromptUser() = %v, want %v", tc.id, answer, tc.expectedAnswer)
		}
	}
}

// TestPromptUser_SharedReader checks that answers queued up front (as when
// input is piped in) are all delivered, one per prompt, in order.
func TestPromptUser_SharedReader(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("y\nn\nraspberrypi\n"))

	prompts := []struct {
		allowed  []string
		expected string
	}{
		{[]string{"y", "n"}, "y"},
		{[]string{"y", "n"}, "n"},
		{nil, "raspberrypi"},
	}

	for i, p := range prompts {
		answer, err := PromptUser(reader, fmt.Sprintf("prompt %d", i), p.allowed)
		fmt.Println("")

		if err != nil {
			t.Fatalf("prompt %d: PromptUser() returned error %v, expected none", i, err)
		}

		if answer != p.expected {
			t.Errorf("prompt %d: PromptUser() = %q, want %q", i, answer, p.expected)
		}
	}
}
