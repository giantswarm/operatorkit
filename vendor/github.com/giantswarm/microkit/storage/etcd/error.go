package etcd

import (
	"github.com/juju/errgo"
)

var createFailedError = errgo.New("create failed")

// IsCreateFailed asserts createFailedError.
func IsCreateFailed(err error) bool {
	return errgo.Cause(err) == createFailedError
}

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError
}

var multipleValuesError = errgo.New("multiple values")

// IsMultipleValuesFound asserts multipleValuesError.
func IsMultipleValuesFound(err error) bool {
	return errgo.Cause(err) == multipleValuesError
}

var notFoundError = errgo.New("not found")

// IsNotFound asserts notFoundError.
func IsNotFound(err error) bool {
	return errgo.Cause(err) == notFoundError
}
