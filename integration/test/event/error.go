//go:build k8srequired
// +build k8srequired

package event

import (
	"github.com/giantswarm/microerror"
)

var eventError = &microerror.Error{
	Kind: "EventError",
	Desc: "Error of an event",
}

// IsEventError asserts eventError.
func IsEventError(err error) bool {
	return microerror.Cause(err) == eventError
}
