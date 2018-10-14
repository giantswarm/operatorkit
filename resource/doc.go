// Package resource provides a set of Kubernetes resources ready to use with
// controller package, i.e. they implement controller.Resource interface. All
// resources in the package follow the same schema. They are configured with
// name and namespace of the reconciled object and desired state function which
// returns the desired object for the given watched object. Configured name and
// namespace values have to match those set in the desired object otherwise
// reconciliation can not be performed and resource will return
// invalidDesiredStateError.
package resource
