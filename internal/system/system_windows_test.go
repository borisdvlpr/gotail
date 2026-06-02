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
			id:            "case_01",
			checker:       MockElevationChecker{elevated: true, err: nil},
			expectedError: nil,
		},
		{
			id:      "case_02",
			checker: MockElevationChecker{elevated: false, err: nil},
			expectedError: ierror.StatusError{
				Status:     "setup must be run as administrator. please relaunch from an elevated shell (run as administrator).",
				StatusCode: 126,
			},
		},
		{
			id:            "case_03",
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
