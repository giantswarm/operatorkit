package testresource

import "github.com/giantswarm/microerror"

var testError = microerror.New("just testing")

// IsTestError asserts testError.
func IsTestError(err error) bool {
	return microerror.Cause(err) == testError
}
