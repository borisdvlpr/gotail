package file

import (
	"os"
	"path/filepath"
	"testing"
)

type GetFilePathTestCase struct {
	id           string
	file         string
	expectedPath string
}

func TestGetFilePath(t *testing.T) {
	testCases := []GetFilePathTestCase{
		{id: "case_01", file: "user-data", expectedPath: "file-test-dir/user-data"},
		{id: "case_02", file: "", expectedPath: ""},
	}

	for _, tc := range testCases {
		tempDir := filepath.Join(".", "file-test-dir")
		err := os.MkdirAll(tempDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test dir: %v", err)
		}

		testFilePath := filepath.Join(tempDir, tc.file)
		if tc.file != "" {
			if _, err := os.Create(testFilePath); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
		}

		path, err := GetFilePath(tempDir, "user-data")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if path != tc.expectedPath {
			t.Errorf("%v: GetFilePath() path = %v, wantPath %v", tc.id, err, tc.expectedPath)
		}

		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
}
