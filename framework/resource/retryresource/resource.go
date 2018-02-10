package retryresource

import (
	"github.com/cenkalti/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
)

type Config struct {
	Logger   micrologger.Logger
	Resource framework.Resource

	BackOff backoff.BackOff
}

type Resource struct {
	logger   micrologger.Logger
	resource framework.Resource

	backOff backoff.BackOff
}

func New(config Config) (framework.Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Resource must not be empty")
	}

	if config.BackOff == nil {
		config.BackOff = backoff.NewExponentialBackOff()
	}

	// If the resource is instance of *framework.CRUDResource we wrap it
	// with a crudResource otherwise we wrap with basicResource.
	crudResource, ok := config.Resource.Underlying().(*framework.CRUDResource)
	if ok {
		c := crudResourceConfig{
			Logger:   config.Logger,
			Resource: crudResource,

			BackOff: config.BackOff,
		}

		r, err := newCRUDResource(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		return r, nil
	} else {
		// TODO uncomment when new Resource interface is created.
		//c := basicResourceConfig{
		//	Logger:   config.Logger,
		//	Resource: config.Resource,
		//
		//	BackOff: config.BackOff,
		//}
		//
		//r, err := newBasicResource(c)
		//if err != nil {
		//	return nil, microerror.Mask(err)
		//}
		//
		//return r, nil
	}

	// TODO remove code below when new Resource interface is created.
	c := crudResourceOpsConfig{
		Logger:   config.Logger,
		Resource: config.Resource,

		BackOff: config.BackOff,
	}

	r, err := newCRUDResourceOps(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}
