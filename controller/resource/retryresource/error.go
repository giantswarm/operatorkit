package retryresource

import (
	"github.com/giantswarm/microerror"
)

var incompatibleUnderlyingResourceError = microerror.New("incompatible underlying resource")

// isIncompatibleUnderlyingResource asserts incompatibleUnderlyingResourceError.
func isIncompatibleUnderlyingResource(err error) bool {
	return microerror.Cause(err) == incompatibleUnderlyingResourceError
}

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}
