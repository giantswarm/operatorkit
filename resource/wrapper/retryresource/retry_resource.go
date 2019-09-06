package retryresource

import (
	"context"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/resource"
	"github.com/giantswarm/operatorkit/resource/crud"
	"github.com/giantswarm/operatorkit/resource/wrapper/internal"
)

type Config struct {
	BackOff  backoff.Interface
	Logger   micrologger.Logger
	Resource resource.Interface
}

type Resource struct {
	resource resource.Interface
}

// New returns a new retry resource according to the configured resource's
// implementation, which might be resource.Interface or crud.Interface. This has
// then different implications on how to retry the different methods of the
// interfaces.
func New(config Config) (*Resource, error) {
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
	}

	var err error
	var wrapped resource.Interface

	var u resource.Interface
	{
		u, err = internal.Underlying(config.Resource)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	// Here we check if the configured resource is actually a CRUD Resource
	// implementation and wrap it accordingly. In this case we have to wrap
	// GetCurrentState, GetDesiredState, NewUpdatePatch, NewDeletePatch,
	// ApplyCreateChange, ApplyDeleteChange and ApplyUpdateChange to execute the
	// retry logic properly.
	ci, ok := u.(crud.Interface)
	if ok {
		var crudResource crud.Interface
		{
			c := crudResourceConfig{
				BackOff: config.BackOff,
				CRUD:    ci,
				Logger:  config.Logger,
			}

			crudResource, err = newCRUDResource(c)
			if err != nil {
				return nil, microerror.Mask(err)
			}
		}

		{
			c := crud.ResourceConfig{
				CRUD:   crudResource,
				Logger: config.Logger,
			}

			wrapped, err = crud.NewResource(c)
			if err != nil {
				return nil, microerror.Mask(err)
			}
		}
	}

	// Here we check if the configured resource is actually a simple Resource
	// implementation and wrap it accordingly. In this case we have to wrap
	// EnsureCreated and EnsureDeleted. to execute the retry logic properly.
	ri, ok := u.(resource.Interface)
	if ok {
		c := simpleResourceConfig{
			BackOff:  config.BackOff,
			Logger:   config.Logger,
			Resource: ri,
		}

		wrapped, err = newSimpleResource(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	r := &Resource{
		resource: wrapped,
	}

	return r, nil
}

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	err := r.resource.EnsureCreated(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	err := r.resource.EnsureDeleted(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) Name() string {
	return r.resource.Name()
}
