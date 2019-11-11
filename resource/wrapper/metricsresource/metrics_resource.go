package metricsresource

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger/microloggertest"

	"github.com/giantswarm/operatorkit/resource"
	"github.com/giantswarm/operatorkit/resource/crud"
	"github.com/giantswarm/operatorkit/resource/wrapper/internal"
)

type Config struct {
	Resource resource.Interface
}

// New returns a new metrics resource according to the configured resource's
// implementation, which might be resource.Interface or crud.Interface. This has
// then different implications on how to measure metrics for the different
// methods of the interfaces.
func New(config Config) (resource.Interface, error) {
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
	}

	var err error

	// If crud.Interface can be extracted from this resource wrap it.
	// In this case GetCurrentState, GetDesiredState, NewUpdatePatch,
	// NewDeletePatch, ApplyCreateChange, ApplyDeleteChange and
	// ApplyUpdateChange are wrapped with retries.
	crudInterface, ok := internal.CRUD(config.Resource)
	if ok {
		var wrappedCRUD *crudResource
		{
			c := crudResourceConfig{
				CRUD: crudInterface,
			}

			wrappedCRUD, err = newCRUDResource(c)
			if err != nil {
				return nil, microerror.Mask(err)
			}
		}

		{
			c := crud.ResourceConfig{
				CRUD: wrappedCRUD,
				// TODO the crud resource requires a logger but in fact does not use it
				// at all. To not break the metricsresource interface we just configure
				// some dummy here.
				Logger: microloggertest.New(),
			}

			r, err := crud.NewResource(c)
			if err != nil {
				return nil, microerror.Mask(err)
			}

			return r, nil
		}
	}

	// If crud.Interface can't be extracted resource wrap only resource.Interface
	// EnsureCreated and EnsureDeleted methods with retries.
	{
		c := basicResourceConfig{
			Resource: config.Resource,
		}

		r, err := newBasicResource(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		return r, nil
	}
}
