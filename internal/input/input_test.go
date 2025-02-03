package input

import (
	"fmt"
	"os"
	"testing"
)

type PromptUserTestCase struct {
	id             string
	input          string
	allowedInputs  []string
	expectedAnswer string
}

func TestPromptUser(t *testing.T) {
	defaultStdin := os.Stdin
	defer func() { os.Stdin = defaultStdin }()

	testCases := []PromptUserTestCase{
		{"case01", "y\n", []string{"y", "n"}, "y"},
		{"case02", "tskey_test_1234_5678\n", nil, "tskey_test_1234_5678"},
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

		if err != nil {
			t.Errorf("promptUser() error = %v, wantErr %v", err, nil)
		}

		if answer != tc.expectedAnswer {
			t.Errorf("promptUser() = %v, want %v", answer, tc.expectedAnswer)
		}
	}
}
