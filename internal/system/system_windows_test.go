//go:build windows

package system

//go:build windows

package system

import (
    "errors"
    "testing"

    ierror "github.com/borisdvlpr/gotail/internal/error"
)

type MockElevationChecker struct {
    elevated bool
    err      error
}

func (m MockElevationChecker) isElevated() (bool, error) {
    return m.elevated, m.err
}

type CheckRootTestCase struct {
    id            string
    checker       MockElevationChecker
    expectedError error
}
