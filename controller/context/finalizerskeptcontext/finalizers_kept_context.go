// Package finalizerskeptcontext stores and accesses the kept in
// context.Context.
package finalizerskeptcontext

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// keptKey is the key for kept values in context.Context. Clients use
// finalizerskeptcontext.NewContext and finalizerskeptcontext.FromContext
// instead of using this key directly.
var keptKey key = "kept"

// NewContext returns a new context.Context that carries value v.
func NewContext(ctx context.Context, v chan struct{}) context.Context {
	if v == nil {
		return ctx
	}

	return context.WithValue(ctx, keptKey, v)
}

// FromContext returns the kept channel, if any.
func FromContext(ctx context.Context) (chan struct{}, bool) {
	v, ok := ctx.Value(keptKey).(chan struct{})
	return v, ok
}

// IsKept checks whether the given context obtains information about the kept
// channel as defined in this package, if any kept channel is present.
//
// NOTE that the kept channel, if any found, is only used to be closed to signal
// cancelation. It is not guaranteed that the channel is buffered or read from.
// Clients must not write to it. Otherwise the kept channel will block
// eventually. It is safe to signal cancelation via SetKept.
func IsKept(ctx context.Context) bool {
	kept, keptExists := FromContext(ctx)
	if keptExists {
		select {
		case <-kept:
			return true
		default:
			// fall thorugh
		}
	}

	return false
}

// SetKept is a safe way to signal cancelation.
func SetKept(ctx context.Context) {
	kept, keptExists := FromContext(ctx)
	if keptExists && !IsKept(ctx) {
		close(kept)
	}
}
