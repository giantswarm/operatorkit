// Package updateallowedcontext stores and accesses the update allowed in
// context.Context.
package updateallowedcontext

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// updateAllowedKey is the key for update allowed values in context.Context.
// Clients use updateallowedcontext.NewContext and
// updateallowedcontext.FromContext instead of using this key directly.
var updateAllowedKey key = "updateallowed"

// NewContext returns a new context.Context that carries value v.
func NewContext(ctx context.Context, v chan struct{}) context.Context {
	if v == nil {
		return ctx
	}

	return context.WithValue(ctx, updateAllowedKey, v)
}

// FromContext returns the update allowed, if any.
func FromContext(ctx context.Context) (chan struct{}, bool) {
	v, ok := ctx.Value(updateAllowedKey).(chan struct{})
	return v, ok
}

// IsUpdateAllowed checks whether the given context obtains information about
// the update allowed as defined in this package, if any update allowed is
// present.
//
// NOTE that the update allowed channel, if any found, is only used to be closed
// to signal updates are allowed. It is not guaranteed that the channel is
// buffered or read from. Clients must not write to it. Otherwise the update
// allowed channel will block eventually. It is save to signal updates are
// allowed via SetUpdateAllowed.
func IsUpdateAllowed(ctx context.Context) bool {
	updateAllowed, updateAllowedExists := FromContext(ctx)
	if updateAllowedExists {
		select {
		case <-updateAllowed:
			return true
		default:
			// fall thorugh
		}
	}

	return false
}

// SetUpdateAllowed is a save way to signal updates are allowed.
func SetUpdateAllowed(ctx context.Context) {
	updateAllowed, updateAllowedExists := FromContext(ctx)
	if updateAllowedExists && !IsUpdateAllowed(ctx) {
		close(updateAllowed)
	}
}
