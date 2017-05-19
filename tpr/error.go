package tpr

import (
	"github.com/juju/errgo"
)

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var tprInitTimeoutError = errgo.New("initialization timeout")

// IsTPRInitTimeout asserts tprInitTimeoutError.
func IsTPRInitTimeout(err error) bool {
	return errgo.Cause(err) == tprInitTimeoutError
}
