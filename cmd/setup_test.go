package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	ierror "github.com/borisdvlpr/gotail/internal/error"
	"github.com/borisdvlpr/gotail/internal/file"
	"github.com/spf13/afero"
)

type MockRootChecker struct {
	ShouldError bool
}

func (m MockRootChecker) CheckRoot() error {
	if m.ShouldError {
		return fmt.Errorf("permission denied: must run as root")
	}
	return nil
}

type MockSearcher struct{}

func (m *MockSearcher) Search(ctx context.Context, fsys afero.Fs, fileName string) (string, error) {
	return file.GetFilePath(fsys, "/", fileName)
}

// captureOutput runs fn with os.Stdout and os.Stderr redirected to pipes and
// returns what was written to each. The command under test must not have
// SetOut/SetErr called, so cobra falls back to the real process streams.
func captureOutput(t *testing.T, fn func()) (string, string) {
	t.Helper()

	origStdout, origStderr := os.Stdout, os.Stderr
	defer func() { os.Stdout, os.Stderr = origStdout, origStderr }()

	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}

	errR, errW, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stderr pipe: %v", err)
	}

	os.Stdout, os.Stderr = outW, errW
	fn()
	_ = outW.Close()
	_ = errW.Close()

	stdout, _ := io.ReadAll(outR)
	stderr, _ := io.ReadAll(errR)

	return string(stdout), string(stderr)
}

func TestSetup_Success(t *testing.T) {
	mockFS := afero.NewMemMapFs()

	configFile := "/tmp/config.yaml"
	_ = afero.WriteFile(mockFS, configFile, []byte("exit_node: n\nsubnet_router: y\nsubnets: \"192.168.1.0/24,192.67.2.0/24,192.23.23.0/24\"\nhostname: raspberrypi\nauth_key: tskey_1234"), 0644)

	userDataPath := "/mnt/sdcard/user-data"
	_ = mockFS.MkdirAll("/mnt/sdcard", 0755)
	_ = afero.WriteFile(mockFS, userDataPath, []byte("# cloud-config\n"), 0644)

	setupDeps := SetupCommand{
		Fsys:        mockFS,
		RootChecker: MockRootChecker{ShouldError: false},
		SystemSearcher: &file.SystemSearcher{
			Fsys:     mockFS,
			Searcher: &MockSearcher{},
		},
	}

	setupCmd := NewSetupCmd(setupDeps)
	setupCmd.SetArgs([]string{"--file", configFile})

	var err error
	stdout, stderr := captureOutput(t, func() { err = setupCmd.Execute() })

	if err != nil {
		t.Fatalf("Expected success, but got error: %v", err)
	}

	if !strings.Contains(stdout, "Tailscale will be installed on boot") || stderr != "" {
		t.Errorf("Expected progress messages on stdout only, got stdout:\n%s\nstderr:\n%s", stdout, stderr)
	}

	contentBytes, _ := afero.ReadFile(mockFS, userDataPath)
	content := string(contentBytes)

	if !strings.Contains(content, `- ["sh", "-c", "sudo hostnamectl hostname raspberrypi"]`) {
		t.Errorf("Expected user-data to contain hostname 'raspberrypi', got:\n%s", content)
	}
}

func TestSetup_SubnetsNormalized(t *testing.T) {
	mockFS := afero.NewMemMapFs()

	configFile := "/tmp/config.yaml"
	_ = afero.WriteFile(mockFS, configFile, []byte("exit_node: n\nsubnet_router: y\nsubnets: \"192.168.1.0/24, 10.0.0.0/16\"\nhostname: raspberrypi\nauth_key: tskey_1234"), 0644)

	userDataPath := "/mnt/sdcard/user-data"
	_ = mockFS.MkdirAll("/mnt/sdcard", 0755)
	_ = afero.WriteFile(mockFS, userDataPath, []byte("# cloud-config\n"), 0644)

	setupDeps := SetupCommand{
		Fsys:        mockFS,
		RootChecker: MockRootChecker{ShouldError: false},
		SystemSearcher: &file.SystemSearcher{
			Fsys:     mockFS,
			Searcher: &MockSearcher{},
		},
	}

	setupCmd := NewSetupCmd(setupDeps)

	var buf bytes.Buffer
	setupCmd.SetOut(&buf)
	setupCmd.SetErr(&buf)
	setupCmd.SetArgs([]string{"--file", configFile})

	if err := setupCmd.Execute(); err != nil {
		t.Fatalf("Expected success, but got error: %v", err)
	}

	contentBytes, _ := afero.ReadFile(mockFS, userDataPath)
	content := string(contentBytes)

	if !strings.Contains(content, `"--advertise-routes=192.168.1.0/24,10.0.0.0/16"`) {
		t.Errorf("Expected user-data to contain whitespace-free advertise-routes, got:\n%s", content)
	}
}

