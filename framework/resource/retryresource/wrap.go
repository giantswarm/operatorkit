package retryresource

import (
	"github.com/cenkalti/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/framework"
)

// WrapConfig is the configuration used to wrap resources with retry resources.
type WrapConfig struct {
	Logger micrologger.Logger

	BackOffFactory func() backoff.BackOff
}

// Wrap wraps each given resource with a retry resource and returns the list of
// wrapped resources.
func Wrap(resources []framework.Resource, config WrapConfig) ([]framework.Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	if config.BackOffFactory == nil {
		config.BackOffFactory = func() backoff.BackOff { return &backoff.ZeroBackOff{} }
	}

	var wrapped []framework.Resource

	for _, r := range resources {
		c := Config{
			BackOff:  config.BackOffFactory(),
			Logger:   config.Logger,
			Resource: r,
		}

		retryResource, err := New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		wrapped = append(wrapped, retryResource)
	}

	return wrapped, nil
}
