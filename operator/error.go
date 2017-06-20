package operator

import (
	"github.com/juju/errgo"
)

var executionFailedError = errgo.New("execution failed")

// IsExecutionFailed asserts executionFailedError.
func IsExecutionFailed(err error) bool {
	return errgo.Cause(err) == executionFailedError
}
