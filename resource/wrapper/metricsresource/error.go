package metricsresource

import (
	"github.com/giantswarm/microerror"
)

var incompatibleUnderlyingResourceError = &microerror.Error{
	Kind: "incompatibleUnderlyingResourceError",
}

// isIncompatibleUnderlyingResource asserts incompatibleUnderlyingResourceError.
func isIncompatibleUnderlyingResource(err error) bool {
	return microerror.Cause(err) == incompatibleUnderlyingResourceError
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}
