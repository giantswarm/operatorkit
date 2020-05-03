// Package cachekeycontext stores and accesses the local context key in
// context.Context.
package updateallowedcontext

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// ctxKey is the key for update allowed values in context.Context. Clients use
// updateallowedcontext.NewContext and updateallowedcontext.FromContext instead
// of using this key directly.
var ctxKey key = "updateallowed"

// NewContext returns a new context.Context that can be used to check if updates
// are allowed.
func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, make(chan struct{}))
}

// FromContext returns the update allowed channel, if any.
func FromContext(ctx context.Context) (chan struct{}, bool) {
	v, ok := ctx.Value(ctxKey).(chan struct{})
	return v, ok
}

// IsUpdateAllowed checks whether the given context obtains information about
// the update allowed channel as defined in this package, if any update allowed
// channel is present.
//
// NOTE that the update allowed channel, if any found, is only used to be closed
// to signal updates are allowed. It is not guaranteed that the channel is
// buffered or read from. Clients must not write to it. Otherwise the update
// allowed channel will block eventually. It is safe to signal updates are
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

// SetUpdateAllowed is a safe way to signal updates are allowed.
func SetUpdateAllowed(ctx context.Context) {
	updateAllowed, updateAllowedExists := FromContext(ctx)
	if updateAllowedExists && !IsUpdateAllowed(ctx) {
		close(updateAllowed)
	}
}
