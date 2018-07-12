package controller

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/example/memcached-operator/logger"
)

// newSimpleResourceSets creates a list if resource sets which handles all
// reconciled objects. In this case with only a single resource set.
func newSimpleResourceSets(resources []controller.Resource) ([]*controller.ResourceSet, error) {
	var err error

	var set *controller.ResourceSet
	{
		c := controller.ResourceSetConfig{
			Logger: logger.Default,

			// Handle all objects.
			Handles: func(obj interface{}) bool { return true },

			Resources: resources,
		}

		set, err = controller.NewResourceSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resourceSets := []*controller.ResourceSet{
		set,
	}

	return resourceSets, nil
}
