// +build k8srequired

package controlflow

import (
	"github.com/giantswarm/microerror"
)

var testError = &microerror.Error{
	Kind: "testError",
}

// IsTestError asserts testError.
func IsTestError(err error) bool {
	return microerror.Cause(err) == testError
}

var waitError = &microerror.Error{
	Kind: "waitError",
}

// IsWait asserts waitError.
func IsWait(err error) bool {
	return microerror.Cause(err) == waitError
}
