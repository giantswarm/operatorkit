package metricsresource

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/resource"
)

type WrapConfig struct {
}

// Wrap wraps each given resource with a metrics resource and returns the list
// of wrapped resources.
func Wrap(resources []resource.Interface, config WrapConfig) ([]resource.Interface, error) {
	var wrapped []resource.Interface

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
