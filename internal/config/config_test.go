package config

import (
	"errors"
	"testing"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

type ValidateConfigTestCase struct {
	id            string
	config        Config
	expectedError error
}

func TestConfigValidate(t *testing.T) {
	testCases := []ValidateConfigTestCase{
		{
			id:            "exit_node_valid",
			config:        Config{ExitNode: "y", SubnetRouter: "n", Hostname: "test-host", AuthKey: "tskey-abcd1234"},
			expectedError: nil,
		},
		{
			id:            "subnet_router_valid",
			config:        Config{ExitNode: "n", SubnetRouter: "y", Subnets: "192.168.1.0/24", Hostname: "subnet-router", AuthKey: "tskey-abcd1234"},
			expectedError: nil,
		},
		{
			id:            "missing_auth_key",
			config:        Config{ExitNode: "y", SubnetRouter: "n", Hostname: "test-host", AuthKey: ""},
			expectedError: ierror.StatusError{Status: "auth key is required", StatusCode: 1},
		},
		{
			id:            "missing_hostname",
			config:        Config{ExitNode: "y", SubnetRouter: "n", Hostname: "", AuthKey: "tskey-abcd1234"},
			expectedError: ierror.StatusError{Status: "hostname is required", StatusCode: 1},
		},
		{
			id:            "subnet_router_without_subnets",
			config:        Config{ExitNode: "n", SubnetRouter: "y", Subnets: "", Hostname: "subnet-router", AuthKey: "tskey-abcd1234"},
			expectedError: ierror.StatusError{Status: "subnets are required when subnet router is enabled", StatusCode: 1},
		},
		{
			id:            "simple_node_valid",
			config:        Config{ExitNode: "n", SubnetRouter: "n", Hostname: "simple-node", AuthKey: "tskey-abcd1234"},
			expectedError: nil,
		},
		{
			id:            "fully_configured_valid",
			config:        Config{ExitNode: "y", SubnetRouter: "y", Subnets: "192.168.1.0/24,10.0.0.0/16", Hostname: "fully-configured-node", AuthKey: "tskey-abcdefghijklmnop"},
			expectedError: nil,
		},
		{
			id:            "tags_valid",
			config:        Config{ExitNode: "n", SubnetRouter: "n", Hostname: "tagged-node", Tags: "server,prod", AuthKey: "tskey-abcd1234"},
			expectedError: nil,
		},
		{
			id:            "tags_do_not_mask_missing_auth_key",
			config:        Config{ExitNode: "n", SubnetRouter: "n", Hostname: "tagged-node", Tags: "server", AuthKey: ""},
			expectedError: ierror.StatusError{Status: "auth key is required", StatusCode: 1},
		},
		{
			id:            "tags_with_full_config_valid",
			config:        Config{ExitNode: "y", SubnetRouter: "y", Subnets: "192.168.1.0/24,10.0.0.0/16", Hostname: "full-node", Tags: "exit,router", AuthKey: "tskey-abcd1234"},
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

		if tc.expectedError != nil && err != nil && !errors.Is(err, tc.expectedError) {
			t.Errorf("%v: Validate() returned error %q, expected %q", tc.id, err, tc.expectedError)
		}
	}
}
