package input

import (
	"errors"
	"fmt"
	"os"
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
	defaultStdin := os.Stdin
	defer func() { os.Stdin = defaultStdin }()

	testCases := []PromptUserTestCase{
		{id: "case_01", input: "y\n", allowedInputs: []string{"y", "n"}, expectedAnswer: "y", expectedError: nil},
		{id: "case_02", input: "tskey_test_1234_5678\n", allowedInputs: nil, expectedAnswer: "tskey_test_1234_5678", expectedError: nil},
		{id: "case_03", input: "asdf\n", allowedInputs: []string{"y", "n"}, expectedAnswer: "y", expectedError: ierror.StatusError{Status: "abort", StatusCode: 1}},
	}

	for _, tc := range testCases {
		r, w, _ := os.Pipe()
		os.Stdin = r

		go func() {
			_, err := w.WriteString(tc.input)
			if err != nil {
				t.Errorf("%v", err)
			}

			err = w.Close()
			if err != nil {
				t.Errorf("%v", err)
			}
		}()

		answer, err := PromptUser(tc.id, tc.allowedInputs)
		fmt.Println("")

		if !errors.Is(err, tc.expectedError) {
			t.Errorf("PromptUser() error = %v, wantErr %v", err, tc.expectedError)
		}

		if tc.expectedError == nil && answer != tc.expectedAnswer {
			t.Errorf("PromptUser() = %v, want %v", answer, tc.expectedAnswer)
		}
	}
}
