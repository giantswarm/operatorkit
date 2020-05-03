// Package cachekeycontext stores and accesses the local context key in
// context.Context.
package cachekeycontext

import (
	"context"
	"strconv"
	"time"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// ctxKey is the key for kept values in context.Context. Clients use
// cachekeycontext.NewContext and cachekeycontext.FromContext instead of using
// this key directly.
var ctxKey key = "kept"

// NewContext returns a new context.Context that carries value v.
func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, strconv.Itoa(int(time.Now().UnixNano())))
}

// FromContext returns the kept channel, if any.
func FromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKey).(string)
	return v, ok
}
