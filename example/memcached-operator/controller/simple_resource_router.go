package controller

import (
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/example/memcached-operator/logger"
)

// newSimpleResourceRouter creates a resource router which handles all objects
// with single set of resources.
func newSimpleResourceRouter(resources []controller.Resource) (*controller.ResourceRouter, error) {
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

	var router *controller.ResourceRouter
	{
		c := controller.ResourceRouterConfig{
			Logger: logger.Default,

			ResourceSets: []*controller.ResourceSet{
				set,
			},
		}

		router, err = controller.NewResourceRouter(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return router, nil
}
