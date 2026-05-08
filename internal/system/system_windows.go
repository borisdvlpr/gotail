//go:build windows

package system

import (
	"fmt"

	ierror "github.com/borisdvlpr/gotail/internal/error"
	"golang.org/x/sys/windows"
)

// CheckRoot verifies the process is running with an elevated token (Run as
// Administrator). Auto-elevation via User Account Control (UAC) would require
// relaunching through ShellExecuteEx with the "runas" verb, which detaches
// the new process from the current console.
func (DefaultRootChecker) CheckRoot() error {
	var token windows.Token
	if err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token); err != nil {
		status := fmt.Sprintf("%s", err)
		return ierror.StatusError{Status: status, StatusCode: 1}
	}
	defer token.Close()

	if token.IsElevated() {
		return nil
	}

	return ierror.StatusError{
		Status:     "setup must be run as administrator. please relaunch from an elevated shell (run as administrator).",
		StatusCode: 126,
	}
}
