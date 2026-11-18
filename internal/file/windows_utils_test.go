//go:build windows

package file

import "testing"

type IsDriveSearchableTestCase struct {
	id       string
	drive    WinDrive
	expected bool
}

func TestIsDriveSearchable(t *testing.T) {
	testCases := []IsDriveSearchableTestCase{
		{
			id:       "removable_fat32",
			drive:    WinDrive{Root: `D:\`, Type: driveTypeRemovable, FileSystem: "FAT32"},
			expected: true},
		{
			id:       "removable_fat",
			drive:    WinDrive{Root: `D:\`, Type: driveTypeRemovable, FileSystem: "FAT"},
			expected: true},
		{
			id:       "removable_exfat",
			drive:    WinDrive{Root: `D:\`, Type: driveTypeRemovable, FileSystem: "exFAT"},
			expected: true},
		{
			id:       "fixed_fat32_internal_sd",
			drive:    WinDrive{Root: `E:\`, Type: driveTypeFixed, FileSystem: "FAT32"},
			expected: true},
		{
			id:       "fixed_ntfs_os_volume",
			drive:    WinDrive{Root: `C:\`, Type: driveTypeFixed, FileSystem: "NTFS"},
			expected: false},
		{
			id:       "removable_no_filesystem",
			drive:    WinDrive{Root: `H:\`, Type: driveTypeRemovable, FileSystem: ""},
			expected: false},
		{
			id:       "removable_unknown_filesystem",
			drive:    WinDrive{Root: `D:\`, Type: driveTypeRemovable, FileSystem: "ext4"},
			expected: false},
		{
			id:       "cdrom_drive",
			drive:    WinDrive{Root: `E:\`, Type: 5, FileSystem: "CDFS"},
			expected: false},
		{
			id:       "network_drive",
			drive:    WinDrive{Root: `Z:\`, Type: 4, FileSystem: "NTFS"},
			expected: false},
		{
			id:       "empty_root",
			drive:    WinDrive{Root: "", Type: driveTypeRemovable, FileSystem: "FAT32"},
			expected: false},
		{
			id:       "missing_backslash",
			drive:    WinDrive{Root: "D:", Type: driveTypeRemovable, FileSystem: "FAT32"},
			expected: false},
		{
			id:       "letter_only",
			drive:    WinDrive{Root: "D", Type: driveTypeRemovable, FileSystem: "FAT32"},
			expected: false},
		{
			id:       "linux_path",
			drive:    WinDrive{Root: "/mnt/usb", Type: driveTypeRemovable, FileSystem: "FAT32"},
			expected: false},
		{
			id:       "forward_slash",
			drive:    WinDrive{Root: `D:/`, Type: driveTypeRemovable, FileSystem: "FAT32"},
			expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			if got := isDriveSearchable(tc.drive); got != tc.expected {
				t.Errorf("isDriveSearchable(%+v) = %v, want %v", tc.drive, got, tc.expected)
			}
		})
	}
}
