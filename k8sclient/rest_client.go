package k8sclient

import (
	"github.com/giantswarm/microerror"
	"k8s.io/client-go/rest"
)

func NewRest(config Config) (rest.Interface, error) {
	err := config.Validate()
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if config.Group == "" {
		return nil, microerror.Maskf(invalidConfigError, "Group must not be empty")
	}
	if config.Version == "" {
		return nil, microerror.Maskf(invalidConfigError, "Version must not be empty")
	}

	restConfig, err := config.ToK8sRestConfig()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	client, err := rest.RESTClientFor(restConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	client.Get().Namespace("...").Resource("kvms").DoRaw()

	return client, nil
}
