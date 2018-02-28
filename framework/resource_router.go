package framework

import (
	"github.com/giantswarm/microerror"
)

type ResourceRouterConfig struct {
	ResourceSets []*ResourceSet
}

type ResourceRouter struct {
	resourceSets []*ResourceSet
}

func NewResourceRouter(c ResourceRouterConfig) (*ResourceRouter, error) {
	if len(c.ResourceSets) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.ResourceSets must not be empty", c)
	}

	r := &ResourceRouter{
		resourceSets: c.ResourceSets,
	}

	return r, nil
}

// ResourceSet tries to lookup the appropriate resource set based on the
// received runtime object. There might be not any resource set for an observed
// runtime object if an operator uses multiple frameworks for reconciliations.
// There must not be multiple resource sets per observed runtime object though.
// If this is the case, ResourceSet returns an error.
func (r *ResourceRouter) ResourceSet(obj interface{}) (*ResourceSet, error) {
	var found []*ResourceSet

	for _, router := range r.resourceSets {
		if router.Handles(obj) {
			found = append(found, router)
		}
	}

	if len(found) > 1 {
		return nil, microerror.Maskf(executionFailedError, "multiple handling resource sets found; only single allowed")
	}

	return found[0], nil
}
