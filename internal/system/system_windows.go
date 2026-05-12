//go:build windows

package system

import (
	"fmt"

	ierror "github.com/borisdvlpr/gotail/internal/error"
	"golang.org/x/sys/windows"
)

// tokenElevationChecker defines the interface for checking whether the current
// process token carries an elevated privilege level.
type tokenElevationChecker interface {
	isElevated() (bool, error)
}

// windowsTokenChecker is the production implementation of tokenElevationChecker.
// It queries the real Windows process token via the Win32 API.
type windowsTokenChecker struct{}

// isElevated opens the current process token and reports whether it is elevated.
func (windowsTokenChecker) isElevated() (bool, error) {
	var token windows.Token
	if err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token); err != nil {
		return false, err
	}
	defer token.Close()
	return token.IsElevated(), nil
}

// CheckRoot verifies the process is running with an elevated token (Run as
// Administrator). Auto-elevation via User Account Control (UAC) would require
// relaunching through ShellExecuteEx with the "runas" verb, which detaches
// the new process from the current console.
func (d DefaultRootChecker) CheckRoot() error {
	return checkRootWithChecker(windowsTokenChecker{})
}

// checkRootWithChecker uses the provided tokenElevationChecker to determine
// whether the process has administrative privileges. It returns a StatusError
// with code 1 if the token cannot be queried, or code 126 if the process is
// not elevated.
func checkRootWithChecker(checker tokenElevationChecker) error {
	elevated, err := checker.isElevated()
	if err != nil {
		return ierror.StatusError{Status: fmt.Sprintf("%s", err), StatusCode: 1}
	}

	if elevated {
		return nil
	}

	return ierror.StatusError{
		Status:     "setup must be run as administrator. please relaunch from an elevated shell (run as administrator).",
		StatusCode: 126,
	}
}
