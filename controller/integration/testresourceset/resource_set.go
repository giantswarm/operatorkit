// +build k8srequired

package testresourceset

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/operatorkit/controller"

	"github.com/giantswarm/operatorkit/controller/integration/testresource"
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
	Resources []controller.Resource

	ProjectName string
}

func New(config Config) (*controller.ResourceSet, error) {
	var err error
	var resources []controller.Resource

	if len(config.Resources) == 0 {
		var tr controller.Resource
		{
			c := testresource.Config{}

			tr, err = testresource.New(c)
			if err != nil {
				return nil, microerror.Mask(err)
			}
		}

		resources = []controller.Resource{
			tr,
		}
	} else {
		resources = config.Resources
	}

	handlesFunc := func(obj interface{}) bool {
		return true
	}

	var resourceSet *controller.ResourceSet
	{
		c := controller.ResourceSetConfig{
			Handles:   handlesFunc,
			Logger:    config.Logger,
			Resources: resources,
		}

		resourceSet, err = controller.NewResourceSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return resourceSet, nil
}
