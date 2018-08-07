package retryresource

import (
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/controller"
)

// WrapConfig is the configuration used to wrap resources with retry resources.
type WrapConfig struct {
	Logger micrologger.Logger

	BackOffFactory func() backoff.Interface
}

// Wrap wraps each given resource with a retry resource and returns the list of
// wrapped resources.
func Wrap(resources []controller.Resource, config WrapConfig) ([]controller.Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	if config.BackOffFactory == nil {
		config.BackOffFactory = func() backoff.Interface { return backoff.NewMaxRetries(3, 1*time.Second) }
	}

	var wrapped []controller.Resource

	for _, r := range resources {
		c := Config{
			Logger:   config.Logger,
			Resource: r,

			BackOff: config.BackOffFactory(),
		}

		retryResource, err := New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		wrapped = append(wrapped, retryResource)
	}

	return wrapped, nil
}
