package retryresource

import (
	"github.com/cenkalti/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/framework"
)

// WrapConfig is the configuration used to wrap resources with retry resources.
type WrapConfig struct {
	BackOffFactory func() backoff.BackOff
	Logger         micrologger.Logger
}

// DefaultWrapConfig provides a default configuration to wrap resource with
// retry resources. by best effort.
func DefaultWrapConfig() WrapConfig {
	var err error

	var newLogger micrologger.Logger
	{
		config := micrologger.DefaultConfig()
		newLogger, err = micrologger.New(config)
		if err != nil {
			panic(err)
		}
	}

	return WrapConfig{
		// Dependencies.
		BackOffFactory: func() backoff.BackOff { return &backoff.ZeroBackOff{} },
		Logger:         newLogger,
	}
}

// Wrap wraps each given resource with a retry resource and returns the list of
// wrapped resources.
func Wrap(resources []framework.Resource, config WrapConfig) ([]framework.Resource, error) {
	if config.BackOffFactory == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.BackOffFactory must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	var wrapped []framework.Resource

	for _, r := range resources {
		resourceConfig := DefaultConfig()
		resourceConfig.BackOff = config.BackOffFactory()
		resourceConfig.Logger = config.Logger
		resourceConfig.Resource = r

		retryResource, err := New(resourceConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		wrapped = append(wrapped, retryResource)
	}

	return wrapped, nil
}
