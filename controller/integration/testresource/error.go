package testresource

import "github.com/giantswarm/microerror"

var testError = &microerror.Error{
	Kind: "testError",
}

// IsTestError asserts testError.
func IsTestError(err error) bool {
	return microerror.Cause(err) == testError
}
