// Package cancelercontext stores and accesses the HTTP Authorization precondition
// context.Context.
package cancelercontext

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// cancelKey is the key for cancel values in context.Context. Clients use
// cancelercontext.NewContext and cancelercontext.FromContext instead of using this
// key directly.
var cancelKey key = "canceler"

// NewContext returns a new context.Context that carries value v.
func NewContext(ctx context.Context, v chan struct{}) context.Context {
	return context.WithValue(ctx, cancelKey, v)
}

// FromContext returns the HTTP Authorization preconditionx, if
// any.
func FromContext(ctx context.Context) (chan struct{}, bool) {
	v, ok := ctx.Value(cancelKey).(chan struct{})
	return v, ok
}
