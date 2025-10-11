package cmd

import (
	"bytes"
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

func TestVersionCommandOutput(t *testing.T) {
	testRootCmd := rootCmd
	var buf bytes.Buffer

	testRootCmd.SetOut(&buf)
	testRootCmd.SetErr(&buf)

	testRootCmd.SetArgs([]string{"version"})
	err := testRootCmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	output := buf.String()
	expectedOutput := "gotail dev\n"

	if output != expectedOutput {
		t.Errorf("Expected output '%s', got '%s'", expectedOutput, output)
	}
}

func TestVersionCommandWithCustomVersion(t *testing.T) {
	testRootCmd := rootCmd
	var buf bytes.Buffer
	originalVersion := version
	defer func() { version = originalVersion }()
	version = "1.2.3"

	testRootCmd.SetOut(&buf)
	testRootCmd.SetErr(&buf)

	testRootCmd.SetArgs([]string{"version"})
	err := testRootCmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	output := buf.String()
	expectedOutput := "gotail 1.2.3\n"

	if output != expectedOutput {
		t.Errorf("Expected output '%s', got '%s'", expectedOutput, output)
	}
}
