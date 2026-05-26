//go:build windows

package file

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
