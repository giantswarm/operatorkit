package k8scrdclient

import (
	"github.com/giantswarm/microerror"
)

var nameConflictError = &microerror.Error{
	Kind: "nameConflictError",
}

// IsNameConflict asserts nameConflictError.
func IsNameConflict(err error) bool {
	return microerror.Cause(err) == nameConflictError
}

var notEstablishedError = &microerror.Error{
	Kind: "notEstablishedError",
}

// IsNotEstablished asserts notEstablishedError.
func IsNotEstablished(err error) bool {
	return microerror.Cause(err) == notEstablishedError
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}
