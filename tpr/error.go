package tpr

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

var malformedNameError = errgo.New("malformed name")

// IsMalformedName asserts malformedNameError.
func IsMalformedName(err error) bool {
	return errgo.Cause(err) == malformedNameError
}

var tprInitTimeoutError = errgo.New("initialization timeout")

// IsTPRInitTimeout asserts tprInitTimeoutError.
func IsTPRInitTimeout(err error) bool {
	return errgo.Cause(err) == tprInitTimeoutError
}
