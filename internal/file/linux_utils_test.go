//go:build linux

package file

import "testing"

type IsMountpointSearchableTestCase struct {
	id         string
	mountpoint string
	expected   bool
}

func TestIsMountpointSearchable(t *testing.T) {
	testCases := []IsMountpointSearchableTestCase{
		{id: "valid_media", mountpoint: "/media/usb", expected: true},
		{id: "valid_mnt", mountpoint: "/mnt/data", expected: true},
		{id: "valid_run_media", mountpoint: "/run/media/user/disk", expected: true},
		{id: "root", mountpoint: "/", expected: false},
		{id: "empty", mountpoint: "", expected: false},
		{id: "invalid_prefix", mountpoint: "/home/user", expected: false},
		{id: "invalid_chars", mountpoint: "/mnt/\x00bad", expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			if got := isMountpointSearchable(tc.mountpoint); got != tc.expected {
				t.Errorf("isMountpointSearchable(%q) = %v, want %v", tc.mountpoint, got, tc.expected)
			}
		})
	}
}
