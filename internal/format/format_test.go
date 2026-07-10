package format

import (
	"reflect"
	"testing"

	"go.yaml.in/yaml/v3"
)

type BuildRunCmdTestCase struct {
	id       string
	input    []string
	expected string
}

func TestBuildRunCmd(t *testing.T) {
	testCases := []BuildRunCmdTestCase{
		{
			id:       "typical_command",
			input:    []string{"tailscale", "up", "--ssh"},
			expected: `["tailscale", "up", "--ssh"]`,
		},
		{
			id:       "single_element",
			input:    []string{"sh"},
			expected: `["sh"]`,
		},
		{
			id:       "comma_in_value_preserved",
			input:    []string{"--advertise-routes=192.168.1.0/24,10.0.0.0/16"},
			expected: `["--advertise-routes=192.168.1.0/24,10.0.0.0/16"]`,
		},
		{
			id:       "embedded_double_quote_escaped",
			input:    []string{`say "hi"`},
			expected: `["say \"hi\""]`,
		},
		{
			id:       "embedded_backslash_escaped",
			input:    []string{`path\to\thing`},
			expected: `["path\\to\\thing"]`,
		},
		{
			id:       "backslash_and_quote_escaped",
			input:    []string{`weird\"mix`},
			expected: `["weird\\\"mix"]`,
		},
		{
			id:       "empty_string_element",
			input:    []string{""},
			expected: `[""]`,
		},
		{
			id:       "empty_slice",
			input:    []string{},
			expected: `[]`,
		},
		{
			id:       "nil_slice",
			input:    nil,
			expected: `[]`,
		},
	}

	for _, tc := range testCases {
		output := BuildRunCmd(tc.input)

		if output != tc.expected {
			t.Errorf("%v: BuildRunCmd(%q) = %q, expected %q", tc.id, tc.input, output, tc.expected)
		}
	}
}

// TestBuildRunCmd_RoundTrip asserts the real contract: the output is a valid
// YAML flow sequence that parses back to the exact input args. This is what
// protects against the comma-splitting bug and catches any escaping error,
// independent of the exact spacing/quoting style.
func TestBuildRunCmd_RoundTrip(t *testing.T) {
	inputs := [][]string{
		{"tailscale", "up", "--ssh", "--advertise-routes=192.168.1.0/24,10.0.0.0/16", "--authkey=tskey-abcd1234"},
		{"sh", "-c", "curl -fsSL https://tailscale.com/install.sh | sh"},
		{"sh", "-c", "echo 'net.ipv4.ip_forward = 1' | sudo tee -a /etc/sysctl.d/99-tailscale.conf"},
		{`arg with "quotes"`, `arg\with\backslashes`, "value,with,commas"},
	}

	for i, args := range inputs {
		out := BuildRunCmd(args)

		var parsed []string
		if err := yaml.Unmarshal([]byte(out), &parsed); err != nil {
			t.Errorf("case %d: output %q is not valid YAML: %v", i, out, err)
			continue
		}

		if !reflect.DeepEqual(parsed, args) {
			t.Errorf("case %d: round-trip mismatch\n  input:  %q\n  output: %s\n  parsed: %q", i, args, out, parsed)
		}
	}
}
