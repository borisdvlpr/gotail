package system

import (
	"fmt"
)

type MockRootChecker struct {
	shouldError bool
	errorMsg    string
}

func (m MockRootChecker) CheckRoot() error {
	if m.shouldError {
		return fmt.Errorf("%v", m.errorMsg)
	}
	return nil
}

type MockRootCheckerTestCase struct {
	id            string
	checker       MockRootChecker
	expectedError error
}
