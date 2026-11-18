package input

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

var (
	tagRegexp      = regexp.MustCompile(`^[a-z0-9-]+$`)
	hostnameRegexp = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9-]{0,61}[A-Za-z0-9])?$`)
)

// ValidateHostname validates a device hostname as a single RFC 1123 label:
// 1-63 letters, digits, and hyphens, not starting or ending with a hyphen.
// The hostname is placed in a shell command run as root on first boot, so
// anything outside that set is rejected. Returns a StatusError with status
// code 1 if the hostname is invalid.
func ValidateHostname(hostname string) error {
	if !hostnameRegexp.MatchString(hostname) {
		return ierror.StatusError{
			Status:     fmt.Sprintf("%q: invalid hostname (expected 1-63 letters, digits, and hyphens, not starting or ending with a hyphen)", hostname),
			StatusCode: 1,
		}
	}

	return nil
}

// ValidateSubnets validates and normalizes a comma-separated list of subnets.
// Each subnet must be in CIDR notation with no host bits set (e.g.
// "192.168.1.0/24", not "192.168.1.1/24"), since `tailscale up` rejects
// those. Whitespace around each subnet is trimmed, because Tailscale does not
// trim it. Returns the normalized, comma-separated subnets, or a StatusError
// with status code 1 describing the first invalid subnet.
func ValidateSubnets(subnets string) (string, error) {
	var normalized []string
	for subnet := range strings.SplitSeq(subnets, ",") {
		trimmed := strings.TrimSpace(subnet)
		ip, ipNet, err := net.ParseCIDR(trimmed)
		if err != nil {
			return "", ierror.StatusError{Status: fmt.Sprintf("%s: invalid subnet format", subnet), StatusCode: 1}
		}

		if !ip.Equal(ipNet.IP) {
			return "", ierror.StatusError{Status: fmt.Sprintf("%s: host bits set, expected %s", trimmed, ipNet), StatusCode: 1}
		}

		normalized = append(normalized, trimmed)
	}

	return strings.Join(normalized, ","), nil
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
