package retryresource

import (
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/handler"
)

type Config struct {
	BackOff  backoff.Interface
	Logger   micrologger.Logger
	Resource handler.Interface
}

// New returns a new retry resource according to the configured resource's
// implementation, which might be handler.Interface or crud.Interface. This has
// then different implications on how to retry the different methods of the
// interfaces.
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
