package file

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

type GetFilePathTestCase struct {
	id           string
	file         string
	expectedPath string
}

func TestGetFilePath(t *testing.T) {
	testCases := []GetFilePathTestCase{
		{
			id:           "case_01",
			file:         "user-data",
			expectedPath: "file-test-dir/user-data",
		},
		{
			id:           "case_02",
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