func TestSetup_ExitNodeEnablesForwarding(t *testing.T) {
	mockFS := afero.NewMemMapFs()

	configFile := "/tmp/config.yaml"
	_ = afero.WriteFile(mockFS, configFile, []byte("exit_node: y\nsubnet_router: n\nhostname: raspberrypi\nauth_key: tskey_1234"), 0644)

	userDataPath := "/mnt/sdcard/user-data"
	_ = mockFS.MkdirAll("/mnt/sdcard", 0755)
	_ = afero.WriteFile(mockFS, userDataPath, []byte("# cloud-config\n"), 0644)

	setupDeps := SetupCommand{
		Fsys:        mockFS,
		RootChecker: MockRootChecker{ShouldError: false},
		SystemSearcher: &file.SystemSearcher{
			Fsys:     mockFS,
			Searcher: &MockSearcher{},
		},
	}

	setupCmd := NewSetupCmd(setupDeps)

	var buf bytes.Buffer
	setupCmd.SetOut(&buf)
	setupCmd.SetErr(&buf)
	setupCmd.SetArgs([]string{"--file", configFile})

	if err := setupCmd.Execute(); err != nil {
		t.Fatalf("Expected success, but got error: %v", err)
	}

	contentBytes, _ := afero.ReadFile(mockFS, userDataPath)
	content := string(contentBytes)

	if got := strings.Count(content, "net.ipv4.ip_forward = 1"); got != 1 {
		t.Errorf("Expected IP forwarding to be enabled once for an exit node, found %d times in:\n%s", got, content)
	}
}

func TestSetup_Fail_InvalidHostname(t *testing.T) {
	mockFS := afero.NewMemMapFs()

	configFile := "/tmp/config.yaml"
	_ = afero.WriteFile(mockFS, configFile, []byte("exit_node: n\nsubnet_router: n\nhostname: \"pi; reboot\"\nauth_key: tskey_1234"), 0644)

	userDataPath := "/mnt/sdcard/user-data"
	_ = mockFS.MkdirAll("/mnt/sdcard", 0755)
	_ = afero.WriteFile(mockFS, userDataPath, []byte("# cloud-config\n"), 0644)

	setupDeps := SetupCommand{
		Fsys:        mockFS,
		RootChecker: MockRootChecker{ShouldError: false},
		SystemSearcher: &file.SystemSearcher{
			Fsys:     mockFS,
			Searcher: &MockSearcher{},
		},
	}

	setupCmd := NewSetupCmd(setupDeps)

	var buf bytes.Buffer
	setupCmd.SetOut(&buf)
	setupCmd.SetErr(&buf)
	setupCmd.SetArgs([]string{"--file", configFile})

	err := setupCmd.Execute()

	if err == nil || !strings.Contains(err.Error(), "invalid hostname") {
		t.Errorf("Expected 'invalid hostname' error, got: %v", err)
	}

	contentBytes, _ := afero.ReadFile(mockFS, userDataPath)
	if content := string(contentBytes); content != "# cloud-config\n" {
		t.Errorf("Expected user-data to be left unchanged, got:\n%s", content)
	}
}

