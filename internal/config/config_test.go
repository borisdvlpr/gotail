package config

import (
	"errors"
	"testing"

	ierr "github.com/borisdvlpr/gotail/internal/error"
)

type ValidateConfigTestCase struct {
	id            string
	config        Config
	expectedError error
}

func TestConfigValidate(t *testing.T) {
	testCases := []ValidateConfigTestCase{
		{
			id:            "case_01",
			config:        Config{ExitNode: "y", SubnetRouter: "n", Hostname: "test-host", AuthKey: "tskey-abcd1234"},
			expectedError: nil,
		},
		{
			id:            "case_02",
			config:        Config{ExitNode: "n", SubnetRouter: "y", Subnets: "192.168.1.0/24", Hostname: "subnet-router", AuthKey: "tskey-abcd1234"},
			expectedError: nil,
		},
		{
			id:            "case_03",
			config:        Config{ExitNode: "y", SubnetRouter: "n", Hostname: "test-host", AuthKey: ""},
			expectedError: ierr.StatusError{Status: "auth key is required", StatusCode: 1},
		},
		{
			id:            "case_04",
			config:        Config{ExitNode: "y", SubnetRouter: "n", Hostname: "", AuthKey: "tskey-abcd1234"},
			expectedError: ierr.StatusError{Status: "hostname is required", StatusCode: 1},
		},
		{
			id:            "case_05",
			config:        Config{ExitNode: "n", SubnetRouter: "y", Subnets: "", Hostname: "subnet-router", AuthKey: "tskey-abcd1234"},
			expectedError: ierr.StatusError{Status: "subnets are required when subnet router is enabled", StatusCode: 1},
		},
		{
			id:            "case_06",
			config:        Config{ExitNode: "n", SubnetRouter: "n", Hostname: "simple-node", AuthKey: "tskey-abcd1234"},
			expectedError: nil,
		},
		{
			id:            "case_07",
			config:        Config{ExitNode: "y", SubnetRouter: "y", Subnets: "192.168.1.0/24,10.0.0.0/16", Hostname: "fully-configured-node", AuthKey: "tskey-abcdefghijklmnop"},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		err := tc.config.Validate()

		if tc.expectedError == nil && err != nil {
			t.Errorf("%v: Validate() returned error %v, expected no error", tc.id, err)
		}

		if tc.expectedError != nil && err == nil {
			t.Errorf("%v: Validate() returned no error, expected error containing %q", tc.id, tc.expectedError)
		}

		if tc.expectedError != nil && !errors.Is(err, tc.expectedError) {
			t.Errorf("%v: Validate() returned error %q, expected %q", tc.id, err.Error(), tc.expectedError)
		}
	}
}
