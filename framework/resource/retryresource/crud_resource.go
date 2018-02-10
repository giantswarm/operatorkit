package retryresource

import (
	"context"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
)

type crudResourceConfig struct {
	Logger   micrologger.Logger
	Resource *framework.CRUDResource

	BackOff backoff.BackOff
}

type crudResource struct {
	resource *framework.CRUDResource
}

func newCRUDResource(config crudResourceConfig) (*crudResource, error) {
	var err error

	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Resource must not be empty")
	}
	if config.BackOff == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.BackOff must not be empty")
	}

	var retryOps *crudResourceOps
	{
		c := crudResourceOpsConfig{
			Logger: config.Logger,
			Ops:    config.Resource.CRUDResourceOps,

			BackOff: config.BackOff,
		}

		retryOps, err = newCRUDResourceOps(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	// Replace ops with wrapped ones.
	config.Resource.CRUDResourceOps = retryOps

	r := &crudResource{
		resource: config.Resource,
	}

	return r, nil
}

func (r *crudResource) EnsureCreated(ctx context.Context, obj interface{}) error {
	err := r.resource.EnsureCreated(ctx, obj)
	return microerror.Mask(err)
}

func (r *crudResource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	err := r.resource.EnsureDeleted(ctx, obj)
	return microerror.Mask(err)
}

func (r *crudResource) Name() string {
	return r.resource.Name()
}

func (r *crudResource) Underlying() framework.Resource {
	return r.resource.Underlying()
}
