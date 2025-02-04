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
		{"case01", "y\n", []string{"y", "n"}, "y", nil},
		{"case02", "tskey_test_1234_5678\n", nil, "tskey_test_1234_5678", nil},
		{"case03", "asdf\n", []string{"y", "n"}, "y", ierror.StatusError{Status: "abort", StatusCode: 1}},
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
			t.Errorf("promptUser() error = %v, wantErr %v", err, tc.expectedError)
		}

		if tc.expectedError == nil && answer != tc.expectedAnswer {
			t.Errorf("promptUser() = %v, want %v", answer, tc.expectedAnswer)
		}
	}
}
