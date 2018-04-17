package controller

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/micrologger/loggermeta"
	"k8s.io/apimachinery/pkg/api/meta"

	"github.com/giantswarm/operatorkit/controller/context/updateallowedcontext"
	"github.com/giantswarm/operatorkit/controller/context/updatenecessarycontext"
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
	Resources []Resource
}

type ResourceSet struct {
	handles   func(obj interface{}) bool
	initCtx   func(ctx context.Context, obj interface{}) (context.Context, error)
	logger    micrologger.Logger
	resources []Resource
}

func NewResourceSet(c ResourceSetConfig) (*ResourceSet, error) {
	if c.InitCtx == nil {
		c.InitCtx = func(ctx context.Context, obj interface{}) (context.Context, error) {
			return ctx, nil
		}
	}

	if c.Handles == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Handles must not be empty", c)
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

func (r *ResourceSet) InitCtx(ctx context.Context, obj interface{}) (context.Context, error) {
	ctx = updateallowedcontext.NewContext(ctx, make(chan struct{}))
	ctx = updatenecessarycontext.NewContext(ctx, make(chan struct{}))

	ctx, err := r.initCtx(ctx, obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	accessor, err := meta.Accessor(obj)
	if err != nil {
		r.logger.LogCtx(ctx, "function", "InitCtx", "level", "warning", "message", "cannot create accessor for object", "object", fmt.Sprintf("%#v", obj), "stack", fmt.Sprintf("%#v", err))
	} else {
		meta, ok := loggermeta.FromContext(ctx)
		if !ok {
			meta = loggermeta.New()
		}
		meta.KeyVals["object"] = accessor.GetSelfLink()

		ctx = loggermeta.NewContext(ctx, meta)
	}

	return ctx, nil
}

func (r *ResourceSet) Handles(obj interface{}) bool {
	return r.handles(obj)
}

func (r *ResourceSet) Resources() []Resource {
	return r.resources
}
