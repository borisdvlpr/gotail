package input

import (
	"fmt"
	"net"
	"strings"

	ierror "github.com/borisdvlpr/gotail/internal/error"
)

// ValidateSubnets validates
func ValidateSubnets(subnets string) error {
	for _, subnet := range strings.Split(subnets, ",") {
		if _, _, err := net.ParseCIDR(subnet); err != nil {
			return ierror.StatusError{Status: fmt.Sprintf("%s: invalid subnet format", subnet), StatusCode: 1}
		}
	}
	return nil
}
