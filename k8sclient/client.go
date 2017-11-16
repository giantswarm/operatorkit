package k8sclient

import (
	"github.com/giantswarm/microerror"
	"k8s.io/client-go/kubernetes"
)

func NewClient(config Config) (kubernetes.Interface, error) {
	err := config.Validate()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	restConfig, err := config.ToK8sRestConfig()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return client, nil
}
