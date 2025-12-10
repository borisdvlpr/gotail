// Package error provides custom error types for the application.
//
// This package was heavily inspired by the Docker CLI error type found at
// https://github.com/docker/cli/blob/master/cli/error.go
package error

// StatusError reports an unsuccessful exit by a command, with its correspondent
// status code
type StatusError struct {
	Status     string
	StatusCode int
}

// Error formats the error for printing. If a custom Status is provided,
// it is returned as-is, otherwise it returns an 'abort' generic error message
func (e StatusError) Error() string {
	if e.Status == "" {
		return "abort"
	}
	return e.Status
}
