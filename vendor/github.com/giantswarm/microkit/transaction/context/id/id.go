// Package id stores and accesses the transaction ID in and from
// context.Context.
package id

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// idKey is the key for context values in context.Context. Clients use
// id.NewContext and id.FromContext instead of using this key directly.
var idKey key = "id"

// NewContext returns a new context.Context that carries value val.
func NewContext(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, idKey, val)
}

// FromContext returns the context value stored in ctx, if any.
func FromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(idKey).(string)
	return val, ok
}
