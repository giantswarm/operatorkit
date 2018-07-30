// +build k8srequired

package parallel

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
