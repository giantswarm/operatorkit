package framework

import "context"

// Resource implements the building blocks of any resource business logic being
// reconciled when observing custom resources. This interface provides
// a guideline for an easier way to follow the rather complex intentions of
// operators in general.
type Resource interface {
	// Name returns the resource's name used for identification.
	Name() string
	// Underlying returns the underlying resource which is wrapped by the calling
	// resource. Underlying must always return a non nil resource. Otherwise
	// proper resource chaining and execution cannot be guaranteed. In case a
	// resource does not wrap any other resource, Underlying must return the
	// resource that does not wrap any resource. The returned resource is then the
	// origin, the underlying resource of the chain. In combination with Name,
	// Underlying can be used for proper identification.
	Underlying() Resource

	// GetCurrentState receives the custom object observed during custom
	// resource watches. Its purpose is to return the current state of the
	// resources being managed by the operator. This can e.g. be some
	// actual data within a configmap as provided by the Kubernetes API.
	// This is not limited to Kubernetes resources though. Another example
	// would be to fetch and return information about Flannel bridges.
	//
	// NOTE GetCurrentState is called on create, delete and update events. When
	// called on create and delete events the provided custom object will be the
	// custom object currently known to the informer. On update events the
	// informer knows about the old and the new custom object. GetCurrentState
	// then receives the new custom object to be able to consume the current state
	// of a system.
	GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error)
	// GetDesiredState receives the custom object observed during custom
	// resource watches. Its purpose is to return the desired state of the
	// resources being managed by the operator. The desired state should
	// always be able to be made up using the information provided by the
	// custom object. This can e.g. be some data within a configmap, how it
	// should be provided by the Kubernetes API. This is not limited to
	// Kubernetes resources though. Another example would be to make up and
	// return information about Flannel bridges, how they should look like
	// on a server host.
	//
	// NOTE GetDesiredState is called on create, delete and update events.
	// When called on create events the provided custom object will be the
	// custom object currently known to the informer. On update events the
	// informer knows about the old and the new custom object.
	// GetDesiredState then receives the new custom object to be able to
	// compute the desired state of a system.
	GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error)

	// NewUpdatePatch is callend upon observed custom object change. It receives
	// the observed custom object, the current state as provided by
	// GetCurrentState and the desired state as provided by
	// GetDesiredState. NewUpdatePatch analyses the current and desired
	// state and returns the patch to be applied by Create, Delete, and
	// Update functions. ApplyCreatePatch, ApplyDeletePatch, and
	// ApplyUpdatePatch are called only when the corresponding patch part
	// was created.
	NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error)
	// NewDeletePatch is called upon observed custom object deletion. It
	// receives the deleted custom object, the current state as provided by
	// GetCurrentState and the desired state as provided by
	// GetDesiredState. NewDeletePatch analyses the current and desired
	// state returns the patch to be applied by Create, Delete, and Update
	// functions. ApplyCreatePatch, ApplyDeletePatch, and
	// ApplyUpdatePatch are called only when the corresponding patch part
	// was created.
	NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error)

	// ApplyCreatePatch receives the new custom object observed during
	// custom resource watches. It also receives the ApplyCreatePatch
	// portion of the Patch provided by NewUpdatePatch or NewDeletePatch.
	// ApplyCreatePatch only has to create resources based on its provided
	// input. All other reconciliation logic and state transformation is
	// already done at this point of the reconciliation loop.
	ApplyCreatePatch(ctx context.Context, obj, createPatch interface{}) error
	// ApplyDeletePatch receives the new custom object observed during
	// custom resource watches. It also receives the ApplyDeletePatch
	// portion of the Patch provided by NewUpdatePatch or NewDeletePatch.
	// ApplyDeletePatch only has to delete resources based on its provided
	// input. All other reconciliation logic and state transformation is
	// already done at this point of the reconciliation loop.
	ApplyDeletePatch(ctx context.Context, obj, deletePatch interface{}) error
	// ApplyUpdatePatch receives the new custom object observed during
	// custom resource watches. It also receives the state intended to be
	// updated as provided by NewUpdatePatch or NewDeletePatch.
	// ApplyUpdatePatch has to update resources based on its provided
	// input. All other reconciliation logic and state transformation is
	// already done at this point of the reconciliation loop.
	ApplyUpdatePatch(ctx context.Context, obj, updatePatch interface{}) error
}
