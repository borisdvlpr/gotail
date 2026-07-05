package format

import "strings"

// BuildRunCmd renders args as a YAML flow sequence with each element
// double-quoted, e.g. ["tailscale", "up", "--advertise-routes=10.0.0.0/24,10.0.1.0/24"].
// Quoting is required: cloud-init parses runcmd entries as YAML, so an unquoted
// element containing a comma (multiple subnets) would be split into separate
// arguments by the YAML parser.
func BuildRunCmd(args []string) string {
	quoted := make([]string, len(args))
	for i, a := range args {
		esc := strings.ReplaceAll(a, `\`, `\\`)
		esc = strings.ReplaceAll(esc, `"`, `\"`)
		quoted[i] = `"` + esc + `"`
	}

	return "[" + strings.Join(quoted, ", ") + "]"
}
