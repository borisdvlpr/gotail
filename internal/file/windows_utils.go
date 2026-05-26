//go:build windows

package file

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
