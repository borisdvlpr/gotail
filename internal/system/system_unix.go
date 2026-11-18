//go:build !windows

package system

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

// rootEnv holds the process and OS operations CheckRoot depends on, so its
// logic can be tested without running as root or replacing the test binary.
type rootEnv struct {
	geteuid    func() int
	lookPath   func(file string) (string, error)
	executable func() (string, error)
	execve     func(argv0 string, argv []string, envv []string) error
	args       []string
	environ    func() []string
}

// CheckRoot checks if the current user has root privileges. If not, it
// re-executes the binary with sudo. If the sudo command is not found or
// fails, an error is returned.
func (DefaultRootChecker) CheckRoot() error {
	return checkRootWithEnv(rootEnv{
		geteuid:    os.Geteuid,
		lookPath:   exec.LookPath,
		executable: os.Executable,
		execve:     syscall.Exec,
		args:       os.Args,
		environ:    os.Environ,
	})
}

// checkRootWithEnv implements CheckRoot using the operations in env. It
// returns nil if already root, a StatusError with code 127 if sudo cannot be
// found, code 1 if the executable path cannot be resolved, or code 126 if
// re-executing through sudo fails. On success, env.execve replaces the
// process and does not return.
func checkRootWithEnv(env rootEnv) error {
	if env.geteuid() == 0 {
		return nil
	}

	fmt.Println("Setup must be run as root. Re-executing with sudo...")

	sudoPath, err := env.lookPath("sudo")
	if err != nil {
		status := fmt.Sprintf("%s", err)
		return ierror.StatusError{Status: status, StatusCode: 127}
	}

	execPath, err := env.executable()
	if err != nil {
		status := fmt.Sprintf("failed to get executable path: %s", err)
		return ierror.StatusError{Status: status, StatusCode: 1}
	}

	args := []string{"sudo", execPath}
	args = append(args, env.args[1:]...)

	err = env.execve(sudoPath, args, env.environ())
	if err != nil {
		status := fmt.Sprintf("failed to execute sudo: %s", err)
		return ierror.StatusError{Status: status, StatusCode: 126}
	}

	return nil
}
