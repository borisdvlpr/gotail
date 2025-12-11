package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/borisdvlpr/gotail/internal/file"
	"github.com/spf13/afero"
)

// MockRootChecker simulates checking for sudo/root permissions.
type MockRootChecker struct {
	ShouldError bool
}

func (m MockRootChecker) CheckRoot() error {
	if m.ShouldError {
		return fmt.Errorf("permission denied: must run as root")
	}
	return nil
}

// MockDeviceLister simulates finding hardware devices (SD cards).
type MockDeviceLister struct {
	Devices *file.BlockDevices
}

func (m *MockDeviceLister) List() (*file.BlockDevices, error) {
	if m.Devices == nil {
		return &file.BlockDevices{}, nil
	}
	return m.Devices, nil
}

func TestSetup_Success(t *testing.T) {
	mockFS := afero.NewMemMapFs()

	configFile := "/tmp/config.yaml"
	_ = afero.WriteFile(mockFS, configFile, []byte("exit_node: n\nsubnet_router: y\nsubnets: \"192.168.1.1/24,192.67.2.2/24,192.23.23.23/24\"\nhostname: raspberrypi\nauth_key: tskey_1234"), 0644)

	userDataPath := "/mnt/sdcard/user-data"
	_ = mockFS.MkdirAll("/mnt/sdcard", 0755)
	_ = afero.WriteFile(mockFS, userDataPath, []byte("#cloud-config\n"), 0644)

	mockRoot := MockRootChecker{ShouldError: false}
	mockLister := &MockDeviceLister{
		Devices: &file.BlockDevices{
			Blockdevices: []struct {
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
			}{
				{
					Name:        "mmcblk0",
					Type:        "disk",
					Mountpoints: []string{"/mnt/sdcard"},
				},
			},
		},
	}

	setupDeps := SetupCommand{
		Fsys:        mockFS,
		RootChecker: mockRoot,
		SystemSearcher: &file.SystemSearcher{
			Fsys:         mockFS,
			DeviceLister: mockLister,
		},
	}

	setupCmd := NewSetupCmd(setupDeps)

	var buf bytes.Buffer
	setupCmd.SetOut(&buf)
	setupCmd.SetErr(&buf)
	setupCmd.SetArgs([]string{"--file", configFile})

	err := setupCmd.Execute()

	if err != nil {
		t.Fatalf("Expected success, but got error: %v", err)
	}

	contentBytes, _ := afero.ReadFile(mockFS, userDataPath)
	content := string(contentBytes)

	if !strings.Contains(content, "- [ sh, -c, sudo hostnamectl hostname raspberrypi ]") {
		t.Errorf("Expected user-data to contain hostname 'raspberrypi', got:\n%s", content)
	}
}

func TestSetup_Fail_RootCheck(t *testing.T) {
	mockFS := afero.NewMemMapFs()

	setupDeps := SetupCommand{
		Fsys:        mockFS,
		RootChecker: MockRootChecker{ShouldError: true},
		SystemSearcher: &file.SystemSearcher{
			Fsys:         mockFS,
			DeviceLister: &MockDeviceLister{},
		},
	}

	setupCmd := NewSetupCmd(setupDeps)

	var buf bytes.Buffer
	setupCmd.SetOut(&buf)
	setupCmd.SetErr(&buf)
	setupCmd.SetArgs([]string{"--file", "dummy.yaml"})

	err := setupCmd.Execute()

	if err == nil {
		t.Error("Expected error due to root check, got nil")
	}

	if err != nil && !strings.Contains(err.Error(), "permission denied") {
		t.Errorf("Expected 'permission denied' error, got: %v", err)
	}
}

func TestSetup_Fail_FileNotFound(t *testing.T) {
	mockFS := afero.NewMemMapFs()

	setupDeps := SetupCommand{
		Fsys:        mockFS,
		RootChecker: MockRootChecker{ShouldError: false},
		SystemSearcher: &file.SystemSearcher{
			Fsys:         mockFS,
			DeviceLister: &MockDeviceLister{Devices: nil},
		},
	}

	setupCmd := NewSetupCmd(setupDeps)

	var buf bytes.Buffer
	setupCmd.SetOut(&buf)
	setupCmd.SetErr(&buf)
	setupCmd.SetArgs([]string{"--file", "dummy.yaml"})

	err := setupCmd.Execute()

	if err == nil {
		t.Error("Expected error because user-data file is missing, got nil")
	}

	if err != nil && !strings.Contains(err.Error(), "cannot access user-data") {
		t.Errorf("Expected 'cannot access user-data' error, got: %v", err)
	}
}
