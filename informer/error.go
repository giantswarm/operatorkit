package informer

import (
	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
)

var alreadyRegisteredError = microerror.New("already registered")

// IsAlreadyRegisteredError asserts alreadyRegisteredError.
func IsAlreadyRegisteredError(err error) bool {
	c := microerror.Cause(err)
	_, ok := c.(prometheus.AlreadyRegisteredError)
	if ok {
		return true
	}
	if c == alreadyRegisteredError {
		return true
	}

	return false
}

var contextCanceledError = microerror.New("context canceled")

// IsContextCanceled asserts contextCanceledError.
func IsContextCanceled(err error) bool {
	return microerror.Cause(err) == contextCanceledError
}

var initializationTimedOutError = microerror.New("initialization timed out")

// IsInitializationTimedOut asserts initializationTimedOutError.
func IsInitializationTimedOut(err error) bool {
	return microerror.Cause(err) == initializationTimedOutError
}

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidEventError = microerror.New("invalid event")

// IsInvalidEvent asserts invalidEventError.
func IsInvalidEvent(err error) bool {
	return microerror.Cause(err) == invalidEventError
}
