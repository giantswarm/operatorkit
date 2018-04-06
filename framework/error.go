package framework

import (
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

var invalidResourceVersionError = errors.StatusError{
	ErrStatus: metav1.Status{
		Status:  "Failure",
		Code:    500,
		Message: "Testing value /metadata/resourceVersion failed",
	},
}

// IsInvalidResourceVersionError asserts invalidResourceVersionError.
func IsInvalidResourceVersionError(err error) bool {
	return err == error(&invalidResourceVersionError)
}

var noResourceSetError = microerror.New("no resource set")

// IsNoResourceSet asserts noResourceSetError.
func IsNoResourceSet(err error) bool {
	return microerror.Cause(err) == noResourceSetError
}
