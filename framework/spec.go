package framework

import "context"

type patchType string

const (
	patchCreate = "create"
	patchDelete = "delete"
	patchUpdate = "update"
)

// Patch is a set of information required in order to reconcile to the desired
// state. Patch is split to three parts: create, delete and update. The parts
// are passed as arguments to Resource's Create, Delete and Update functions
// respectively. Patch is guaranteed to be applied in that order (i.e. create,
// update, delete).
type Patch struct {
	data map[patchType]interface{}
}

func NewPatch() *Patch {
	return &Patch{
		data: make(map[patchType]interface{}, 3),
	}
}

func (p *Patch) getCreate() (interface{}, bool) {
	create, ok := p.data[patchCreate]
	return create, ok
}

func (p *Patch) getDelete() (interface{}, bool) {
	delete, ok := p.data[patchDelete]
	return delete, ok
}
func (p *Patch) getUpdate() (interface{}, bool) {
	update, ok := p.data[patchUpdate]
	return update, ok
}

func (p *Patch) SetCreate(create interface{}) { p.data[patchCreate] = create }
func (p *Patch) SetDelete(delete interface{}) { p.data[patchDelete] = delete }
func (p *Patch) SetUpdate(update interface{}) { p.data[patchUpdate] = update }

// Resource implements the building blocks of any resource business logic being
// reconciled when observing TPRs. This interface provides a guideline for an
// easier way to follow the rather complex intentions of operators in general.
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
	// NOTE GetDesiredState is called on create and update events. When
	// called on create events the provided custom object will be the
	// custom object currently known to the informer. On update events the
	// informer knows about the old and the new custom object.
	// GetDesiredState then receives the new custom object to be able to
	// compute the desired state of a system.
	GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error)

	// NewUpdatePatch is callend upon observed CRO/TPO change. It receives the
	// observed custom object, the current state as provided by GetCurrentState and
	// the desired state as provided by GetDesiredState. NewUpdatePatch analyses
	// the current and desired state and returns the patch state to be
	// applied by Create, Delete, and Update functions. Create, Delete, and
	// Update are called only when the corresponding patch state was
	// created.
	NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error)
	// NewDeletePatch is called upon observed CRO/TPO deleteion. It
	// receives the deleted custom object, the current state as provided by
	// GetCurrentState. NewDeletePatch analyses the current state
	// returns the patch state to be applied by Create, Delete, and Update
	// functions. Create, Delete, and Update are called only when the
	// corresponding patch state was created.
	NewDeletePatch(ctx context.Context, obj, currentState interface{}) (*Patch, error)

	// Create receives the new custom object observed during TPR watches.
	// It also receives the Create portion of the Patch provided by
	// NewUpdatePatch or NewDeletePatch. Create only has to create
	// resources based on its provided input. All other reconciliation
	// logic and state transformation is already done at this point of the
	// reconciliation loop.
	Create(ctx context.Context, obj, createPatch interface{}) error
	// Delete receives the new custom object observed during TPR watches.
	// It also receives the Delete portion of the Patch provided by
	// NewUpdatePatch or NewDeletePatch. Delete only has to delete
	// resources based on its provided input. All other reconciliation
	// logic and state transformation is already done at this point of the
	// reconciliation loop.
	Delete(ctx context.Context, obj, deletePatch interface{}) error
	// Update receives the new custom object observed during TPR watches.
	// It also receives the state intended to be updated as provided by
	// NewUpdatePatch or NewDeletePatch. Update has to update resources
	// based on its provided input. All other reconciliation logic and
	// state transformation is already done at this point of the
	// reconciliation loop.
	Update(ctx context.Context, obj, updatePatch interface{}) error
}
