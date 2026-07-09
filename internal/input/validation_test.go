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

type ValidateTagsTestCase struct {
	id             string
	input          string
	expectedOutput string
	expectError    bool
}

func TestValidateSubnet(t *testing.T) {
	testCases := []ValidateSubnetsTestCase{
		{
			id:            "case_01",
			subnet:        "192.168.1.1/24",
			expectedError: nil,
		},
		{
			id:            "case_02",
			subnet:        "192.168.1.1/24,192.168.2.2/24",
			expectedError: nil,
		},
		{
			id:            "case_03",
			subnet:        "2001:db8::/32",
			expectedError: nil,
		},
		{
			id:            "case_04",
			subnet:        "2001:db8::/32,2001:db8::/32",
			expectedError: nil,
		},
		{
			id:            "case_05",
			subnet:        "",
			expectedError: ierror.StatusError{Status: ": invalid subnet format", StatusCode: 1},
		},
		{
			id:            "case_06",
			subnet:        "192.168.1.1",
			expectedError: ierror.StatusError{Status: "192.168.1.1: invalid subnet format", StatusCode: 1},
		},
		{
			id:            "case_07",
			subnet:        "192.168.1.1/24,",
			expectedError: ierror.StatusError{Status: ": invalid subnet format", StatusCode: 1},
		},
		{
			id:            "case_08",
			subnet:        ",192.168.1.1",
			expectedError: ierror.StatusError{Status: ": invalid subnet format", StatusCode: 1},
		},
		{
			id:            "case_09",
			subnet:        "192.168.1.1/24,192.168.2.2",
			expectedError: ierror.StatusError{Status: "192.168.2.2: invalid subnet format", StatusCode: 1},
		},
		{
			id:            "case_10",
			subnet:        "192.168.1.",
			expectedError: ierror.StatusError{Status: "192.168.1.: invalid subnet format", StatusCode: 1},
		},
		{
			id:            "case_11",
			subnet:        "2001:db8::",
			expectedError: ierror.StatusError{Status: "2001:db8::: invalid subnet format", StatusCode: 1},
		},
		{
			id:            "case_12",
			subnet:        "2001:db8::/32,2001:db8::",
			expectedError: ierror.StatusError{Status: "2001:db8::: invalid subnet format", StatusCode: 1},
		},
	}

	for _, tc := range testCases {
		err := ValidateSubnets(tc.subnet)
		if !errors.Is(err, tc.expectedError) {
			t.Errorf("%v: ValidateSubnets() error = %q, wantErr %q", tc.id, err, tc.expectedError)
		}
	}
}

func TestValidateTags(t *testing.T) {
	testCases := []ValidateTagsTestCase{
		{
			id:             "single_tag",
			input:          "server",
			expectedOutput: "tag:server",
			expectError:    false,
		},
		{
			id:             "multiple_tags",
			input:          "server,prod",
			expectedOutput: "tag:server,tag:prod",
			expectError:    false,
		},
		{
			id:             "hyphenated_name",
			input:          "prod-2",
			expectedOutput: "tag:prod-2",
			expectError:    false},
		{
			id:             "uppercase_normalized",
			input:          "Server",
			expectedOutput: "tag:server",
			expectError:    false},
		{
			id:             "whitespace_trimmed",
			input:          " server , prod ",
			expectedOutput: "tag:server,tag:prod",
			expectError:    false},
		{
			id:          "prefix_included",
			input:       "tag:server",
			expectError: true},
		{
			id:          "space_in_name",
			input:       "my server",
			expectError: true},
		{
			id:          "underscore",
			input:       "prod_env",
			expectError: true},
		{
			id:          "special_char",
			input:       "server!",
			expectError: true},
		{
			id:          "trailing_comma",
			input:       "server,",
			expectError: true},
		{
			id:          "whitespace_only",
			input:       "   ",
			expectError: true},
	}

	for _, tc := range testCases {
		output, err := ValidateTags(tc.input)

		if tc.expectError {
			if err == nil {
				t.Errorf("%v: ValidateTags(%q) returned no error, expected one", tc.id, tc.input)
				continue
			}

			var statusErr ierror.StatusError
			if !errors.As(err, &statusErr) {
				t.Errorf("%v: ValidateTags(%q) returned %T, expected ierr.StatusError", tc.id, tc.input, err)

			} else if statusErr.StatusCode != 1 {
				t.Errorf("%v: ValidateTags(%q) returned status code %d, expected 1", tc.id, tc.input, statusErr.StatusCode)
			}

			continue
		}

		if err != nil {
			t.Errorf("%v: ValidateTags(%q) returned error %v, expected none", tc.id, tc.input, err)
			continue
		}

		if output != tc.expectedOutput {
			t.Errorf("%v: ValidateTags(%q) = %q, expected %q", tc.id, tc.input, output, tc.expectedOutput)
		}
	}
}
