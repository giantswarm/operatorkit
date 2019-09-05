package testresourceset

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/controller/integration/testresource"
	"github.com/giantswarm/operatorkit/resource"
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
	Resources []resource.Interface

	ProjectName string
}

func New(config Config) (*controller.ResourceSet, error) {
	var err error
	var resources []resource.Interface

	if len(config.Resources) == 0 {
		var tr resource.Interface
		{
			c := testresource.Config{}

			tr, err = testresource.New(c)
			if err != nil {
				return nil, microerror.Mask(err)
			}
		}

		resources = []resource.Interface{
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
