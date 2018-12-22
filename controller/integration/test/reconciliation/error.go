// +build k8srequired

package reconciliation

import (
	"github.com/giantswarm/microerror"
)

var waitError = &microerror.Error{
	Kind: "waitError",
}

// IsWait asserts waitError.
func IsWait(err error) bool {
	return microerror.Cause(err) == waitError
}
