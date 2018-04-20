// +build k8srequired

package controlflow

import (
	"github.com/giantswarm/microerror"
)

var countMismatchError = microerror.New("count mismatch")

// IsCountMismatch asserts countMismatchError.
func IsCountMismatch(err error) bool {
	return microerror.Cause(err) == countMismatchError
}
