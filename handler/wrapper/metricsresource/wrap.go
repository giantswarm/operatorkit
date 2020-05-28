package metricsresource

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/handler"
)

// WrapConfig is the configuration used to wrap handlers with metrics handlers.
type WrapConfig struct {
}

// Wrap wraps each given resource with a metrics resource and returns the list of
// wrapped handlers.
func Wrap(handlers []handler.Interface, config WrapConfig) ([]handler.Interface, error) {
	var wrapped []handler.Interface

	for _, r := range handlers {
		c := Config{
			Resource: r,
		}

		metricsResource, err := New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		wrapped = append(wrapped, metricsResource)
	}

	return wrapped, nil
}
