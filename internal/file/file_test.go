package file

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	ierror "github.com/borisdvlpr/gotail/internal/error"
	"github.com/spf13/afero"
)

// MockUnreadableDirFs wraps an afero.Fs so that opening the unreadable
// directory fails with a permission error, simulating a mounted volume the
// current user is not allowed to list.
type MockUnreadableDirFs struct {
	afero.Fs
	unreadable string
}

func (m MockUnreadableDirFs) Open(name string) (afero.File, error) {
	if name == m.unreadable {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrPermission}
	}
	return m.Fs.Open(name)
}

type GetFilePathTestCase struct {
	id           string
	file         string
	expectedPath string
}

func TestGetFilePath(t *testing.T) {
	testCases := []GetFilePathTestCase{
		{
			id:           "valid_path",
			file:         "user-data",
			expectedPath: filepath.Join("file-test-dir", "user-data"),
		},
		{
			id:           "empty_path",
			file:         "",
			expectedPath: "",
		},
	}

	for _, tc := range testCases {
		fs := afero.NewMemMapFs()

		tempDir := filepath.Join(".", "file-test-dir")
		err := fs.MkdirAll(tempDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test dir: %v", err)
		}

		testFilePath := filepath.Join(tempDir, tc.file)
		if tc.file != "" {
			if _, err := fs.Create(testFilePath); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
		}

		// Pass the filesystem to GetFilePath
		path, err := GetFilePath(fs, tempDir, "user-data")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if path != tc.expectedPath {
			t.Errorf("%v: GetFilePath() path = %v, wantPath %v", tc.id, path, tc.expectedPath)
		}
	}
}

type SearchVolumesTestCase struct {
	id           string
	unreadable   string
	files        []string
	expectedPath string
	expectError  bool
}

func TestSearchVolumes(t *testing.T) {
	volumesDir := filepath.FromSlash("/Volumes")

	testCases := []SearchVolumesTestCase{
		{
			id:           "found_in_volume",
			files:        []string{filepath.Join(volumesDir, "bootfs", "user-data")},
			expectedPath: filepath.Join(volumesDir, "bootfs", "user-data"),
		},
		{
			// "Backup" sorts before "bootfs", so it is searched first.
			id:           "unreadable_volume_skipped",
			unreadable:   filepath.Join(volumesDir, "Backup"),
			files:        []string{filepath.Join(volumesDir, "Backup", "notes.txt"), filepath.Join(volumesDir, "bootfs", "user-data")},
			expectedPath: filepath.Join(volumesDir, "bootfs", "user-data"),
		},
		{
			id:           "not_found",
			files:        []string{filepath.Join(volumesDir, "Backup", "notes.txt")},
			expectedPath: "",
		},
		{
			id:          "volumes_dir_missing",
			files:       nil,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		memFs := afero.NewMemMapFs()
		for _, f := range tc.files {
			if err := afero.WriteFile(memFs, f, []byte{}, 0644); err != nil {
				t.Fatalf("%v: failed to create test file: %v", tc.id, err)
			}
		}

		fsys := MockUnreadableDirFs{Fs: memFs, unreadable: tc.unreadable}

		path, err := searchVolumes(fsys, volumesDir, "user-data")

		if tc.expectError {
			if err == nil {
				t.Errorf("%v: searchVolumes() returned no error, expected one", tc.id)
			}
			continue
		}

		if err != nil {
			t.Errorf("%v: searchVolumes() returned error %v, expected none", tc.id, err)
		}

		if path != tc.expectedPath {
			t.Errorf("%v: searchVolumes() path = %q, expected %q", tc.id, path, tc.expectedPath)
		}
	}
}

// MockDriveSearcher returns a fixed result from Search.
type MockDriveSearcher struct {
	path string
	err  error
}

func (m MockDriveSearcher) Search(ctx context.Context, fsys afero.Fs, fileName string) (string, error) {
	return m.path, m.err
}

type FindUserDataTestCase struct {
	id           string
	searcher     DriveSearcher
	expectedPath string
	expectedErr  error
	expectedCode int
}

func TestFindUserData(t *testing.T) {
	searchErr := errors.New("lsblk failed")

	testCases := []FindUserDataTestCase{
		{
			id:           "found",
			searcher:     MockDriveSearcher{path: "/media/user/bootfs/user-data"},
			expectedPath: "/media/user/bootfs/user-data",
		},
		{
			id:           "unsupported_os",
			searcher:     nil,
			expectedCode: 71,
		},
		{
			id:           "not_found",
			searcher:     MockDriveSearcher{path: ""},
			expectedCode: 2,
		},
		{
			id:          "search_error",
			searcher:    MockDriveSearcher{err: searchErr},
			expectedErr: searchErr,
		},
	}

	for _, tc := range testCases {
		s := &SystemSearcher{Fsys: afero.NewMemMapFs(), Searcher: tc.searcher}
		path, err := s.FindUserData()

		if path != tc.expectedPath {
			t.Errorf("%v: FindUserData() path = %q, expected %q", tc.id, path, tc.expectedPath)
		}

		switch {
		case tc.expectedCode != 0:
			var statusErr ierror.StatusError
			if !errors.As(err, &statusErr) || statusErr.StatusCode != tc.expectedCode {
				t.Errorf("%v: FindUserData() error = %v, expected StatusError with code %d", tc.id, err, tc.expectedCode)
			}
		case tc.expectedErr != nil:
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("%v: FindUserData() error = %v, expected %v", tc.id, err, tc.expectedErr)
			}
		default:
			if err != nil {
				t.Errorf("%v: FindUserData() returned error %v, expected none", tc.id, err)
			}
		}
	}
}

type DrainSearchChanTestCase struct {
	id             string
	results        []SearchResult
	expectedPath   string
	expectError    bool
	expectedCancel bool
}

func TestDrainSearchChan(t *testing.T) {
	testCases := []DrainSearchChanTestCase{
		{
			id:             "first_hit_returned",
			results:        []SearchResult{{Path: ""}, {Path: "/mnt/sdcard/user-data"}, {Path: "/mnt/other/user-data"}},
			expectedPath:   "/mnt/sdcard/user-data",
			expectedCancel: true,
		},
		{
			id:             "error_returned",
			results:        []SearchResult{{Err: errors.New("read failed")}},
			expectError:    true,
			expectedCancel: true,
		},
		{
			id:           "closed_without_hit",
			results:      []SearchResult{{Path: ""}, {Path: ""}},
			expectedPath: "",
		},
	}

	for _, tc := range testCases {
		c := make(chan SearchResult, len(tc.results))
		for _, r := range tc.results {
			c <- r
		}
		close(c)

		ctx, cancel := context.WithCancel(context.Background())
		path, err := drainSearchChan(cancel, c)

		if tc.expectError != (err != nil) {
			t.Errorf("%v: drainSearchChan() error = %v, expectError %v", tc.id, err, tc.expectError)
		}

		if path != tc.expectedPath {
			t.Errorf("%v: drainSearchChan() path = %q, expected %q", tc.id, path, tc.expectedPath)
		}

		if cancelled := ctx.Err() != nil; cancelled != tc.expectedCancel {
			t.Errorf("%v: drainSearchChan() cancelled = %v, expected %v", tc.id, cancelled, tc.expectedCancel)
		}
		cancel()
	}
}
