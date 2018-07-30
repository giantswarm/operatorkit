package controller

import (
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/errors"
)

var executionFailedError = &microerror.Error{
	Kind: "executionFailedError",
}

// IsExecutionFailed asserts executionFailedError.
func IsExecutionFailed(err error) bool {
	return microerror.Cause(err) == executionFailedError
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidRESTClientError = &microerror.Error{
	Kind: "invalidRESTClientError",
}

// IsInvalidRESTClient asserts invalidRESTClientError.
func IsInvalidRESTClient(err error) bool {
	return microerror.Cause(err) == invalidRESTClientError
}

var noResourceSetError = &microerror.Error{
	Kind: "noResourceSetError",
}

// IsNoResourceSet asserts noResourceSetError.
func IsNoResourceSet(err error) bool {
	return microerror.Cause(err) == noResourceSetError
}

var statusForbiddenError = &microerror.Error{
	Kind: "statusForbiddenError",
}

// IsStatusForbiddenError asserts statusForbiddenError and apimachinery
// StatusError with StatusReasonForbidden.
func IsStatusForbidden(err error) bool {
	if err == nil {
		return false
	}

	c := microerror.Cause(err)

	if c == statusForbiddenError {
		return true
	}

	if errors.IsForbidden(c) {
		return true
	}

	return false
}
