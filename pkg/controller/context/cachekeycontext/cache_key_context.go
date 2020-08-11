// Package cachekeycontext stores and accesses the local context key in
// context.Context.
package cachekeycontext

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// ctxKey is the key for cache key values in context.Context. Clients use
// cachekeycontext.NewContext and cachekeycontext.FromContext instead of using
// this key directly.
var ctxKey key = "cache-key"

// NewContext returns a new context.Context that can be used to access the
// current reconciliation loop's cache key.
func NewContext(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxKey, v)
}

// FromContext returns the kept channel, if any.
func FromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKey).(string)
	return v, ok
}
