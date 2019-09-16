package retryresource

import (
	"context"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/resource"
	"github.com/giantswarm/operatorkit/resource/wrapper/internal"
)

// crudResourceWrapper is a specialized wrapper which wraps
// *controller.CRUDResource.
type crudResourceWrapper struct {
	logger   micrologger.Logger
	resource resource.Interface

	backOff backoff.Interface
}

func newCRUDResourceWrapper(config Config) (*crudResourceWrapper, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
	}

	if config.BackOff == nil {
		config.BackOff = backoff.NewExponential(2*time.Minute, 10*time.Second)
	}

	// Wrap underlying resource Ops with retry logic. Underlying resource
	// is a pointer so we can modify it in place.
	{
		underlying, err := internal.Underlying(config.Resource)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		underlyingCRUD, ok := underlying.(*controller.CRUDResource)
		if !ok {
			return nil, microerror.Maskf(incompatibleUnderlyingResourceError, "expected %T", underlyingCRUD)
		}

		c := crudResourceOpsWrapperConfig{
			Logger: config.Logger,
			Ops:    underlyingCRUD.CRUDResourceOps,

			BackOff: config.BackOff,
		}

		wrappedOps, err := newCRUDResourceWrapperOps(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		underlyingCRUD.CRUDResourceOps = wrappedOps
	}

	r := &crudResourceWrapper{
		logger: config.Logger.With(
			"underlyingResource", config.Resource.Name(),
		),
		resource: config.Resource,

		backOff: config.BackOff,
	}

	return r, nil
}

func (r *crudResourceWrapper) EnsureCreated(ctx context.Context, obj interface{}) error {
	// Just call wrapped resource. Wrapped crudResourceOpsWrapper will do
	// the retries.
	err := r.resource.EnsureCreated(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *crudResourceWrapper) EnsureDeleted(ctx context.Context, obj interface{}) error {
	// Just call wrapped resource. Wrapped crudResourceOpsWrapper will do
	// the retries.
	err := r.resource.EnsureDeleted(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *crudResourceWrapper) Name() string {
	return r.resource.Name()
}

// Wrapped implements internal.Wrapper interface.
func (r *crudResourceWrapper) Wrapped() resource.Interface {
	return r.resource
}
