package input

import (
	"fmt"
	"net"
	"strings"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

// ValidateSubnets validates a comma-separated list of subnets.
// Each subnet must be in CIDR notation. If any subnet is not in
// the correct format, an error is returned indicating the invalid
// subnet and a status code of 1.
func ValidateSubnets(subnets string) error {
	for _, subnet := range strings.Split(subnets, ",") {
		if _, _, err := net.ParseCIDR(subnet); err != nil {
			return ierror.StatusError{Status: fmt.Sprintf("%s: invalid subnet format", subnet), StatusCode: 1}
		}
	}
	return nil
}
