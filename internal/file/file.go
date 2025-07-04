// Package file implements utility routines for file operations and system interactions.
//
// It provides functions to search for files, list block devices, and search mountpoints on
// both macOS and Linux systems.
package file

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

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
func GetFilePath(rootDir string, fileName string) (string, error) {
	filePath := ""

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		if !d.IsDir() && d.Name() == fileName {
			filePath = path
			return fs.SkipDir
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return filePath, nil
}

// FindUserData searches for the config file on the system.
// On macOS, it searches within the "/Volumes" directory.
// On Linux, it lists block devices and searches their mountpoints and their children mountpoints.
// If the file is found, its path is returned. If an error occurs during the search, it is returned.
func FindUserData() (string, error) {
	const fileName = "user-data"
	var filePath string
	var err error

	switch runtime.GOOS {
	case "darwin":
		filePath, err = GetFilePath("/Volumes", fileName)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		if filePath != "" {
			return filePath, nil
		}

	case "linux":
		searchChan := make(chan SearchResult)
		var wg sync.WaitGroup

		devices, err := ListBlockDevices()
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
					SearchMountpoints(mounts, fileName, searchChan)
				}(device.Mountpoints)
			}

			if device.Children != nil {
				for _, child := range device.Children {
					if child.Mountpoints != nil {
						wg.Add(1)
						go func(mounts []string) {
							defer wg.Done()
							SearchMountpoints(mounts, fileName, searchChan)
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
		return "", ierror.StatusError{Status: status, StatusCode: 71} // EX_OSERR
	}

	status := fmt.Sprintf("cannot access %s: could not find %s file, please try again", fileName, fileName)
	return "", ierror.StatusError{Status: status, StatusCode: 2}
}
