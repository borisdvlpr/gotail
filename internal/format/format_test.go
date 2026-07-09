package format

import (
	"testing"
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
