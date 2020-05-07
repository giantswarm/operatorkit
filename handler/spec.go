package handler

import "context"

type Request struct {
	Obj interface{}
}

type Response struct {
}

// Interface defines the building blocks of an operator's reconciliation logic.
// Note there can be multiple Handlers reconciling the same object in a chain.
// In that case they are guaranteed to be executed in order one after another.
type Interface interface {
	// EnsureCreated is called when the observed runtime object is created or
	// updated. The runtime object observed is contained in req. After the
	// successful execution of EnsureCreated, systems being managed have created
	// or updated system resources. This method must be idempotent.
	EnsureCreated(ctx context.Context, req Request) (*Response, error)
	// EnsureDeleted is called when the observed runtime object is requested to be
	// deleted, which means its DeletionTimestamp is set, but the runtime object
	// itself is not garbage collected yet. The runtime object observed is
	// contained in req. After the execution of EnsureDeleted, systems being
	// managed have deleted system resources. If deletion could not be done
	// successfully handler implementations must request to keep finalizers using
	// the available controller context control flow primitives. In case
	// EnsureDeleted returns an error, finalizers are kept automatically. This
	// method must be idempotent.
	EnsureDeleted(ctx context.Context, req Request) (*Response, error)
}
