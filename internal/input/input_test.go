package input

import (
	"os"
	"testing"
)

func TestPromptUser(t *testing.T) {
	defaultStdin := os.Stdin
	defer func() { os.Stdin = defaultStdin }()

	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		_, err := w.WriteString("y\n")
		if err != nil {
			t.Errorf("%v", err)
		}

		err = w.Close()
		if err != nil {
			t.Errorf("%v", err)
		}
	}()

	answer, err := PromptUser("Would you like this device to be an exit node?", []string{"y", "n"})
	if err != nil {
		t.Errorf("promptUser() error = %v, wantErr %v", err, nil)
	}

	if answer != "y" {
		t.Errorf("promptUser() = %v, want %v", answer, "y")
	}
}
