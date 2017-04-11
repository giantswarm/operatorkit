package daemon

import (
	"github.com/juju/errgo"
)

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var invalidFlagError = errgo.New("invalid flag")

// IsInvalidFlag asserts invalidFlagError.
func IsInvalidFlag(err error) bool {
	return errgo.Cause(err) == invalidFlagError
}
