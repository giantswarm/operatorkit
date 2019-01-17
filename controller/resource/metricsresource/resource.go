package metricsresource

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/operatorkit/controller"
)

type Config struct {
	Resource controller.Resource
}

func New(config Config) (controller.Resource, error) {
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
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
