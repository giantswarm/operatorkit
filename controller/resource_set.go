package controller

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/controller/context/finalizerskeptcontext"
	"github.com/giantswarm/operatorkit/controller/context/updateallowedcontext"
	"github.com/giantswarm/operatorkit/resource"
)

type ResourceSetConfig struct {
	// Handles determines if this resource set handles the reconciliation of the
	// object.
	Handles func(obj interface{}) bool
	// InitCtx is to prepare the given context for a single reconciliation loop.
	// Operators can implement common context packages to enable communication
	// between resources. These context packages can be set up within this context
	// initializer function. InitCtx receives the runtime object being reconciled
	// as second argument. Information provided by the runtime object can be used
	// to initialize the context.
	InitCtx func(ctx context.Context, obj interface{}) (context.Context, error)
	// Logger is a usual micrologger instance to emit log messages, if any.
	Logger micrologger.Logger
	// Resources is the list of controller resources being executed on runtime
	// object reconciliation if Handles returns true when asked by the
	// controller. Resources are executed in given order.
	Resources []resource.Interface
}

type ResourceSet struct {
	handles   func(obj interface{}) bool
	initCtx   func(ctx context.Context, obj interface{}) (context.Context, error)
	logger    micrologger.Logger
	resources []resource.Interface
}

func NewResourceSet(c ResourceSetConfig) (*ResourceSet, error) {
	if c.Handles == nil {
		c.Handles = func(obj interface{}) bool { return true }
	}
	if c.InitCtx == nil {
		c.InitCtx = func(ctx context.Context, obj interface{}) (context.Context, error) {
			return ctx, nil
		}
	}
	if c.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", c)
	}
	if len(c.Resources) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resources must not be empty", c)
	}

	r := &ResourceSet{
		handles:   c.Handles,
		initCtx:   c.InitCtx,
		logger:    c.Logger,
		resources: c.Resources,
	}

	return r, nil
}

func (r *ResourceSet) Handles(obj interface{}) bool {
	return r.handles(obj)
}

func (r *ResourceSet) InitCtx(ctx context.Context, obj interface{}) (context.Context, error) {
	ctx = finalizerskeptcontext.NewContext(ctx, make(chan struct{}))
	ctx = updateallowedcontext.NewContext(ctx, make(chan struct{}))

	ctx, err := r.initCtx(ctx, obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return ctx, nil
}

func (r *ResourceSet) Resources() []resource.Interface {
	return r.resources
}
