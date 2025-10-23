package system

import (
	"errors"
	"os"
	"testing"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

type MockRootChecker struct {
	shouldError bool
	errorMsg    string
}

func (m MockRootChecker) CheckRoot() error {
	if m.shouldError {
		return ierror.StatusError{Status: m.errorMsg, StatusCode: 1}
	}

	execPath := "fake/path/to/executable"
	args := []string{"sudo", execPath}
	args = append(args, os.Args[1:]...)

	return nil
}

type MockRootCheckerTestCase struct {
	id            string
	checker       MockRootChecker
	expectedError error
}

func TestMockRootChecker(t *testing.T) {
	testCases := []MockRootCheckerTestCase{
		{
			id:            "case_01",
			checker:       MockRootChecker{shouldError: false},
			expectedError: nil,
		},
		{
			id:            "case_02",
			checker:       MockRootChecker{shouldError: true, errorMsg: "error: permission denied"},
			expectedError: ierror.StatusError{Status: "error: permission denied", StatusCode: 1},
		},
		{
			id:            "case_03",
			checker:       MockRootChecker{shouldError: true, errorMsg: ""},
			expectedError: ierror.StatusError{Status: "", StatusCode: 1},
		},
	}

	for _, tc := range testCases {
		err := tc.checker.CheckRoot()

		if tc.expectedError == nil && err != nil {
			t.Errorf("%v: TestMockRootChecker() returned error %v, expected no error", tc.id, err)
		}

		if tc.expectedError != nil && err == nil {
			t.Errorf("%v: TestMockRootChecker() returned no error, expected error containing %q", tc.id, tc.expectedError)
		}

		if tc.expectedError != nil && !errors.Is(err, tc.expectedError) {
			t.Errorf("%v: TestMockRootChecker() returned error %q, expected %q", tc.id, err, tc.expectedError)
		}
	}
}
