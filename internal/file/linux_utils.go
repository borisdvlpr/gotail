package file

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

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
