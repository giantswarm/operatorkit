package controller

import "context"

// Resource is an interface. Resources are the building blocks of the
// operator's reconciliation logic. Note there can be multiple Resources
// reconciling the same object in the chain. In that case they are guaranteed
// to be executed in order one after another.
type Resource interface {
	// EnsureCreated is called when observed object is created or updated.
	// The object is in state after cration or modification. This method must
	// be idempotent.
	EnsureCreated(ctx context.Context, obj interface{}) error
	// EnsureDeleted is called when observed object is deleted. The object
	// is in last observed state before the deletion. This method must be
	// idempotent.
	EnsureDeleted(ctx context.Context, obj interface{}) error
	// Name returns the resource's name used for identification.
	Name() string
}
