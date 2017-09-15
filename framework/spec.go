package framework

import "context"

// Resource implements the building blocks of any resource business logic being
// reconciled when observing TPRs. This interface provides a guideline for an
// easier way to follow the rather complex intentions of operators in general.
type Resource interface {
	// GetCurrentState receives the custom object observed during TPR watches. Its
	// purpose is to return the current state of the resources being managed by
	// the operator. This can e.g. be some actual data within a configmap as
	// provided by the Kubernetes API. This is not limited to Kubernetes resources
	// though. Another example would be to fetch and return information about
	// Flannel bridges.
	//
	// NOTE GetCurrentState is called on create, delete and update events. When
	// called on create and delete events the provided custom object will be the
	// custom object currently known to the informer. On update events the
	// informer knows about the old and the new custom object. GetCurrentState
	// then receives the new custom object to be able to consume the current state
	// of a system.
	GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error)
	// GetDesiredState receives the custom object observed during TPR watches. Its
	// purpose is to return the desired state of the resources being managed by
	// the operator. The desired state should always be able to be made up using
	// the information provided by the TPO. This can e.g. be some data within a
	// configmap, how it should be provided by the Kubernetes API. This is not
	// limited to Kubernetes resources though. Another example would be to make up
	// and return information about Flannel bridges, how they should look like on
	// a server host.
	//
	// NOTE GetDesiredState is called on create, delete and update events. When
	// called on create and delete events the provided custom object will be the
	// custom object currently known to the informer. On update events the
	// informer knows about the old and the new custom object. GetDesiredState
	// then receives the new custom object to be able to compute the desired state
	// of a system.
	GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error)
	// GetCreateState receives the custom object observed during TPR watches. It
	// also receives the current state as provided by GetCurrentState and the
	// desired state as provided by GetDesiredState. GetCreateState analyses the
	// current and desired state and returns the state intended to be created by
	// ProcessCreateState.
	GetCreateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error)
	// GetDeleteState receives the custom object observed during TPR watches. It
	// also receives the current state as provided by GetCurrentState and the
	// desired state as provided by GetDesiredState. GetDeleteState analyses the
	// current and desired state and returns the state intended to be deleted by
	// ProcessDeleteState.
	GetDeleteState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error)
	// GetUpdateState receives the new custom object observed during TPR watches.
	// It also receives the current state as provided by GetCurrentState and the
	// desired state as provided by GetDesiredState. GetUpdateState analyses the
	// current and desired state and returns the states intended to be created,
	// deleted and updated. The returned create state will be given to
	// ProcessCreateState. The returned delete state will be given to
	// ProcessDeleteState. The returned update state will be given to
	// ProcessUpdateState.
	//
	// NOTE simple resources not concerned with being updated do not have to
	// implement anything but just fulfil the resource interface. More complex
	// resources, e.g. those managing multiple entities of themselves at once may
	// require a more complex update mechanism. Then multiple entities might be
	// added, removed and/or modified over the course of the resource's lifecycle.
	// This transformation has to be reflected by different states which are
	// returned by GetUpdateState. The first value being returned is the create
	// state, the second the delete state and the third the update state.
	GetUpdateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, interface{}, interface{}, error)
	// Name returns the resource's name used for identification.
	Name() string
	// ProcessCreateState receives the new custom object observed during TPR
	// watches. It also receives the state intended to be created as provided by
	// GetCreateState. ProcessCreateState only has to create resources based on
	// its provided input. All other reconciliation logic and state transformation
	// is already done at this point of the reconciliation loop.
	//
	// NOTE ProcessCreateState is called on create and update events. When called
	// on create events the provided custom object will be the custom object
	// currently known to the informer. On update events the informer knows about
	// the old and the new custom object. ProcessCreateState then receives the new
	// custom object to be able to process the create state of a system.
	ProcessCreateState(ctx context.Context, obj, createState interface{}) error
	// ProcessDeleteState receives the new custom object observed during TPR
	// watches. It also receives the state intended to be deleted as provided by
	// GetDeleteState. ProcessDeleteState only has to delete resources based on
	// its provided input. All other reconciliation logic and state transformation
	// is already done at this point of the reconciliation loop.
	//
	// NOTE ProcessDeleteState is called on delete and update events. When called
	// on delete events the provided custom object will be the custom object
	// currently known to the informer. On update events the informer knows about
	// the old and the new custom object. ProcessDeleteState then receives the new
	// custom object to be able to process the delete state of a system.
	ProcessDeleteState(ctx context.Context, obj, deleteState interface{}) error
	// ProcessUpdateState receives the new custom object observed during TPR
	// watches. It also receives the state intended to be updated as provided by
	// GetUpdateState. ProcessUpdateState has to update resources based on its
	// provided input. All other reconciliation logic and state transformation is
	// already done at this point of the reconciliation loop.
	ProcessUpdateState(ctx context.Context, obj, updateState interface{}) error
	// Underlying returns the underlying resource which is wrapped by the calling
	// resource. Underlying must always return a non nil resource. Otherwise
	// proper resource chaining and execution cannot be guaranteed. In case a
	// resource does not wrap any other resource, Underlying must return the
	// resource that does not wrap any resource. The returned resource is then the
	// origin, the underlying resource of the chain. In combination with Name,
	// Underlying can be used for proper identification.
	Underlying() Resource
}
