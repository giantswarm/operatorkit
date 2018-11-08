package informer

import (
	"github.com/giantswarm/microerror"
)

var contextCanceledError = &microerror.Error{
	Kind: "contextCanceledError",
}

// IsContextCanceled asserts contextCanceledError.
func IsContextCanceled(err error) bool {
	return microerror.Cause(err) == contextCanceledError
}

var initializationTimedOutError = &microerror.Error{
	Kind: "initializationTimedOutError",
}

// IsInitializationTimedOut asserts initializationTimedOutError.
func IsInitializationTimedOut(err error) bool {
	return microerror.Cause(err) == initializationTimedOutError
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidEventError = &microerror.Error{
	Kind: "invalidEventError",
}

// IsInvalidEvent asserts invalidEventError.
func IsInvalidEvent(err error) bool {
	return microerror.Cause(err) == invalidEventError
}
