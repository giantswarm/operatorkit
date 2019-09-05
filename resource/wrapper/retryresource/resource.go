package retryresource

import (
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/resource"
)

type Config struct {
	Logger   micrologger.Logger
	Resource resource.Interface

	BackOff backoff.Interface
}

func New(config Config) (resource.Interface, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
	}

	if config.BackOff == nil {
		config.BackOff = backoff.NewExponential(2*time.Minute, 10*time.Second)
	}

	var err error
	var r resource.Interface

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
