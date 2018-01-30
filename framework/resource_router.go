package framework

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/micrologger/loggermeta"
	"k8s.io/apimachinery/pkg/api/meta"

	"github.com/giantswarm/operatorkit/framework/context/reconciliationcanceledcontext"
	"github.com/giantswarm/operatorkit/framework/context/resourcecanceledcontext"
	"github.com/giantswarm/operatorkit/framework/context/updateallowedcontext"
	"github.com/giantswarm/operatorkit/framework/context/updatenecessarycontext"
)

type ResourceRouterConfig struct {
	CtxFunc      func(ctx context.Context, obj interface{}) (context.Context, error)
	Handles      func(obj interface{}) bool
	Logger       micrologger.Logger
	ResourceFunc func(ctx context.Context, obj interface{}) ([]Resource, error)
}

func DefaultResourceRouterResourceFunc(rs []Resource) func(ctx context.Context, obj interface{}) ([]Resource, error) {
	return func(ctx context.Context, obj interface{}) ([]Resource, error) {
		return rs, nil
	}
}

type ResourceRouter struct {
	ctxFunc      func(ctx context.Context, obj interface{}) (context.Context, error)
	handles      func(obj interface{}) bool
	logger       micrologger.Logger
	resourceFunc func(ctx context.Context, obj interface{}) ([]Resource, error)
}

func NewResourceRouter(c ResourceRouterConfig) (*ResourceRouter, error) {
	if c.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", c)
	}
	if c.Handles == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Handles must not be empty", c)
	}
	if c.ResourceFunc == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.ResourceFunc must not be empty", c)
	}

	if c.CtxFunc == nil {
		c.CtxFunc = func(ctx context.Context, obj interface{}) (context.Context, error) {
			return ctx, nil
		}
	}

	r := &ResourceRouter{
		ctxFunc:      c.CtxFunc,
		handles:      c.Handles,
		logger:       c.Logger,
		resourceFunc: c.ResourceFunc,
	}

	return r, nil
}

func (r *ResourceRouter) CtxFunc() func(ctx context.Context, obj interface{}) (context.Context, error) {
	return func(ctx context.Context, obj interface{}) (context.Context, error) {
		ctx = reconciliationcanceledcontext.NewContext(ctx, make(chan struct{}))
		ctx = resourcecanceledcontext.NewContext(ctx, make(chan struct{}))
		ctx = updateallowedcontext.NewContext(ctx, make(chan struct{}))
		ctx = updatenecessarycontext.NewContext(ctx, make(chan struct{}))

		ctx, err := r.ctxFunc(ctx, obj)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		accessor, err := meta.Accessor(obj)
		if err != nil {
			r.logger.LogCtx(ctx, "warning", fmt.Sprintf("cannot create accessor for object %#v", obj))
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
}

func (r *ResourceRouter) Handles(obj interface{}) bool {
	return r.handles(obj)
}

func (r *ResourceRouter) ResourceFunc() func(ctx context.Context, obj interface{}) ([]Resource, error) {
	return r.resourceFunc
}
