package framework

import "context"

// Resource are the building blocks of the operator's reconciliation logic.
// Note there can be multiple Resources reconciling the same object in the
// chain. In that case they are guaranteed to be executed in order one one
// after another.
type Resource interface {
	// Name returns the resource's name used for identification.
	Name() string

	// EnsureCreated is called whtn observed object is created or updated.
	// The object is in state after cration or modification. This method must
	// be idempotent.
	EnsureCreated(ctx context.Context, obj interface{}) error
	// EnsureDeleted is called whtn observed object is deleted. The object
	// is in last observed state before the deletion. This method must be
	// idempotent.
	EnsureDeleted(ctx context.Context, obj interface{}) error
}
