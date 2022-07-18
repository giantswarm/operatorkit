package wrapper

import (
	"context"

	"github.com/giantswarm/operatorkit/v8/pkg/controller"
)

type Interface interface {
	// CRUD functions for objects.
	CreateObject(ctx context.Context, namespace string, obj interface{}) (interface{}, error)
	DeleteObject(ctx context.Context, name, namespace string) error
	GetObject(ctx context.Context, name, namespace string) (interface{}, error)
	UpdateObject(ctx context.Context, namespace string, obj interface{}) (interface{}, error)

	// Functions for test setup and teardown.
	MustSetup(ctx context.Context, namespace string)
	MustTeardown(ctx context.Context, namespace string)

	// Getters and setters.
	Controller() *controller.Controller
}
