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

var tprInitTimeoutError = errgo.New("initialization timeout")

// IsTPRInitTimeout asserts tprInitTimeoutError.
func IsTPRInitTimeout(err error) bool {
	return errgo.Cause(err) == tprInitTimeoutError
}

var unexpectedlyShortResourceNameError = errgo.New("unexpectedly short resource name")

// IsUnexpectedlyShortResourceName asserts unexpectedlyShortResourceNameError.
func IsUnexpectedlyShortResourceName(err error) bool {
	return errgo.Cause(err) == unexpectedlyShortResourceNameError
}
