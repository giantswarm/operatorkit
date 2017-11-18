package k8sclient

import (
	"github.com/giantswarm/microerror"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

func NewExt(config Config) (apiextensionsclient.Interface, error) {
	err := config.Validate()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	restConfig, err := config.ToK8sRestConfig()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	client, err := apiextensionsclient.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return client, nil
}
