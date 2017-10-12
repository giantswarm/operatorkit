package metricsresource

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/operatorkit/framework"
)

// WrapConfig is the configuration used to wrap resources with metrics
// resources.
type WrapConfig struct {
	// Settings.
	Name string
}

// DefaultWrapConfig provides a default configuration to wrap resource with
// metrics resources. by best effort.
func DefaultWrapConfig() WrapConfig {
	return WrapConfig{
		// Settings.
		Name: "",
	}
}

// Wrap wraps each given resource with a metrics resource and returns the list
// of wrapped resources.
func Wrap(resources []framework.Resource, config WrapConfig) ([]framework.Resource, error) {
	var wrapped []framework.Resource

	for _, r := range resources {
		resourceConfig := DefaultConfig()

		resourceConfig.Resource = r

		resourceConfig.Name = config.Name

		prometheusResource, err := New(resourceConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		wrapped = append(wrapped, prometheusResource)
	}

	return wrapped, nil
}
