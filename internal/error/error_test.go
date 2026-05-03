package error

import (
	"testing"
)

type StatusErrorTestCase struct {
	id       string
	err      StatusError
	expected string
}

func TestStatusError(t *testing.T) {
	testCases := []StatusErrorTestCase{
		{
			id:       "with_status",
			err:      StatusError{Status: "file not found", StatusCode: 2},
			expected: "file not found",
		},
		{
			id:       "empty_status",
			err:      StatusError{StatusCode: 2},
			expected: "abort",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			if got := tc.err.Error(); got != tc.expected {
				t.Errorf("StatusError.Error() = %q, want %q", got, tc.expected)
			}
		})
	}
}
