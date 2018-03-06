package metricsresource

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/operatorkit/framework"
)

type WrapConfig struct {
	Name string
}

// Wrap wraps each given resource with a metrics resource and returns the list
// of wrapped resources.
func Wrap(resources []framework.Resource, config WrapConfig) ([]framework.Resource, error) {
	var wrapped []framework.Resource

	for _, r := range resources {
		c := Config{
			Resource: r,

			Name: config.Name,
		}

		prometheusResource, err := New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		wrapped = append(wrapped, prometheusResource)
	}

	return wrapped, nil
}
