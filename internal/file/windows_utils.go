//go:build windows

package file

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf16"

	"golang.org/x/sys/windows"
)

var driveRootRegexp = regexp.MustCompile(`^[A-Za-z]:\\$`)

const (
	driveTypeRemovable = uint32(2)
	driveTypeFixed     = uint32(3)
)

// WinDrive represents a single logical drive on a Windows system, identified
// by its root path (e.g. "D:\"), classified by its Windows drive type constant,
// and tagged with the volume's filesystem name (e.g. "FAT32", "NTFS"). The
// filesystem name may be empty if the volume could not be inspected, which
// happens for empty card readers and similar "no media present" cases.
type WinDrive struct {
	Root       string
	Type       uint32
	FileSystem string
}

// WinDrives holds the list of logical drives enumerated on a Windows system.
type WinDrives struct {
	Drives []WinDrive
}

// WinDriveLister defines the interface for listing logical drives on a Windows system.
type WinDriveLister interface {
	List() (*WinDrives, error)
}

// DefaultWinDriveLister is the primary implementation of WinDriveLister.
type DefaultWinDriveLister struct{}

// WindowsSearcher implements DriveSearcher for Windows by fanning out
// goroutines across all logical drives reported by the WinDriveLister.
type WindowsSearcher struct {
	DriveLister WinDriveLister
}

// List enumerates all logical drives on the Windows system using the
// GetLogicalDriveStrings, GetDriveType, and GetVolumeInformation Win32 APIs.
// Each drive's root path, drive type, and filesystem name are recorded.
// GetVolumeInformation can fail (typically for empty card readers); that
// case is non-fatal — the drive is returned with an empty filesystem name
// and gets filtered out downstream by isDriveSearchable.
func (r *DefaultWinDriveLister) List() (*WinDrives, error) {
	buf := make([]uint16, 256)

	n, err := windows.GetLogicalDriveStrings(uint32(len(buf)), &buf[0])
	if err != nil {
		return &WinDrives{}, fmt.Errorf("GetLogicalDriveStrings: %w", err)
	}

	raw := buf[:n]
	var drives []WinDrive

	for len(raw) > 0 {
		end := 0
		for end < len(raw) && raw[end] != 0 {
			end++
		}

		if end > 0 {
			root := string(utf16.Decode(raw[:end]))
			rootPtr, err := windows.UTF16PtrFromString(root)
			if err != nil {
				return &WinDrives{}, fmt.Errorf("UTF16PtrFromString: %w", err)
			}
			drives = append(drives, WinDrive{
				Root:       root,
				Type:       windows.GetDriveType(rootPtr),
				FileSystem: getFileSystemName(rootPtr),
			})
		}

		if end >= len(raw) {
			break
		}

		raw = raw[end+1:]
	}

	return &WinDrives{Drives: drives}, nil
}

// getFileSystemName returns the filesystem name reported by the OS for the
// volume at root (e.g. "FAT32", "NTFS"). Returns "" if the volume cannot be
// probed, which is the normal case for empty removable drives.
func getFileSystemName(root *uint16) string {
	fsBuf := make([]uint16, windows.MAX_PATH+1)
	err := windows.GetVolumeInformation(
		root,
		nil, 0,
		nil, nil, nil,
		&fsBuf[0], uint32(len(fsBuf)),
	)
	if err != nil {
		return ""
	}

	end := 0
	for end < len(fsBuf) && fsBuf[end] != 0 {
		end++
	}

	return string(utf16.Decode(fsBuf[:end]))
}

// isDriveSearchable returns true if drive is a candidate for hosting a
// flashed Ubuntu boot partition. Both removable drives (USB SD readers) and
// fixed drives (built-in SD slots that enumerate as fixed) are accepted, but
// only when the filesystem is FAT-family — that excludes OS volumes like C:\
// (NTFS) and skips drives whose filesystem could not be probed.
func isDriveSearchable(drive WinDrive) bool {
	if drive.Root == "" {
		return false
	}

	if !driveRootRegexp.MatchString(drive.Root) {
		return false
	}

	if drive.Type != driveTypeRemovable && drive.Type != driveTypeFixed {
		return false
	}

	switch strings.ToUpper(drive.FileSystem) {
	case "FAT", "FAT32", "EXFAT":
		return true
	default:
		return false
	}
}
