package retryresource

import (
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/handler"
)

// WrapConfig is the configuration used to wrap handlers with retry handlers.
type WrapConfig struct {
	BackOffFactory func() backoff.Interface
	Logger         micrologger.Logger
}

// Wrap wraps each given resource with a retry resource and returns the list of
// wrapped handlers.
func Wrap(handlers []handler.Interface, config WrapConfig) ([]handler.Interface, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.BackOffFactory == nil {
		config.BackOffFactory = func() backoff.Interface { return backoff.NewMaxRetries(3, 1*time.Second) }
	}

	var wrapped []handler.Interface

	for _, r := range handlers {
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
