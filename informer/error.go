package informer

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidEventError = microerror.New("invalid event")

// IsInvalidEvent asserts invalidEventError.
func IsInvalidEvent(err error) bool {
	return microerror.Cause(err) == invalidEventError
}
