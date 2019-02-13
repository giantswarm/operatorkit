package metricsresource

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/operatorkit/controller"
)

type WrapConfig struct {
}

// Wrap wraps each given resource with a metrics resource and returns the list
// of wrapped resources.
func Wrap(resources []controller.Resource, config WrapConfig) ([]controller.Resource, error) {
	var wrapped []controller.Resource

	for _, r := range resources {
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
