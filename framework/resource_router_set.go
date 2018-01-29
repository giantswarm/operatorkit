package framework

import (
	"github.com/giantswarm/microerror"
)

type ResourceRouterSetConfig struct {
	ResourceRouters []*ResourceRouter
}

type ResourceRouterSet struct {
	resourceRouters []*ResourceRouter
}

func NewResourceRouterSet(c ResourceRouterSetConfig) (*ResourceRouterSet, error) {
	if len(c.ResourceRouters) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.ResourceRouters must not be empty", c)
	}

	r := &ResourceRouterSet{
		resourceRouters: c.ResourceRouters,
	}

	return r, nil
}

// VersionedResourceRouter tries to lookup the versioned resource router based
// on the received custom object.
func (r *ResourceRouterSet) VersionedResourceRouter(obj interface{}) (*ResourceRouter, error) {
	var found []*ResourceRouter

	for _, router := range r.resourceRouters {
		v, err := router.CustomObjectVersionFunc()(obj)
		if IsCustomObjectVersionNotFound(err) {
			continue
		} else if err != nil {
			return nil, microerror.Mask(err)
		}
		if router.VersionBundleVersion() == v {
			found = append(found, router)
		}
	}

	if len(found) == 0 {
		return nil, microerror.Maskf(executionFailedError, "could not find any resource router")
	}
	if len(found) > 1 {
		return nil, microerror.Maskf(executionFailedError, "multiple resource routers not allowed")
	}

	return found[0], nil
}
