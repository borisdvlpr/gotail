//go:build linux

package file

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/spf13/afero"
)

// MockBlockDeviceLister replaces the lsblk call with a fixed device list,
// so the real LinuxSearcher fan-out runs against an in-memory filesystem.
type MockBlockDeviceLister struct {
	lsblkJSON string
	err       error
}

func (m MockBlockDeviceLister) List() (*BlockDevices, error) {
	if m.err != nil {
		return &BlockDevices{}, m.err
	}

	var devices BlockDevices
	if err := json.Unmarshal([]byte(m.lsblkJSON), &devices); err != nil {
		return &BlockDevices{}, err
	}

	return &devices, nil
}

type LinuxSearcherSearchTestCase struct {
	id           string
	lsblkJSON    string
	listErr      error
	files        []string
	expectedPath string
	expectError  bool
}

func TestLinuxSearcherSearch(t *testing.T) {
	testCases := []LinuxSearcherSearchTestCase{
		{
			id: "found_on_partition",
			lsblkJSON: `{"blockdevices": [{"name": "sdb", "type": "disk", "mountpoints": [null], "children": [
				{"name": "sdb1", "type": "part", "mountpoints": ["/media/user/bootfs"]},
				{"name": "sdb2", "type": "part", "mountpoints": ["/media/user/rootfs"]}]}]}`,
			files:        []string{"/media/user/bootfs/user-data"},
			expectedPath: "/media/user/bootfs/user-data",
		},
		{
			id:           "found_on_device_mountpoint",
			lsblkJSON:    `{"blockdevices": [{"name": "sdb", "type": "disk", "mountpoints": ["/mnt/sdcard"]}]}`,
			files:        []string{"/mnt/sdcard/user-data"},
			expectedPath: "/mnt/sdcard/user-data",
		},
		{
			id:           "loop_device_skipped",
			lsblkJSON:    `{"blockdevices": [{"name": "loop0", "type": "loop", "mountpoints": ["/mnt/snap"]}]}`,
			files:        []string{"/mnt/snap/user-data"},
			expectedPath: "",
		},
		{
			id: "unsearchable_mountpoint_skipped",
			lsblkJSON: `{"blockdevices": [{"name": "sda", "type": "disk", "mountpoints": [null], "children": [
				{"name": "sda1", "type": "part", "mountpoints": ["/home/user"]}]}]}`,
			files:        []string{"/home/user/user-data"},
			expectedPath: "",
		},
		{
			id: "not_found",
			lsblkJSON: `{"blockdevices": [{"name": "sdb", "type": "disk", "mountpoints": [null], "children": [
				{"name": "sdb1", "type": "part", "mountpoints": ["/media/user/bootfs"]}]}]}`,
			files:        []string{"/media/user/bootfs/config.txt"},
			expectedPath: "",
		},
		{
			// The first mountpoint does not exist, so the walk errors; the
			// search must move on to the next mountpoint of the partition.
			id: "missing_mountpoint_skipped",
			lsblkJSON: `{"blockdevices": [{"name": "sdb", "type": "disk", "mountpoints": [null], "children": [
				{"name": "sdb1", "type": "part", "mountpoints": ["/media/user/gone", "/media/user/bootfs"]}]}]}`,
			files:        []string{"/media/user/bootfs/user-data"},
			expectedPath: "/media/user/bootfs/user-data",
		},
		{
			id: "multiple_devices_one_match",
			lsblkJSON: `{"blockdevices": [
				{"name": "loop0", "type": "loop", "mountpoints": ["/mnt/snap"]},
				{"name": "sda", "type": "disk", "mountpoints": [null], "children": [
					{"name": "sda1", "type": "part", "mountpoints": ["/"]},
					{"name": "sda2", "type": "part", "mountpoints": ["/mnt/data"]}]},
				{"name": "sdb", "type": "disk", "mountpoints": [null], "children": [
					{"name": "sdb1", "type": "part", "mountpoints": ["/run/media/user/bootfs"]}]}]}`,
			files:        []string{"/mnt/snap/user-data", "/mnt/data/notes.txt", "/run/media/user/bootfs/user-data"},
			expectedPath: "/run/media/user/bootfs/user-data",
		},
		{
			id:          "lister_error",
			listErr:     errors.New("lsblk: executable file not found in $PATH"),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		fsys := afero.NewMemMapFs()
		for _, f := range tc.files {
			if err := afero.WriteFile(fsys, f, []byte{}, 0644); err != nil {
				t.Fatalf("%v: failed to create test file: %v", tc.id, err)
			}
		}

		searcher := &LinuxSearcher{DeviceLister: MockBlockDeviceLister{lsblkJSON: tc.lsblkJSON, err: tc.listErr}}
		path, err := searcher.Search(context.Background(), fsys, "user-data")

		if tc.expectError {
			if err == nil {
				t.Errorf("%v: Search() returned no error, expected one", tc.id)
			}
			continue
		}

		if err != nil {
			t.Errorf("%v: Search() returned error %v, expected none", tc.id, err)
		}

		if path != tc.expectedPath {
			t.Errorf("%v: Search() path = %q, expected %q", tc.id, path, tc.expectedPath)
		}
	}
}

func TestNewSystemSearcher(t *testing.T) {
	fsys := afero.NewMemMapFs()
	s := NewSystemSearcher(fsys)

	if s.Fsys != fsys {
		t.Errorf("NewSystemSearcher() Fsys = %v, expected the given filesystem", s.Fsys)
	}

	ls, ok := s.Searcher.(*LinuxSearcher)
	if !ok {
		t.Fatalf("NewSystemSearcher() Searcher = %T, expected *LinuxSearcher", s.Searcher)
	}

	if _, ok := ls.DeviceLister.(*DefaultBlockDeviceLister); !ok {
		t.Errorf("NewSystemSearcher() DeviceLister = %T, expected *DefaultBlockDeviceLister", ls.DeviceLister)
	}
}

type IsMountpointSearchableTestCase struct {
	id         string
	mountpoint string
	expected   bool
}

func TestIsMountpointSearchable(t *testing.T) {
	testCases := []IsMountpointSearchableTestCase{
		{
			id:         "valid_media",
			mountpoint: "/media/usb",
			expected:   true},
		{
			id:         "valid_mnt",
			mountpoint: "/mnt/data",
			expected:   true},
		{
			id:         "valid_run_media",
			mountpoint: "/run/media/user/disk",
			expected:   true},
		{
			id:         "root",
			mountpoint: "/",
			expected:   false},
		{
			id:         "empty",
			mountpoint: "",
			expected:   false},
		{
			id:         "invalid_prefix",
			mountpoint: "/home/user",
			expected:   false},
		{
			id:         "invalid_chars",
			mountpoint: "/mnt/\x00bad",
			expected:   false},
		{
			id:         "exact_prefix",
			mountpoint: "/mnt",
			expected:   true},
		{
			id:         "prefix_without_separator_mnt",
			mountpoint: "/mntfoo",
			expected:   false},
		{
			id:         "prefix_without_separator_media",
			mountpoint: "/mediabackup",
			expected:   false},
		{
			id:         "prefix_without_separator_run_media",
			mountpoint: "/run/media-old",
			expected:   false},
	}

	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			if got := isMountpointSearchable(tc.mountpoint); got != tc.expected {
				t.Errorf("isMountpointSearchable(%q) = %v, want %v", tc.mountpoint, got, tc.expected)
			}
		})
	}
}
