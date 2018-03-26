// +build k8srequired

package testresourceset

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/operatorkit/framework/integration/testresource"
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	ProjectName string
}

func New(config Config) (*framework.ResourceSet, error) {
	var err error
	var tr framework.Resource
	{
		c := testresource.Config{}

		tr, err = testresource.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []framework.Resource{
		tr,
	}

	handlesFunc := func(obj interface{}) bool {
		return true
	}

	var resourceSet *framework.ResourceSet
	{
		c := framework.ResourceSetConfig{
			Handles:   handlesFunc,
			Logger:    config.Logger,
			Resources: resources,
		}

		resourceSet, err = framework.NewResourceSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return resourceSet, nil
}
