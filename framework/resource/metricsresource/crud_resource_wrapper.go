package metricsresource

import (
	"context"

	"github.com/giantswarm/microerror"

	originalframework "github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/resource/internal"
	"github.com/giantswarm/operatorkit/framework/resource/internal/framework"
)

// crudResourceWrapper is a specialized wrapper which wraps
// *framework.CRUDResource.
type crudResourceWrapper struct {
	resource framework.Resource
}

func newCRUDResourceWrapper(config Config) (*crudResourceWrapper, error) {
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Resource must not be empty")
	}

	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Name must not be empty")
	}

	// Wrap underlying resource Ops with retry logic. Underlying resource
	// is a pointer so we can modify it in place.
	{
		underlying, ok := internal.Underlying(config.Resource).(*originalframework.CRUDResource)
		if !ok {
			return nil, microerror.Maskf(incompatibleUnderlyingResourceError, "expected %T", underlying)
		}

		c := crudResourceOpsWrapperConfig{
			Ops: underlying.CRUDResourceOps,

			Name: config.Name,
		}

		wrappedOps, err := newCRUDResourceWrapperOps(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		underlying.CRUDResourceOps = wrappedOps
	}

	r := &crudResourceWrapper{
		resource: config.Resource,
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
func (r *crudResourceWrapper) Wrapped() framework.Resource {
	return r.resource
}
