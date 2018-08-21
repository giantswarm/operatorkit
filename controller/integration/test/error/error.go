// +build k8srequired

package error

import (
	"github.com/giantswarm/microerror"
)

var countMismatchError = &microerror.Error{
	Kind: "countMismatchError",
}

// IsCountMismatch asserts countMismatchError.
func IsCountMismatch(err error) bool {
	return microerror.Cause(err) == countMismatchError
}

var testError = &microerror.Error{
	Kind: "testError",
}

// IsTestError asserts testError.
func IsTestError(err error) bool {
	return microerror.Cause(err) == testError
}
