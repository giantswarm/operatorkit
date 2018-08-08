package retryresource

import (
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/controller"
)

type Config struct {
	Logger   micrologger.Logger
	Resource controller.Resource

	BackOff backoff.Interface
}

func New(config Config) (controller.Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Resource must not be empty")
	}

	if config.BackOff == nil {
		config.BackOff = backoff.NewExponential(2*time.Minute, 10*time.Second)
	}

	var err error
	var r controller.Resource

	// CRUD resource special case.
	r, err = newCRUDResourceWrapper(config)
	if isIncompatibleUnderlyingResource(err) {
		// Fall trough. Try wrap with next wrapper.
	} else if err != nil {
		return nil, microerror.Mask(err)
	} else {
		return r, nil
	}

	// Direct resource implmementation.
	r, err = newResourceWrapper(config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}
