// Package file implements utility routines for file operations and system interactions.
//
// It provides functions to search for files, list block devices, and search mountpoints on
// both macOS and Linux systems.
package file

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	ierror "github.com/borisdvlpr/gotail/internal/error"
	"github.com/spf13/afero"
)

// DriveSearcher is the platform-specific strategy for locating a file across
// the drives or mountpoints visible on the current OS. Each platform provides
// its own implementation via NewSystemSearcher; darwin leaves this nil and
// falls back to the inline /Volumes walk.
type DriveSearcher interface {
	Search(ctx context.Context, fsys afero.Fs, fileName string, c chan SearchResult)
}

// SystemSearcher holds the external interfaces used by the file module. It
// provides access to a virtual filesystem and block-device listings, enabling
// the module to inspect mountpoints and search for files within them.
type SystemSearcher struct {
	Fsys         afero.Fs
	DeviceLister BlockDeviceLister
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
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil && !errors.Is(err, filepath.SkipDir) {
		return "", fmt.Errorf("%w", err)
	}

	return foundPath, nil
}

// FindUserData searches for the config file on the system.
// On macOS, it searches within the "/Volumes" directory.
// On Linux, it lists block devices and searches their mountpoints and their children mountpoints.
// If the file is found, its path is returned. If an error occurs during the search, it is returned.
func (s *SystemSearcher) FindUserData() (string, error) {
	const fileName = "user-data"
	var filePath string
	var err error

	switch runtime.GOOS {
	case "darwin":
		filePath, err = GetFilePath(s.Fsys, "/Volumes", fileName)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		if filePath != "" {
			return filePath, nil
		}

	case "linux":
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		searchChan := make(chan SearchResult, 1)
		var wg sync.WaitGroup

		devices, err := s.DeviceLister.List()
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		for _, device := range devices.Blockdevices {
			if device.Type == "loop" {
				continue
			}

			if device.Mountpoints != nil {
				wg.Add(1)
				go func(mounts []string) {
					defer wg.Done()
					SearchMountpoints(ctx, s.Fsys, mounts, fileName, searchChan)
				}(device.Mountpoints)
			}

			if device.Children != nil {
				for _, child := range device.Children {
					if child.Mountpoints != nil {
						wg.Add(1)
						go func(mounts []string) {
							defer wg.Done()
							SearchMountpoints(ctx, s.Fsys, mounts, fileName, searchChan)
						}(child.Mountpoints)
					}
				}
			}
		}

		go func() {
			wg.Wait()
			close(searchChan)
		}()

		path, err := drainSearchChan(cancel, searchChan)
		if err != nil || path != "" {
			return path, err
		}

	default:
		status := fmt.Sprintf("unsupported operating system: %s", runtime.GOOS)
		return "", ierror.StatusError{Status: status, StatusCode: 71}
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
