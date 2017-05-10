package cert

import (
	"github.com/juju/errgo"
)

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var secretsRetrievalFailedError = errgo.New("secrets retreival failed")

// IsSecretsRetrievalFailed asserts secretsRetrievalFailedError.
func IsSecretsRetrievalFailed(err error) bool {
	return errgo.Cause(err) == secretsRetrievalFailedError
}
