package controller

import (
	"github.com/giantswarm/microerror"
)

var executionFailedError = microerror.New("execution failed")

// IsExecutionFailed asserts executionFailedError.
func IsExecutionFailed(err error) bool {
	return microerror.Cause(err) == executionFailedError
}

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidRESTClientError = microerror.New("invalid REST client")

// IsInvalidRESTClient asserts invalidRESTClientError.
func IsInvalidRESTClient(err error) bool {
	return microerror.Cause(err) == invalidRESTClientError
}

var noResourceSetError = microerror.New("no resource set")

// IsNoResourceSet asserts noResourceSetError.
func IsNoResourceSet(err error) bool {
	return microerror.Cause(err) == noResourceSetError
}

var statusForbiddenError = microerror.New("status forbidden")

// IsStatusForbiddenError asserts statusForbiddenError.
func IsStatusForbidden(err error) bool {
	return microerror.Cause(err) == statusForbiddenError
}
