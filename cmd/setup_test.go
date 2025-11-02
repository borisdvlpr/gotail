package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	ierror "github.com/borisdvlpr/gotail/internal/error"
	"github.com/spf13/cobra"
)

type MockRootChecker struct {
	shouldError bool
}

func (m MockRootChecker) CheckRoot() error {
	errorMsg := "default check root error"

	if m.shouldError {
		return ierror.StatusError{Status: errorMsg, StatusCode: 1}
	}

	execPath := "fake/path/to/executable"
	args := []string{"sudo", execPath}
	args = append(args, os.Args[1:]...)

	return nil
}

func makeSetupCommand() (*cobra.Command, *bytes.Buffer) {
	testRootCmd := rootCmd
	var buf bytes.Buffer

	testRootCmd.SetOut(&buf)
	testRootCmd.SetErr(&buf)
	testRootCmd.SetArgs([]string{"setup"})

	return testRootCmd, &buf
}

func TestSetupCommandProperties(t *testing.T) {
	testSetupCmd := setupCmd
	cmdUse := "setup"
	cmdShort := "Setup Tailscale on a new device"

	if testSetupCmd.Use != cmdUse {
		t.Errorf("Expected Use to be '%s', got '%s'", cmdUse, testSetupCmd.Use)
	}

	if testSetupCmd.Short != cmdShort {
		t.Errorf("Expected Short to be '%s', got '%s'", cmdShort, testSetupCmd.Short)
	}
}

func TestSetupCommandNoRoot(t *testing.T) {
	testSetupCmd, buf := makeSetupCommand()

	originalChecker := rootChecker
	defer func() { rootChecker = originalChecker }()
	rootChecker = MockRootChecker{shouldError: true}

	err := testSetupCmd.Execute()
	if err == nil {
		t.Errorf("Expected error %v, got nil", err)
	}

	output := strings.TrimSpace(buf.String())
	if !strings.Contains(output, "Error: default check root error") {
		t.Errorf("Expected output to be 'Error: default check root error', got '%s'", output)
	}
}
