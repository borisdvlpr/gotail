//go:build !windows

package system

import (
	"errors"
	"slices"
	"testing"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

type CheckRootWithEnvTestCase struct {
	id            string
	euid          int
	lookPathErr   error
	executableErr error
	execErr       error
	expectedError error
	expectExec    bool
}

func TestCheckRootWithEnv(t *testing.T) {
	testCases := []CheckRootWithEnvTestCase{
		{
			id:            "already_root",
			euid:          0,
			expectedError: nil,
			expectExec:    false,
		},
		{
			id:            "sudo_not_found",
			euid:          1000,
			lookPathErr:   errors.New("sudo not found"),
			expectedError: ierror.StatusError{Status: "sudo not found", StatusCode: 127},
			expectExec:    false,
		},
		{
			id:            "executable_path_unavailable",
			euid:          1000,
			executableErr: errors.New("no such file"),
			expectedError: ierror.StatusError{Status: "failed to get executable path: no such file", StatusCode: 1},
			expectExec:    false,
		},
		{
			id:            "exec_fails",
			euid:          1000,
			execErr:       errors.New("permission denied"),
			expectedError: ierror.StatusError{Status: "failed to execute sudo: permission denied", StatusCode: 126},
			expectExec:    true,
		},
		{
			id:            "reexecs_through_sudo",
			euid:          1000,
			expectedError: nil,
			expectExec:    true,
		},
	}

	for _, tc := range testCases {
		var execCalled bool
		var gotArgv0 string
		var gotArgv, gotEnv []string

		env := rootEnv{
			geteuid: func() int { return tc.euid },
			lookPath: func(file string) (string, error) {
				if tc.lookPathErr != nil {
					return "", tc.lookPathErr
				}
				return "/usr/bin/" + file, nil
			},
			executable: func() (string, error) {
				if tc.executableErr != nil {
					return "", tc.executableErr
				}
				return "/usr/local/bin/gotail", nil
			},
			execve: func(argv0 string, argv []string, envv []string) error {
				execCalled = true
				gotArgv0, gotArgv, gotEnv = argv0, argv, envv
				return tc.execErr
			},
			args:    []string{"gotail", "setup", "--file", "config.yaml"},
			environ: func() []string { return []string{"HOME=/home/pi"} },
		}

		err := checkRootWithEnv(env)

		if tc.expectedError == nil && err != nil {
			t.Errorf("%v: checkRootWithEnv() returned error %v, expected no error", tc.id, err)
		}

		if tc.expectedError != nil && err == nil {
			t.Errorf("%v: checkRootWithEnv() returned no error, expected %q", tc.id, tc.expectedError)
		}

		if tc.expectedError != nil && err != nil && !errors.Is(err, tc.expectedError) {
			t.Errorf("%v: checkRootWithEnv() returned %q, expected %q", tc.id, err, tc.expectedError)
		}

		if execCalled != tc.expectExec {
			t.Errorf("%v: checkRootWithEnv() re-exec called = %v, expected %v", tc.id, execCalled, tc.expectExec)
		}

		if tc.expectExec {
			wantArgv := []string{"sudo", "/usr/local/bin/gotail", "setup", "--file", "config.yaml"}
			wantEnv := []string{"HOME=/home/pi"}

			if gotArgv0 != "/usr/bin/sudo" || !slices.Equal(gotArgv, wantArgv) || !slices.Equal(gotEnv, wantEnv) {
				t.Errorf("%v: checkRootWithEnv() re-exec(%q, %q, %q), expected (%q, %q, %q)",
					tc.id, gotArgv0, gotArgv, gotEnv, "/usr/bin/sudo", wantArgv, wantEnv)
			}
		}
	}
}
