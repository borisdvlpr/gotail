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

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

// GetFilePath searches for a file with the specified name starting from the rootDir.
// It traverses the directory tree and returns the path of the first matching file found.
// Hidden directories and files are skipped during the search.
// If the file is found, its path is returned. If an error occurs during the search, it is returned.
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

	if runtime.GOOS == "darwin" {
		filePath, err = GetFilePath("/Volumes", fileName)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		if filePath != "" {
			return filePath, nil
		}
	}

	if runtime.GOOS == "linux" {
		devices, err := ListBlockDevices()
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		for _, device := range devices.Blockdevices {
			if device.Type == "loop" {
				continue
			}

			if device.Mountpoints != nil {
				filePath, err = SearchMountpoints(device.Mountpoints, fileName)
				if err != nil {
					return "", fmt.Errorf("%w", err)
				}

				if filePath != "" {
					return filePath, nil
				}
			}

			if device.Children != nil {
				for _, child := range device.Children {
					if child.Mountpoints != nil {
						filePath, err = SearchMountpoints(child.Mountpoints, fileName)
						if err != nil {
							return "", fmt.Errorf("%w", err)
						}

						if filePath != "" {
							return filePath, nil
						}
					}
				}
			}
		}
	}

	status := fmt.Sprintf("cannot access %s: could not find %s file, please try again", fileName, fileName)
	return "", ierror.StatusError{Status: status, StatusCode: 2}
}
