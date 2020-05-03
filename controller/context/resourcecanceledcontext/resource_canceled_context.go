// Package resourcecanceledcontext stores and accesses the canceled in
// context.Context.
package resourcecanceledcontext

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// ctxKey is the key for canceled values in context.Context. Clients use
// resourcecanceledcontext.NewContext and resourcecanceledcontext.FromContext
// instead of using this key directly.
var ctxKey key = "canceled"

// NewContext returns a new context.Context that carries value v.
func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, make(chan struct{}))
}

// FromContext returns the canceled channel, if any.
func FromContext(ctx context.Context) (chan struct{}, bool) {
	v, ok := ctx.Value(ctxKey).(chan struct{})
	return v, ok
}

// IsCanceled checks whether the given context obtains information about the
// canceled channel as defined in this package, if any canceled channel is
// present.
//
// NOTE that the canceled channel, if any found, is only used to be closed to
// signal cancelation. It is not guaranteed that the channel is buffered or read
// from. Clients must not write to it. Otherwise the canceled channel will block
// eventually. It is safe to signal cancelation via SetCanceled.
func IsCanceled(ctx context.Context) bool {
	canceled, canceledExists := FromContext(ctx)
	if canceledExists {
		select {
		case <-canceled:
			return true
		default:
			// fall thorugh
		}
	}

	return false
}

// SetCanceled is a safe way to signal cancelation.
func SetCanceled(ctx context.Context) {
	canceled, canceledExists := FromContext(ctx)
	if canceledExists && !IsCanceled(ctx) {
		close(canceled)
	}
}
