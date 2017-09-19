// Package cancelercontext stores and accesses the canceler in context.Context.
package cancelercontext

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// cancelerKey is the key for canceler values in context.Context. Clients use
// cancelercontext.NewContext and cancelercontext.FromContext instead of using this
// key directly.
var cancelerKey key = "canceler"

// NewContext returns a new context.Context that carries value v.
func NewContext(ctx context.Context, v chan struct{}) context.Context {
	if v == nil {
		return ctx
	}

	return context.WithValue(ctx, cancelerKey, v)
}

// FromContext returns the canceler, if any.
func FromContext(ctx context.Context) (chan struct{}, bool) {
	v, ok := ctx.Value(cancelerKey).(chan struct{})
	return v, ok
}

// IsCanceled checks whether the given context obtains information about the
// canceler as defined in this package, if any canceler is present.
//
// NOTE that the canceler, if any found, is only used to be closed to signal
// cancelation. It is not guaranteed that the channel is buffered or read from.
// Clients must not write to it. Otherwise the canceler will block eventually.
// It is save to signal cancelation via SetCanceled.
func IsCanceled(ctx context.Context) bool {
	canceler, cancelerExists := FromContext(ctx)
	if cancelerExists {
		select {
		case <-canceler:
			return true
		default:
			// fall thorugh
		}
	}

	return false
}

// SetCanceled is a save way to signal cancelation.
func SetCanceled(ctx context.Context) {
	canceler, cancelerExists := FromContext(ctx)
	if cancelerExists && !IsCanceled(ctx) {
		close(canceler)
	}
}
