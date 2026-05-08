// Package system provides utility functions for system checks. It includes
// functions to check for root or administrative privileges so the setup
// command can ensure it has the access it needs to write to a mounted SD card.
package system

// RootChecker defines the interface for verifying if the application is currently
// running with sufficient root or administrative privileges.
type RootChecker interface {
	CheckRoot() error
}

// DefaultRootChecker represents the primary implementation of the RootChecker interface.
type DefaultRootChecker struct{}
