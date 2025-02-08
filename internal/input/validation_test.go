package input

import (
	"errors"
	"testing"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

type ValidateSubnetsTestCase struct {
	id            string
	subnet        string
	expectedError error
}

func TestValidateSubnet(t *testing.T) {
	testCases := []ValidateSubnetsTestCase{
		{id: "case_01", subnet: "192.168.1.1/24", expectedError: nil},
		{id: "case_02", subnet: "192.168.1.1/24,192.168.2.2/24", expectedError: nil},
		{id: "case_03", subnet: "2001:db8::/32", expectedError: nil},
		{id: "case_04", subnet: "2001:db8::/32,2001:db8::/32", expectedError: nil},
		{id: "case_05", subnet: "", expectedError: ierror.StatusError{Status: ": invalid subnet format", StatusCode: 1}},
		{id: "case_06", subnet: "192.168.1.1", expectedError: ierror.StatusError{Status: "192.168.1.1: invalid subnet format", StatusCode: 1}},
		{id: "case_07", subnet: "192.168.1.1/24,", expectedError: ierror.StatusError{Status: ": invalid subnet format", StatusCode: 1}},
		{id: "case_08", subnet: ",192.168.1.1", expectedError: ierror.StatusError{Status: ": invalid subnet format", StatusCode: 1}},
		{id: "case_09", subnet: "192.168.1.1/24,192.168.2.2", expectedError: ierror.StatusError{Status: "192.168.2.2: invalid subnet format", StatusCode: 1}},
		{id: "case_10", subnet: "192.168.1.", expectedError: ierror.StatusError{Status: "192.168.1.: invalid subnet format", StatusCode: 1}},
		{id: "case_11", subnet: "2001:db8::", expectedError: ierror.StatusError{Status: "2001:db8::: invalid subnet format", StatusCode: 1}},
		{id: "case_12", subnet: "2001:db8::/32,2001:db8::", expectedError: ierror.StatusError{Status: "2001:db8::: invalid subnet format", StatusCode: 1}},
	}

	for _, tc := range testCases {
		err := ValidateSubnets(tc.subnet)
		if !errors.Is(err, tc.expectedError) {
			t.Errorf("%v: ValidateSubnets() error = %v, wantErr %v", tc.id, err, tc.expectedError)
		}
	}
}
