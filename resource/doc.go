// Package resource provides a set of Kubernetes resources ready to use with
// controller package, i.e. they implement controller.Resource interface. All
// resources in the package follow the same schema. They are configured with
// desired state function which returns the desired object set for the given
// carnation of watched object.
//
// All objects created by this package's resources have
// "operatorkit.giantswarm.io/resource" annotation set. It is used during the
// deletion phase to make sure all objects created by the resource are properly
// garbage collected.
//
// NOTE: Uninitialized objects (i.e. objects having initializers set) are
// deleted during the deletion phase.
package resource
