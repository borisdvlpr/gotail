package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

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
