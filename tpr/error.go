package tpr

import (
	"github.com/giantswarm/microerror"
)

var alreadyExistsError = microerror.New("already exists")

// IsAlreadyExists asserts alreadyExistsError.
func IsAlreadyExists(err error) bool {
	return microerror.Cause(err) == alreadyExistsError
}

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var malformedNameError = microerror.New("malformed name")

// IsMalformedName asserts malformedNameError.
func IsMalformedName(err error) bool {
	return microerror.Cause(err) == malformedNameError
}

var tprInitTimeoutError = microerror.New("initialization timeout")

// IsTPRInitTimeout asserts tprInitTimeoutError.
func IsTPRInitTimeout(err error) bool {
	return microerror.Cause(err) == tprInitTimeoutError
}
