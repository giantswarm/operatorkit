// Package deletionallowedcontext stores and accesses the deletion allowed in
// context.Context.
package deletionallowedcontext

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// deletionAllowedKey is the key for deletion allowed values in context.Context.
// Clients use deletionallowedcontext.NewContext and
// deletionallowedcontext.FromContext instead of using this key directly.
var deletionAllowedKey key = "deletionallowed"

// NewContext returns a new context.Context that carries value v.
func NewContext(ctx context.Context, v chan struct{}) context.Context {
	if v == nil {
		return ctx
	}

	return context.WithValue(ctx, deletionAllowedKey, v)
}

// FromContext returns the deletion allowed channel, if any.
func FromContext(ctx context.Context) (chan struct{}, bool) {
	v, ok := ctx.Value(deletionAllowedKey).(chan struct{})
	return v, ok
}

// IsDeletionAllowed checks whether the given context obtains information about
// the deletion allowed channel as defined in this package, if any deletion
// allowed channel is present.
//
// NOTE that the deletion allowed channel, if any found, is only used to be
// closed to signal deletions are allowed. It is not guaranteed that the channel
// is buffered or read from. Clients must not write to it. Otherwise the
// deletion allowed channel will block eventually. It is safe to signal
// deletions are allowed via SetDeletionAllowed.
func IsDeletionAllowed(ctx context.Context) bool {
	deletionAllowed, deletionAllowedExists := FromContext(ctx)
	if deletionAllowedExists {
		select {
		case <-deletionAllowed:
			return true
		default:
			// fall thorugh
		}
	}

	return false
}

// SetDeletionAllowed is a safe way to signal deletions are allowed.
func SetDeletionAllowed(ctx context.Context) {
	deletionAllowed, deletionAllowedExists := FromContext(ctx)
	if deletionAllowedExists && !IsDeletionAllowed(ctx) {
		close(deletionAllowed)
	}
}
