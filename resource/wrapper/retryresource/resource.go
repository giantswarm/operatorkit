package retryresource

import (
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/resource"
)

type Config struct {
	BackOff  backoff.Interface
	Logger   micrologger.Logger
	Resource resource.Interface
}

func New(config Config) (resource.Interface, error) {
	if config.BackOff == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.BackOff must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
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
