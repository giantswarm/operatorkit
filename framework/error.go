package framework

import (
	"github.com/giantswarm/microerror"
)

var customObjectVersionNotFoundError = microerror.New("custom object version not found")

// IsCustomObjectVersionNotFound asserts customObjectVersionNotFoundError.
func IsCustomObjectVersionNotFound(err error) bool {
	return microerror.Cause(err) == customObjectVersionNotFoundError
}

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
