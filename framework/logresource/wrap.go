package logresource

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/framework"
)

// WrapConfig is the configuration used to wrap resources with log resources.
type WrapConfig struct {
	// Dependencies.
	Logger micrologger.Logger
}

// DefaultWrapConfig provides a default configuration to wrap resource with log
// resources. by best effort.
func DefaultWrapConfig() WrapConfig {
	return WrapConfig{
		// Dependencies.
		Logger: nil,
	}
}

// Wrap wraps each given resource with a log resource and returns the list of
// wrapped resources.
func Wrap(resources []framework.Resource, config WrapConfig) ([]framework.Resource, error) {
	var wrapped []framework.Resource

	for _, r := range resources {
		resourceConfig := DefaultConfig()
		resourceConfig.Logger = config.Logger
		resourceConfig.Resource = r

		logResource, err := New(resourceConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		wrapped = append(wrapped, logResource)
	}

	return wrapped, nil
}
