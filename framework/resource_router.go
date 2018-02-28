package framework

import (
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/apimachinery/pkg/api/meta"
)

type ResourceRouterConfig struct {
	Logger micrologger.Logger

	ResourceSets []*ResourceSet
}

type ResourceRouter struct {
	logger micrologger.Logger

	resourceSets []*ResourceSet
}

func NewResourceRouter(c ResourceRouterConfig) (*ResourceRouter, error) {
	if c.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", c)
	}

	if len(c.ResourceSets) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.ResourceSets must not be empty", c)
	}

	r := &ResourceRouter{
		logger: c.Logger,

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

	if len(found) == 0 {
		accessor, err := meta.Accessor(obj)
		if err != nil {
			r.logger.Log("function", "ResourceSet", "level", "warning", "message", "cannot create accessor for object", "object", fmt.Sprintf("%#v", obj), "stack", fmt.Sprintf("%#v", err))
		} else {
			r.logger.Log("level", "debug", "message", "no resource set for reconciled object", "object", accessor.GetSelfLink())
		}

		return nil, microerror.Mask(noResourceRouterError)
	}
	if len(found) > 1 {
		return nil, microerror.Maskf(executionFailedError, "multiple handling resource sets found; only single allowed")
	}

	return found[0], nil
}
