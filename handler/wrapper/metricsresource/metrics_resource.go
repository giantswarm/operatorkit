package metricsresource

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/operatorkit/handler"
)

type Config struct {
	Resource handler.Interface
}

// New returns a new metrics resource according to the configured resource's
// implementation, which might be handler.Interface or crud.Interface. This has
// then different implications on how to measure metrics for the different
// methods of the interfaces.
func New(config Config) (handler.Interface, error) {
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
	}

	var err error

	// If crud.Interface can't be extracted resource wrap only handler.Interface
	// EnsureCreated and EnsureDeleted methods with retries.
	{
		c := basicResourceConfig(config)

		r, err := newBasicResource(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		return r, nil
	}
}
