// Package file implements utility routines for file operations and system interactions.
//
// It provides functions to search for files, list block devices or logical drives, and search
// mountpoints on macOS, Linux, and Windows systems.
package file

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	ierror "github.com/borisdvlpr/gotail/internal/error"
	"github.com/spf13/afero"
)

var errFound = errors.New("file found")

// DriveSearcher is the platform-specific strategy for locating a file across
// the drives or mountpoints visible on the current OS. Each platform provides
// its own implementation via NewSystemSearcher; an unsupported operating system
// leaves this nil.
type DriveSearcher interface {
	Search(ctx context.Context, fsys afero.Fs, fileName string) (string, error)
}

// SystemSearcher holds the external interfaces used by the file module. It
// provides access to a virtual filesystem and a platform-specific DriveSearcher,
// enabling the module to locate files on removable media regardless of OS.
type SystemSearcher struct {
	Fsys     afero.Fs
	Searcher DriveSearcher
}

// SearchResult represents the outcome of a file search operation, containing either
// the path where the file was found or an error that occurred during the search process.
// It's used primarily for concurrent file search operations across multiple locations.
type SearchResult struct {
	Path string
	Err  error
}

// GetFilePath searches for a file with the specified name starting from the rootDir. It
// traverses the directory tree and returns the path of the first matching file found.
// Hidden directories and files are skipped during the search. If the file is found, its
// path is returned. If an error occurs during the search, it is returned.
func GetFilePath(fsys afero.Fs, rootDir string, fileName string) (string, error) {
	foundPath := ""

	err := afero.Walk(fsys, rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() && info.Name() == fileName {
			foundPath = path
			return errFound
		}

		return nil
	})

	if err != nil && !errors.Is(err, errFound) {
		return "", fmt.Errorf("%w", err)
	}

	return foundPath, nil
}

// FindUserData searches for the cloud-init configuration file on the current system.
// On macOS, it searches directly within /Volumes. On Linux and Windows, it
// delegates to the platform-specific DriveSearcher wired in by NewSystemSearcher.
// Returns the file path on success, or a StatusError if not found or unsupported.
func (s *SystemSearcher) FindUserData() (string, error) {
	const fileName = "user-data"

	if s.Searcher == nil {
		status := fmt.Sprintf("unsupported operating system: %s", runtime.GOOS)
		return "", ierror.StatusError{Status: status, StatusCode: 71}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	path, err := s.Searcher.Search(ctx, s.Fsys, fileName)
	if err != nil {
		return "", err
	}

	if path != "" {
		return path, nil
	}

	status := fmt.Sprintf("cannot access %s: could not find %s file, please try again", fileName, fileName)
	return "", ierror.StatusError{Status: status, StatusCode: 2}
}

// drainSearchChan reads results from c until it is closed or a non-empty
// path is found. The first hit cancels remaining goroutines via ctx.
func drainSearchChan(cancel context.CancelFunc, c chan SearchResult) (string, error) {
	for result := range c {
		if result.Err != nil {
			cancel()
			return "", fmt.Errorf("%w", result.Err)
		}
		if result.Path != "" {
			cancel()
			return result.Path, nil
		}
	}

	return "", nil
}
