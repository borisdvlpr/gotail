package input

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

var tagRegexp = regexp.MustCompile(`^[a-z0-9-]+$`)

// ValidateSubnets validates a comma-separated list of subnets.
// Each subnet must be in CIDR notation. If any subnet is not in
// the correct format, an error is returned indicating the invalid
// subnet and a status code of 1.
func ValidateSubnets(subnets string) error {
	for subnet := range strings.SplitSeq(subnets, ",") {
		if _, _, err := net.ParseCIDR(strings.TrimSpace(subnet)); err != nil {
			return ierror.StatusError{Status: fmt.Sprintf("%s: invalid subnet format", subnet), StatusCode: 1}
		}
	}
	return nil
}

// ValidateTags validates and normalizes a comma-separated list of Tailscale
// tag names. Users provide bare names (e.g. "server,prod-2"); each must contain
// only lowercase letters, digits, and hyphens. Tailscale lowercases tag names,
// so input is trimmed and lowercased before validation, and the "tag:" prefix
// is added afterward. Returns the normalized, comma-separated tags (e.g.
// "tag:server,tag:prod-2"), or a StatusError describing the first invalid tag.
func ValidateTags(tags string) (string, error) {
	var normalized []string
	for tag := range strings.SplitSeq(tags, ",") {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if !tagRegexp.MatchString(tag) {
			return "", ierror.StatusError{
				Status:     fmt.Sprintf("%q: invalid tag (expected letters, digits, and hyphens)", tag),
				StatusCode: 1,
			}
		}

		normalized = append(normalized, "tag:"+tag)
	}

	return strings.Join(normalized, ","), nil
}
