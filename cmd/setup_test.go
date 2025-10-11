package cmd

import "testing"

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
