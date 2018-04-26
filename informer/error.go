package informer

import (
	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
)

var alreadyRegisteredError = microerror.New("duplicate metrics collector registration attempted")

// IsAlreadyRegisteredError asserts alreadyRegisteredError.
func IsAlreadyRegisteredError(err error) bool {
	c := microerror.Cause(err)
	are, ok := c.(prometheus.AlreadyRegisteredError)
	if !ok {
		return false
	}
	if are == alreadyRegisteredError {
		return true
	}

	return false
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
