// Package updatenecessarycontext stores and accesses the update necessary in
// context.Context.
package updatenecessarycontext

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// updateNecessaryKey is the key for update necessary values in context.Context.
// Clients use updatenecessarycontext.NewContext and
// updatenecessarycontext.FromContext instead of using this key directly.
var updateNecessaryKey key = "updatenecessary"

// NewContext returns a new context.Context that carries value v.
func NewContext(ctx context.Context, v chan struct{}) context.Context {
	if v == nil {
		return ctx
	}

	return context.WithValue(ctx, updateNecessaryKey, v)
}

// FromContext returns the update necessary channel, if any.
func FromContext(ctx context.Context) (chan struct{}, bool) {
	v, ok := ctx.Value(updateNecessaryKey).(chan struct{})
	return v, ok
}

// IsUpdateNecessary checks whether the given context obtains information about
// the update necessary channel as defined in this package, if any update
// necessary channel is present.
//
// NOTE that the update necessary channel, if any found, is only used to be closed
// to signal updates are necessary. It is not guaranteed that the channel is
// buffered or read from. Clients must not write to it. Otherwise the update
// necessary channel will block eventually. It is safe to signal updates are
// necessary via SetUpdateNecessary.
func IsUpdateNecessary(ctx context.Context) bool {
	updateNecessary, updateNecessaryExists := FromContext(ctx)
	if updateNecessaryExists {
		select {
		case <-updateNecessary:
			return true
		default:
			// fall thorugh
		}
	}

	return false
}

// SetUpdateNecessary is a safe way to signal updates are necessary.
func SetUpdateNecessary(ctx context.Context) {
	updateNecessary, updateNecessaryExists := FromContext(ctx)
	if updateNecessaryExists && !IsUpdateNecessary(ctx) {
		close(updateNecessary)
	}
}
