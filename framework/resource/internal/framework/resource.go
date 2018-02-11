package framework

import (
	"context"
)

// TODO Remove this file and change imports when original framework.Resource is updated.

// NOTE this is temporary code for transition period. It will be removed as
// soon as new github.com/operatorkit/framework.Resource is created. It should
// look like below.

type Resource interface {
	Name() string

	EnsureCreated(ctx context.Context, obj interface{}) error
	EnsureDeleted(ctx context.Context, obj interface{}) error
}
