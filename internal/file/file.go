// Package file implements utility routines for file operations and system interactions.
//
// It provides functions to search for files, list block devices, and search mountpoints on
// both macOS and Linux systems.
package file

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	ierror "github.com/borisdvlpr/gotail/internal/error"
	"github.com/spf13/afero"
)

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

	if err != nil && err != filepath.SkipDir {
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
		searchChan := make(chan SearchResult)
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
					SearchMountpoints(s.Fsys, mounts, fileName, searchChan)
				}(device.Mountpoints)
			}

			if device.Children != nil {
				for _, child := range device.Children {
					if child.Mountpoints != nil {
						wg.Add(1)
						go func(mounts []string) {
							defer wg.Done()
							SearchMountpoints(s.Fsys, mounts, fileName, searchChan)
						}(child.Mountpoints)
					}
				}
			}
		}

		go func() {
			wg.Wait()
			close(searchChan)
		}()

		for result := range searchChan {
			if result.Err != nil {
				return "", fmt.Errorf("%w", result.Err)
			}

			if result.Path != "" {
				return result.Path, nil
			}
		}

	default:
		status := fmt.Sprintf("unsupported operating system: %s", runtime.GOOS)
		return "", ierror.StatusError{Status: status, StatusCode: 71}
	}

	status := fmt.Sprintf("cannot access %s: could not find %s file, please try again", fileName, fileName)
	return "", ierror.StatusError{Status: status, StatusCode: 2}
}
