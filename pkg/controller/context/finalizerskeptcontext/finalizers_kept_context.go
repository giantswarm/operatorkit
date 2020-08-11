// Package finalizerskeptcontext stores and accesses the kept in
// context.Context.
package finalizerskeptcontext

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// keptKey is the key for kept values in context.Context. Clients use
// finalizerskeptcontext.NewContext and finalizerskeptcontext.FromContext
// instead of using this key directly.
var keptKey key = "kept"

// NewContext returns a new context.Context that carries value v.
func NewContext(ctx context.Context, v chan struct{}) context.Context {
	if v == nil {
		return ctx
	}

	return context.WithValue(ctx, keptKey, v)
}

// FromContext returns the kept channel, if any.
func FromContext(ctx context.Context) (chan struct{}, bool) {
	v, ok := ctx.Value(keptKey).(chan struct{})
	return v, ok
}

// IsKept checks whether the given context obtains information about the kept
// channel as defined in this package, if any kept channel is present.
//
// NOTE that the kept channel, if any found, is only used to be closed to signal
// keeping finalizers. It is not guaranteed that the channel is buffered or read
// from. Clients must not write to it. Otherwise the kept channel will block
// eventually. It is safe to signal keeping finalizers via SetKept.
func IsKept(ctx context.Context) bool {
	kept, keptExists := FromContext(ctx)
	if keptExists {
		select {
		case <-kept:
			return true
		default:
			// fall thorugh
		}
	}

	return false
}

// SetKept is a safe way to signal keeping finalizers. When operators manage the
// deletion of resources they use finalizers out of the box when using the
// operatorkit controller. On deletion operators might want to replay the
// deletion process. This is achieved by not removing finalizers from the
// observed runtime object. So instead of returning an error and abusing errors
// for control flow, SetKept can be used to signal the operatorkit controller to
// not remove finalizers at the end of the current reconciliation loop. Most
// likely resource cancelation is desired to be used in combination with
// SetKept. Context based functionality can be composed like the following.
//
//     finalizerskeptcontext.SetKept(ctx)
//     resourcecanceledcontext.SetCanceled(ctx)
//
func SetKept(ctx context.Context) {
	kept, keptExists := FromContext(ctx)
	if keptExists && !IsKept(ctx) {
		close(kept)
	}
}
