package framework

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
	GetCurrentState(obj interface{}) (interface{}, error)
	// GetDesiredState receives the custom object observed during TPR watches. Its
	// purpose is to return the desired state of the resources being managed by
	// the operator. The desired state should always be able to be made up using
	// the information provided by the TPO. This can e.g. be some data within a
	// configmap, how it should be provided by the Kubernetes API. This is not
	// limited to Kubernetes resources though. Another example would be to make up
	// and return information about Flannel bridges, how they should look like on
	// a server host.
	GetDesiredState(obj interface{}) (interface{}, error)
	// GetCreateState receives the custom object observed during TPR watches. It
	// also receives the current state as provided by GetCurrentState and the
	// desired state as provided by GetDesiredState. GetCreateState analyses the
	// current and desired state and returns the state intended to be created by
	// ProcessCreateState.
	GetCreateState(obj, currentState, desiredState interface{}) (interface{}, error)
	// GetDeleteState receives the custom object observed during TPR watches. It
	// also receives the current state as provided by GetCurrentState and the
	// desired state as provided by GetDesiredState. GetDeleteState analyses the
	// current and desired state and returns the state intended to be deleted by
	// ProcessDeleteState.
	GetDeleteState(obj, currentState, desiredState interface{}) (interface{}, error)
	// Name returns the resource's name used for identification.
	Name() string
	// ProcessCreateState receives the custom object observed during TPR watches.
	// It also receives the state intended to be created as provided by
	// GetCreateState. ProcessCreateState only has to create resources based on
	// its provided input. All other reconciliation logic and state transformation
	// is already done at this point of the reconciliation loop.
	ProcessCreateState(obj, createState interface{}) error
	// ProcessDeleteState receives the custom object observed during TPR watches.
	// It also receives the state intended to be deleted as provided by
	// GetDeleteState. ProcessDeleteState only has to delete resources based on
	// its provided input. All other reconciliation logic and state transformation
	// is already done at this point of the reconciliation loop.
	ProcessDeleteState(obj, deleteState interface{}) error
	// Underlying returns the underlying resource which is wrapped by the calling
	// resource. Underlying must always return a non nil resource. Otherwise
	// proper resource chaining and execution cannot be guaranteed. In case a
	// resource does not wrap any other resource, Underlying must return the
	// resource that does not wrap any resource. The returned resource is then the
	// origin, the underlying resource of the chain. In combination with Name,
	// Underlying can be used for proper identification.
	Underlying() Resource
}
