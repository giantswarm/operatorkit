// Package tracked stores and accesses the information of if a transaction is
// already tracked in and from context.Context.
package tracked

import (
	"golang.org/x/net/context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// trackedKey is the key for context values in context.Context. Clients use
// tracked.NewContext and tracked.FromContext instead of using this key directly.
var trackedKey key = "tracked"

// NewContext returns a new context.Context that carries value val.
func NewContext(ctx context.Context, val bool) context.Context {
	return context.WithValue(ctx, trackedKey, val)
}

// FromContext returns the context value stored in ctx, if any.
func FromContext(ctx context.Context) (bool, bool) {
	val, ok := ctx.Value(trackedKey).(bool)
	return val, ok
}