func TestSetup_UserDataWithoutTrailingNewline(t *testing.T) {
	mockFS := afero.NewMemMapFs()

	configFile := "/tmp/config.yaml"
	_ = afero.WriteFile(mockFS, configFile, []byte("exit_node: n\nsubnet_router: n\nhostname: raspberrypi\nauth_key: tskey_1234"), 0644)

	userDataPath := "/mnt/sdcard/user-data"
	_ = mockFS.MkdirAll("/mnt/sdcard", 0755)
	_ = afero.WriteFile(mockFS, userDataPath, []byte("# cloud-config"), 0644)

	setupDeps := SetupCommand{
		Fsys:        mockFS,
		RootChecker: MockRootChecker{ShouldError: false},
		SystemSearcher: &file.SystemSearcher{
			Fsys:     mockFS,
			Searcher: &MockSearcher{},
		},
	}

	setupCmd := NewSetupCmd(setupDeps)

	var buf bytes.Buffer
	setupCmd.SetOut(&buf)
	setupCmd.SetErr(&buf)
	setupCmd.SetArgs([]string{"--file", configFile})

	if err := setupCmd.Execute(); err != nil {
		t.Fatalf("Expected success, but got error: %v", err)
	}

	contentBytes, _ := afero.ReadFile(mockFS, userDataPath)
	content := string(contentBytes)

	if !strings.HasPrefix(content, "# cloud-config\n") || !strings.Contains(content, "\nruncmd:\n") {
		t.Errorf("Expected 'runcmd:' on its own line after the existing content, got:\n%s", content)
	}
}

