//go:build !windows

package system

import (
	"errors"
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
			id:            "root_check_passes",
			checker:       MockRootChecker{shouldError: false},
			expectedError: nil,
		},
		{
			id:            "root_check_fails_with_message",
			checker:       MockRootChecker{shouldError: true, errorMsg: "error: permission denied"},
			expectedError: ierror.StatusError{Status: "error: permission denied", StatusCode: 1},
		},
		{
			id:            "root_check_fails_without_message",
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

		if tc.expectedError != nil && err != nil && !errors.Is(err, tc.expectedError) {
			t.Errorf("%v: TestMockRootChecker() returned error %q, expected %q", tc.id, err, tc.expectedError)
		}
	}
}
