package crdclient

import (
	"github.com/giantswarm/microerror"
)

var nameConflictError = microerror.New("name conflict")

// IsNameConflict asserts nameConflictError.
func IsNameConflict(err error) bool {
	return microerror.Cause(err) == nameConflictError
}

var notEstablishedError = microerror.New("not established")

// IsNotEstablished asserts notEstablishedError.
func IsNotEstablished(err error) bool {
	return microerror.Cause(err) == notEstablishedError
}

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}