func TestSetup_Fail_RootCheck(t *testing.T) {
	mockFS := afero.NewMemMapFs()

	setupDeps := SetupCommand{
		Fsys:        mockFS,
		RootChecker: MockRootChecker{ShouldError: true},
		SystemSearcher: &file.SystemSearcher{
			Fsys:     mockFS,
			Searcher: &MockSearcher{},
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
			Fsys:     mockFS,
			Searcher: &MockSearcher{}, // walks empty FS, will find nothing
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

// MockFailingWriteFs wraps an afero.Fs so that every file opened through
// OpenFile rejects writes, simulating an I/O failure on the SD card.
type MockFailingWriteFs struct {
	afero.Fs
}

func (m MockFailingWriteFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := m.Fs.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return mockFailingWriteFile{f}, nil
}

type mockFailingWriteFile struct {
	afero.File
}

func (mockFailingWriteFile) Write([]byte) (int, error) {
	return 0, errors.New("input/output error")
}

type SetupIOErrorTestCase struct {
	id           string
	config       string
	wrapFs       func(afero.Fs) afero.Fs
	expectedCode int
}

func TestSetup_IOErrorExitCodes(t *testing.T) {
	validConfig := "exit_node: n\nsubnet_router: n\nhostname: raspberrypi\nauth_key: tskey_1234"
	noWrap := func(fs afero.Fs) afero.Fs { return fs }

	testCases := []SetupIOErrorTestCase{
		{
			id:           "config_file_missing",
			config:       "",
			wrapFs:       noWrap,
			expectedCode: 66, // EX_NOINPUT
		},
		{
			id:           "config_file_invalid_yaml",
			config:       "exit_node: [",
			wrapFs:       noWrap,
			expectedCode: 65, // EX_DATAERR
		},
		{
			id:           "user_data_open_fails",
			config:       validConfig,
			wrapFs:       afero.NewReadOnlyFs,
			expectedCode: 74, // EX_IOERR
		},
		{
			id:           "user_data_write_fails",
			config:       validConfig,
			wrapFs:       func(fs afero.Fs) afero.Fs { return MockFailingWriteFs{fs} },
			expectedCode: 74, // EX_IOERR
		},
	}

	for _, tc := range testCases {
		mockFS := afero.NewMemMapFs()

		configFile := "/tmp/config.yaml"
		if tc.config != "" {
			_ = afero.WriteFile(mockFS, configFile, []byte(tc.config), 0644)
		}

		_ = mockFS.MkdirAll("/mnt/sdcard", 0755)
		_ = afero.WriteFile(mockFS, "/mnt/sdcard/user-data", []byte("# cloud-config\n"), 0644)

		setupDeps := SetupCommand{
			Fsys:        tc.wrapFs(mockFS),
			RootChecker: MockRootChecker{ShouldError: false},
			SystemSearcher: &file.SystemSearcher{
				Fsys:     mockFS,
				Searcher: &MockSearcher{},
			},
		}

		setupCmd := NewSetupCmd(setupDeps)

		var buf bytes.Buffer
		setupCmd.SetOut(&buf)
		setupCmd.SetErr(&buf)
		setupCmd.SetArgs([]string{"--file", configFile})

		err := setupCmd.Execute()

		var statusErr ierror.StatusError
		if !errors.As(err, &statusErr) {
			t.Errorf("%v: Execute() returned %v (%T), expected ierror.StatusError", tc.id, err, err)
			continue
		}

		if statusErr.StatusCode != tc.expectedCode {
			t.Errorf("%v: Execute() returned status code %d (%q), expected %d", tc.id, statusErr.StatusCode, statusErr, tc.expectedCode)
		}
	}
}

func TestSetup_Interactive(t *testing.T) {
	mockFS := afero.NewMemMapFs()

	userDataPath := "/mnt/sdcard/user-data"
	_ = mockFS.MkdirAll("/mnt/sdcard", 0755)
	_ = afero.WriteFile(mockFS, userDataPath, []byte("# cloud-config\n"), 0644)

	setupDeps := SetupCommand{
		Fsys:        mockFS,
		RootChecker: MockRootChecker{ShouldError: false},
		SystemSearcher: &file.SystemSearcher{
			Fsys:     mockFS,
			Searcher: &MockSearcher{},
		},
	}

	setupCmd := NewSetupCmd(setupDeps)
	// All answers are available up front, as when they are piped in.
	setupCmd.SetIn(strings.NewReader("n\ny\n10.0.0.0/24, 192.168.1.0/24\nraspberrypi\nServer,prod-2\ntskey_1234\n"))
	setupCmd.SetArgs([]string{})

	var err error
	stdout, _ := captureOutput(t, func() { err = setupCmd.Execute() })

	if err != nil {
		t.Fatalf("Expected success, but got error: %v", err)
	}

	if !strings.Contains(stdout, "Setup device as an exit node? [y/n]") {
		t.Errorf("Expected the exit node prompt on stdout, got:\n%s", stdout)
	}

	if !strings.Contains(stdout, "This device will be tagged: tag:server,tag:prod-2.") {
		t.Errorf("Expected the normalized tags on stdout, got:\n%s", stdout)
	}

	contentBytes, _ := afero.ReadFile(mockFS, userDataPath)
	content := string(contentBytes)

	expected := `- ["tailscale", "up", "--ssh", "--advertise-routes=10.0.0.0/24,192.168.1.0/24", "--hostname=raspberrypi", "--advertise-tags=tag:server,tag:prod-2", "--authkey=tskey_1234"]`
	if !strings.Contains(content, expected) {
		t.Errorf("Expected user-data to contain\n%s\ngot:\n%s", expected, content)
	}

	if got := strings.Count(content, "net.ipv4.ip_forward = 1"); got != 1 {
		t.Errorf("Expected IP forwarding to be enabled once for a subnet router, found %d times in:\n%s", got, content)
	}
}

type SetupInteractiveFailTestCase struct {
	id           string
	answers      string
	expectedErr  error
	expectedCode int
}

func TestSetup_Fail_Interactive(t *testing.T) {
	testCases := []SetupInteractiveFailTestCase{
		{
			id:           "exit_node_invalid_reply",
			answers:      "maybe\n",
			expectedCode: 1,
		},
		{
			id:          "exit_node_eof",
			answers:     "",
			expectedErr: io.EOF,
		},
		{
			id:           "subnet_router_invalid_reply",
			answers:      "n\nmaybe\n",
			expectedCode: 1,
		},
		{
			id:          "subnets_eof",
			answers:     "n\ny\n",
			expectedErr: io.EOF,
		},
		{
			id:          "hostname_eof",
			answers:     "n\nn\n",
			expectedErr: io.EOF,
		},
		{
			id:          "tags_eof",
			answers:     "n\nn\nraspberrypi\n",
			expectedErr: io.EOF,
		},
		{
			id:          "authkey_eof",
			answers:     "n\nn\nraspberrypi\n\n",
			expectedErr: io.EOF,
		},
		{
			id:           "authkey_empty",
			answers:      "n\nn\nraspberrypi\n\n\n",
			expectedCode: 1,
		},
	}

	for _, tc := range testCases {
		mockFS := afero.NewMemMapFs()

		userDataPath := "/mnt/sdcard/user-data"
		_ = mockFS.MkdirAll("/mnt/sdcard", 0755)
		_ = afero.WriteFile(mockFS, userDataPath, []byte("# cloud-config\n"), 0644)

		setupDeps := SetupCommand{
			Fsys:        mockFS,
			RootChecker: MockRootChecker{ShouldError: false},
			SystemSearcher: &file.SystemSearcher{
				Fsys:     mockFS,
				Searcher: &MockSearcher{},
			},
		}

		setupCmd := NewSetupCmd(setupDeps)
		setupCmd.SetIn(strings.NewReader(tc.answers))
		setupCmd.SetArgs([]string{})

		var err error
		captureOutput(t, func() { err = setupCmd.Execute() })

		if tc.expectedErr != nil && !errors.Is(err, tc.expectedErr) {
			t.Errorf("%v: Execute() returned %v, expected %v", tc.id, err, tc.expectedErr)
		}

		if tc.expectedCode != 0 {
			var statusErr ierror.StatusError
			if !errors.As(err, &statusErr) || statusErr.StatusCode != tc.expectedCode {
				t.Errorf("%v: Execute() returned %v, expected StatusError with code %d", tc.id, err, tc.expectedCode)
			}
		}

		contentBytes, _ := afero.ReadFile(mockFS, userDataPath)
		if content := string(contentBytes); content != "# cloud-config\n" {
			t.Errorf("%v: Expected user-data to be left unchanged, got:\n%s", tc.id, content)
		}
	}
}

type SetupInvalidConfigTestCase struct {
	id            string
	config        string
	expectedError string
}

func TestSetup_Fail_InvalidConfig(t *testing.T) {
	testCases := []SetupInvalidConfigTestCase{
		{
			id:            "missing_auth_key",
			config:        "exit_node: n\nsubnet_router: n\nhostname: raspberrypi",
			expectedError: "auth key is required",
		},
		{
			id:            "subnet_host_bits_set",
			config:        "subnet_router: y\nsubnets: 192.168.1.1/24\nhostname: raspberrypi\nauth_key: tskey_1234",
			expectedError: "host bits set",
		},
		{
			id:            "invalid_tag",
			config:        "hostname: raspberrypi\ntags: \"server,prod_2\"\nauth_key: tskey_1234",
			expectedError: "invalid tag",
		},
	}

	for _, tc := range testCases {
		mockFS := afero.NewMemMapFs()

		configFile := "/tmp/config.yaml"
		_ = afero.WriteFile(mockFS, configFile, []byte(tc.config), 0644)

		userDataPath := "/mnt/sdcard/user-data"
		_ = mockFS.MkdirAll("/mnt/sdcard", 0755)
		_ = afero.WriteFile(mockFS, userDataPath, []byte("# cloud-config\n"), 0644)

		setupDeps := SetupCommand{
			Fsys:        mockFS,
			RootChecker: MockRootChecker{ShouldError: false},
			SystemSearcher: &file.SystemSearcher{
				Fsys:     mockFS,
				Searcher: &MockSearcher{},
			},
		}

		setupCmd := NewSetupCmd(setupDeps)

		var buf bytes.Buffer
		setupCmd.SetOut(&buf)
		setupCmd.SetErr(&buf)
		setupCmd.SetArgs([]string{"--file", configFile})

		err := setupCmd.Execute()

		var statusErr ierror.StatusError
		if !errors.As(err, &statusErr) || statusErr.StatusCode != 1 || !strings.Contains(err.Error(), tc.expectedError) {
			t.Errorf("%v: Execute() returned %v, expected StatusError with code 1 containing %q", tc.id, err, tc.expectedError)
		}

		contentBytes, _ := afero.ReadFile(mockFS, userDataPath)
		if content := string(contentBytes); content != "# cloud-config\n" {
			t.Errorf("%v: Expected user-data to be left unchanged, got:\n%s", tc.id, content)
		}
	}
}
