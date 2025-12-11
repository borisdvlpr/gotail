package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func makeVersionCommand() (*cobra.Command, *bytes.Buffer) {
	testRootCmd := rootCmd
	var buf bytes.Buffer

	testRootCmd.SetOut(&buf)
	testRootCmd.SetErr(&buf)
	testRootCmd.SetArgs([]string{"version"})

	return testRootCmd, &buf
}

func TestVersion_CommandProperties(t *testing.T) {
	testVersionCmd := versionCmd
	cmdUse := "version"
	cmdShort := "Show the version of your CLI tool"

	if testVersionCmd.Use != cmdUse {
		t.Errorf("Expected Use to be '%s', got '%s'", cmdUse, testVersionCmd.Use)
	}

	if testVersionCmd.Short != cmdShort {
		t.Errorf("Expected Short to be '%s', got '%s'", cmdShort, testVersionCmd.Short)
	}
}

func TestVersion_CommandOutput(t *testing.T) {
	testRootCmd, buf := makeVersionCommand()

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

func TestVersion_CommandWithCustomVersion(t *testing.T) {
	testRootCmd, buf := makeVersionCommand()

	originalVersion := version
	defer func() { version = originalVersion }()
	version = "1.2.3"

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

func TestVersion_CommandFormat(t *testing.T) {
	testRootCmd, buf := makeVersionCommand()

	err := testRootCmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	output := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(output, "gotail ") {
		t.Errorf("Expected output to start with 'gotail ', got '%s'", output)
	}

	parts := strings.Split(output, " ")
	if len(parts) != 2 {
		t.Errorf("Expected output format 'gotail <version>', got '%s'", output)
	}
}
