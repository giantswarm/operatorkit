package transaction

import (
	"github.com/juju/errgo"
)

var alreadyExistsError = errgo.New("already exists")

// IsAlreadyExists asserts alreadyExistsError.
func IsAlreadyExists(err error) bool {
	return errgo.Cause(err) == alreadyExistsError
}

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var invalidExecutionError = errgo.New("invalid execution")

// IsInvalidExecution asserts invalidExecutionError.
func IsInvalidExecution(err error) bool {
	return errgo.Cause(err) == invalidExecutionError
}

var notFoundError = errgo.New("not found")

// IsNotFound asserts notFoundError.
func IsNotFound(err error) bool {
	return errgo.Cause(err) == notFoundError
}
