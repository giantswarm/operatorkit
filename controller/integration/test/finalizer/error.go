// +build k8srequired

package parallel

import (
	"github.com/giantswarm/microerror"
)

var testError = &microerror.Error{
	Kind: "testError",
}

// IsTest asserts testError.
func IsTest(err error) bool {
	return microerror.Cause(err) == testError
}

var waitError = &microerror.Error{
	Kind: "waitError",
}

// IsWait asserts waitError.
func IsWait(err error) bool {
	return microerror.Cause(err) == waitError
}
