package file

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

// BlockDevices is the representation of the block devices on a Linux
// system as returned by the `lsblk` command, containing information about each block device.
// If the block device has children (e.g., partitions), they are also included with similar information.
type BlockDevices struct {
	Blockdevices []struct {
		Name        string   `json:"name"`
		MajMin      string   `json:"maj:min"`
		Rm          bool     `json:"rm"`
		Size        string   `json:"size"`
		Ro          bool     `json:"ro"`
		Type        string   `json:"type"`
		Mountpoints []string `json:"mountpoints"`
		Children    []struct {
			Name        string   `json:"name"`
			MajMin      string   `json:"maj:min"`
			Rm          bool     `json:"rm"`
			Size        string   `json:"size"`
			Ro          bool     `json:"ro"`
			Type        string   `json:"type"`
			Mountpoints []string `json:"mountpoints"`
		} `json:"children,omitempty"`
	} `json:"blockdevices"`
}

// ListBlockDevices returns a list of block devices on a Linux system by executing the
// `lsblk` command with the `--json` flag, parsing it's output into a BlockDevices struct.
// If the command execution or JSON parsing fails, an error is returned.
func ListBlockDevices() (BlockDevices, error) {
	lsblkCmd := exec.Command("lsblk", "--json")
	lsblkOut, err := lsblkCmd.Output()
	if err != nil {
		return BlockDevices{}, fmt.Errorf("%w", err)
	}

	var lsblk BlockDevices
	if err = json.Unmarshal(lsblkOut, &lsblk); err != nil {
		return BlockDevices{}, fmt.Errorf("lsblk parsing: %w", err)
	}

	return lsblk, nil
}

// SearchMountpoints searches for the "user-data" file in the provided
// mountpoints iterating over them, ignoring certain paths.
// For valid mountpoints, it calls GetFilePath to find the "user-data" file.
// If the file is found, its path is returned. If an error occurs during the search, it is returned.
func SearchMountpoints(mountpoints []string) (string, error) {
	ignorePaths := []string{"/boot", "/home", "/snap"}

	for _, mountpoint := range mountpoints {
		if mountpoint != "" {
			validPath := !slices.ContainsFunc(ignorePaths, func(s string) bool {
				return strings.HasPrefix(mountpoint, s)
			})

			if mountpoint != "/" && validPath {
				filePath, err := GetFilePath(mountpoint, "user-data")
				if err != nil {
					return "", fmt.Errorf("%w", err)
				}

				if filePath != "" {
					return filePath, nil
				}
			}
		}
	}

	return "", nil
}
