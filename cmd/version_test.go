package cmd

import (
	"testing"
)

func TestVersionCommandProperties(t *testing.T) {
	testVersionCmd := versionCmd
	cmdUse := "version"
	cmdShort := "Show the version of your CLI tool"

	if testVersionCmd.Use != cmdUse {
		t.Errorf("Expected Use to be '%s', got '%s'", cmdUse, versionCmd.Use)
	}

	if testVersionCmd.Short != cmdShort {
		t.Errorf("Expected Short to be '%s', got '%s'", cmdShort, versionCmd.Short)
	}
}
