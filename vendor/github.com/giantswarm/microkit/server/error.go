package server

import (
	"github.com/juju/errgo"
)

func errorTrace(err error) string {
	switch e := err.(type) {
	case *errgo.Err:
		return e.GoString()
	}
	return "n/a"
}

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var invalidContextError = errgo.New("invalid context")

// IsInvalidContext asserts invalidContextError.
func IsInvalidContext(err error) bool {
	return errgo.Cause(err) == invalidContextError
}

var invalidTransactionIDError = errgo.New("invalid transaction ID")

// IsInvalidTransactionID asserts invalidTransactionIDError.
func IsInvalidTransactionID(err error) bool {
	return errgo.Cause(err) == invalidTransactionIDError
}
