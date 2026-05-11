package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"strings"

	"github.com/spf13/afero"
)

var (
	validMountPrefixes = []string{"/run/media", "/media", "/mnt"}
	pathRegexp         = regexp.MustCompile(`^/[^\x00-\x1f\x7f]*$`)
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

// BlockDeviceLister defines the interface for listing the block devices
// on a Linux system.
type BlockDeviceLister interface {
	List() (*BlockDevices, error)
}

// DefaultBlockDeviceLister represents the primary implementation of the
// BlockDeviceLister interface
type DefaultBlockDeviceLister struct{}

// List returns a list of block devices on a Linux system by executing the
// `lsblk` command with the `--json` flag, parsing it's output into a BlockDevices struct.
// If the command execution or JSON parsing fails, an error is returned.
func (r *DefaultBlockDeviceLister) List() (*BlockDevices, error) {
	lsblkCmd := exec.Command("lsblk", "--json")
	lsblkOut, err := lsblkCmd.Output()
	if err != nil {
		return &BlockDevices{}, fmt.Errorf("%w", err)
	}

	var lsblk BlockDevices
	if err = json.Unmarshal(lsblkOut, &lsblk); err != nil {
		return &BlockDevices{}, fmt.Errorf("lsblk parsing: %w", err)
	}

	return &lsblk, nil
}

// SearchMountpoints searches for a file with the specified name across the
// provided mountpoints, skipping ones rejected by isMountpointSearchable.
// Walk errors on individual mountpoints (e.g. transient I/O issues, races
// with unmount) are treated as "no match" so other mountpoints still get
// scanned. If the file is found, the path is sent on c; context cancellation
// is honored on send.
func SearchMountpoints(ctx context.Context, fs afero.Fs, mountpoints []string, fileName string, c chan SearchResult) {
	for _, mountpoint := range mountpoints {
		if !isMountpointSearchable(mountpoint) {
			continue
		}

		filePath, err := GetFilePath(fs, mountpoint, fileName)
		if err != nil || filePath == "" {
			continue
		}

		select {
		case c <- SearchResult{Path: filePath, Err: nil}:
		case <-ctx.Done():
		}

		return
	}
}

// isMountpointSearchable determines whether a mountpoint is a valid candidate
// for a file search. It rejects empty paths, the root path, paths with invalid
// characters, and paths that don't fall under a known valid mount prefix.
func isMountpointSearchable(mountpoint string) bool {
	if mountpoint == "" || mountpoint == "/" {
		return false
	}

	if !pathRegexp.MatchString(mountpoint) {
		return false
	}

	return slices.ContainsFunc(validMountPrefixes, func(s string) bool {
		return strings.HasPrefix(mountpoint, s)
	})
}
