package storage

import (
	"github.com/juju/errgo"

	"github.com/giantswarm/microkit/storage/etcd"
	"github.com/giantswarm/microkit/storage/etcdv2"
	"github.com/giantswarm/microkit/storage/memory"
)

var invalidConfigError = errgo.New("invalid config")

// IsInvalidConfig represents the error matcher for public use. Services using
// the storage service internally should use this public key matcher to verify
// if some storage error is of type "key not found", instead of using a specific
// error matching of some specific storage implementation. This public error
// matcher groups all necessary error matchers of more specific storage
// implementations.
func IsInvalidConfig(err error) bool {
	return errgo.Cause(err) == invalidConfigError || etcd.IsInvalidConfig(err) || etcdv2.IsInvalidConfig(err)
}

// IsNotFound represents the error matcher for public use. Services using the
// storage service internally should use this public key matcher to verify if
// some storage error is of type "key not found", instead of using a specific
// error matching of some specific storage implementation. This public error
// matcher groups all necessary error matchers of more specific storage
// implementations.
func IsNotFound(err error) bool {
	return etcd.IsNotFound(err) || etcdv2.IsNotFound(err) || memory.IsNotFound(err)
}
