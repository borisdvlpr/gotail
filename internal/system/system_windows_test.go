//go:build windows

package system

import (
	"errors"
	"testing"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

type MockElevationChecker struct {
	elevated bool
	err      error
}

func (m MockElevationChecker) isElevated() (bool, error) {
	return m.elevated, m.err
}

type CheckRootTestCase struct {
	id            string
	checker       MockElevationChecker
	expectedError error
}

func TestCheckRootWithChecker(t *testing.T) {
	testCases := []CheckRootTestCase{
		{
			id:            "elevated_passes",
			checker:       MockElevationChecker{elevated: true, err: nil},
			expectedError: nil,
		},
		{
			id:      "not_elevated_requires_admin",
			checker: MockElevationChecker{elevated: false, err: nil},
			expectedError: ierror.StatusError{
				Status:     "setup must be run as administrator. please relaunch from an elevated shell (run as administrator).",
				StatusCode: 126,
			},
		},
		{
			id:            "elevation_check_errors",
			checker:       MockElevationChecker{elevated: false, err: errors.New("access denied")},
			expectedError: ierror.StatusError{Status: "access denied", StatusCode: 1},
		},
	}

	for _, tc := range testCases {
		err := checkRootWithChecker(tc.checker)

		if tc.expectedError == nil && err != nil {
			t.Errorf("%v: checkRootWithChecker() returned error %v, expected no error", tc.id, err)
		}

		if tc.expectedError != nil && err == nil {
			t.Errorf("%v: checkRootWithChecker() returned no error, expected %q", tc.id, tc.expectedError)
		}

		if tc.expectedError != nil && err != nil && !errors.Is(err, tc.expectedError) {
			t.Errorf("%v: checkRootWithChecker() returned %q, expected %q", tc.id, err, tc.expectedError)
		}
	}
}
