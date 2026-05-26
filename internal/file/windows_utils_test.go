//go:build windows

package file

import "testing"

type IsDriveSearchableTestCase struct {
	name     string
	drive    WinDrive
	expected bool
}

func TestIsDriveSearchable(t *testing.T) {
	testCases := []IsDriveSearchableTestCase{
		{name: "removable_fat32", drive: WinDrive{Root: `D:\`, Type: driveTypeRemovable, FileSystem: "FAT32"}, expected: true},
		{name: "removable_fat", drive: WinDrive{Root: `D:\`, Type: driveTypeRemovable, FileSystem: "FAT"}, expected: true},
		{name: "removable_exfat", drive: WinDrive{Root: `D:\`, Type: driveTypeRemovable, FileSystem: "exFAT"}, expected: true},
		{name: "fixed_fat32_internal_sd", drive: WinDrive{Root: `E:\`, Type: driveTypeFixed, FileSystem: "FAT32"}, expected: true},
		{name: "fixed_ntfs_os_volume", drive: WinDrive{Root: `C:\`, Type: driveTypeFixed, FileSystem: "NTFS"}, expected: false},
		{name: "removable_no_filesystem", drive: WinDrive{Root: `H:\`, Type: driveTypeRemovable, FileSystem: ""}, expected: false},
		{name: "removable_unknown_filesystem", drive: WinDrive{Root: `D:\`, Type: driveTypeRemovable, FileSystem: "ext4"}, expected: false},
		{name: "cdrom_drive", drive: WinDrive{Root: `E:\`, Type: 5, FileSystem: "CDFS"}, expected: false},
		{name: "network_drive", drive: WinDrive{Root: `Z:\`, Type: 4, FileSystem: "NTFS"}, expected: false},
		{name: "empty_root", drive: WinDrive{Root: "", Type: driveTypeRemovable, FileSystem: "FAT32"}, expected: false},
		{name: "missing_backslash", drive: WinDrive{Root: "D:", Type: driveTypeRemovable, FileSystem: "FAT32"}, expected: false},
		{name: "letter_only", drive: WinDrive{Root: "D", Type: driveTypeRemovable, FileSystem: "FAT32"}, expected: false},
		{name: "linux_path", drive: WinDrive{Root: "/mnt/usb", Type: driveTypeRemovable, FileSystem: "FAT32"}, expected: false},
		{name: "forward_slash", drive: WinDrive{Root: `D:/`, Type: driveTypeRemovable, FileSystem: "FAT32"}, expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isDriveSearchable(tc.drive); got != tc.expected {
				t.Errorf("isDriveSearchable(%+v) = %v, want %v", tc.drive, got, tc.expected)
			}
		})
	}
}
