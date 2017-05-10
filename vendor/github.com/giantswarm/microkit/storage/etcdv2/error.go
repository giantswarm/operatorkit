package etcdv2

import (
	"github.com/coreos/etcd/client"
	"github.com/juju/errgo"
)

var createFailedError = errgo.New("create failed")

// IsCreateFailed asserts createFailedError.
func IsCreateFailed(err error) bool {
	return errgo.Cause(err) == createFailedError
}

// IsEtcdKeyAlreadyExists is an error matcher for the v2 etcd client.
func IsEtcdKeyAlreadyExists(err error) bool {
	if cErr, ok := err.(client.Error); ok {
		return cErr.Code == client.ErrorCodeNodeExist
	}
	return false
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
