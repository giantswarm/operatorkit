package framework

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger/loggermeta"
	"github.com/giantswarm/operatorkit/framework/context/reconciliationcanceledcontext"
	"github.com/giantswarm/operatorkit/framework/context/resourcecanceledcontext"
)

type CRUDResourceConfig struct {
	Ops CRUDResourceOps
}

type CRUDResourceOps interface {
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

	// NewUpdatePatch is called upon observed custom object change. It receives
	// the observed custom object, the current state as provided by
	// GetCurrentState and the desired state as provided by
	// GetDesiredState. NewUpdatePatch analyses the current and desired
	// state and returns the patch to be applied by Create, Delete, and
	// Update functions. ApplyCreateChange, ApplyDeleteChange, and
	// ApplyUpdateChange are called only when the corresponding patch part
	// was created.
	NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error)
	// NewDeletePatch is called upon observed custom object deletion. It
	// receives the deleted custom object, the current state as provided by
	// GetCurrentState and the desired state as provided by
	// GetDesiredState. NewDeletePatch analyses the current and desired
	// state returns the patch to be applied by Create, Delete, and Update
	// functions. ApplyCreateChange, ApplyDeleteChange, and
	// ApplyUpdateChange are called only when the corresponding patch part
	// was created.
	NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error)

	// ApplyCreateChange receives the new custom object observed during
	// custom resource watches. It also receives the create portion of the
	// Patch provided by NewUpdatePatch or NewDeletePatch.
	// ApplyCreateChange only has to create resources based on its provided
	// input. All other reconciliation logic and state transformation is
	// already done at this point of the reconciliation loop.
	ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error
	// ApplyDeleteChange receives the new custom object observed during
	// custom resource watches. It also receives the delete portion of the
	// Patch provided by NewUpdatePatch or NewDeletePatch.
	// ApplyDeleteChange only has to delete resources based on its provided
	// input. All other reconciliation logic and state transformation is
	// already done at this point of the reconciliation loop.
	ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error
	// ApplyUpdateChange receives the new custom object observed during
	// custom resource watches. It also receives the update portion of the
	// Patch provided by NewUpdatePatch or NewDeletePatch.
	// ApplyUpdateChange has to update resources based on its provided
	// input. All other reconciliation logic and state transformation is
	// already done at this point of the reconciliation loop.
	ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error
}

type CRUDResource struct {
	ops CRUDResourceOps
}

func NewCRUDResource(config CRUDResourceConfig) (*CRUDResource, error) {
	if config.Ops == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Ops must not be empty")
	}
	if config.Ops.Name() == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Ops.Name() must not be empty")
	}
	if config.Ops.Underlying() == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Ops.Underlying() must not be empty")
	}

	r := &CRUDResource{
		ops: config.Ops,
	}

	return r, nil
}

func (r *CRUDResource) EnsureCreated(ctx context.Context, obj interface{}) error {
	var err error

	var currentState interface{}
	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
			ctx = resourcecanceledcontext.NewContext(ctx, make(chan struct{}))
			return nil
		}

		meta, ok := loggermeta.FromContext(ctx)
		if ok {
			meta.KeyVals["function"] = "GetCurrentState"
			defer delete(meta.KeyVals, "function")
		}
		currentState, err = r.ops.GetCurrentState(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var desiredState interface{}
	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
			ctx = resourcecanceledcontext.NewContext(ctx, make(chan struct{}))
			return nil
		}

		meta, ok := loggermeta.FromContext(ctx)
		if ok {
			meta.KeyVals["function"] = "GetDesiredState"
			defer delete(meta.KeyVals, "function")
		}
		desiredState, err = r.ops.GetDesiredState(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var patch *Patch
	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
			ctx = resourcecanceledcontext.NewContext(ctx, make(chan struct{}))
			return nil
		}

		meta, ok := loggermeta.FromContext(ctx)
		if ok {
			meta.KeyVals["function"] = "NewUpdatePatch"
			defer delete(meta.KeyVals, "function")
		}
		patch, err = r.ops.NewUpdatePatch(ctx, obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
			ctx = resourcecanceledcontext.NewContext(ctx, make(chan struct{}))
			return nil
		}

		if patch != nil {
			createState, ok := patch.getCreateChange()
			if ok {
				meta, ok := loggermeta.FromContext(ctx)
				if ok {
					meta.KeyVals["function"] = "ApplyCreateChange"
					defer delete(meta.KeyVals, "function")
				}
				err := r.ops.ApplyCreateChange(ctx, obj, createState)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}
	}

	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
			ctx = resourcecanceledcontext.NewContext(ctx, make(chan struct{}))
			return nil
		}

		if patch != nil {
			deleteState, ok := patch.getDeleteChange()
			if ok {
				meta, ok := loggermeta.FromContext(ctx)
				if ok {
					meta.KeyVals["function"] = "ApplyDeleteChange"
					defer delete(meta.KeyVals, "function")
				}
				err := r.ops.ApplyDeleteChange(ctx, obj, deleteState)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}
	}

	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
			ctx = resourcecanceledcontext.NewContext(ctx, make(chan struct{}))
			return nil
		}

		if patch != nil {
			updateState, ok := patch.getUpdateChange()
			if ok {
				meta, ok := loggermeta.FromContext(ctx)
				if ok {
					meta.KeyVals["function"] = "ApplyUpdateChange"
					defer delete(meta.KeyVals, "function")
				}
				err := r.ops.ApplyUpdateChange(ctx, obj, updateState)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}
	}

	return nil
}

func (r *CRUDResource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	var err error

	var currentState interface{}
	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
			ctx = resourcecanceledcontext.NewContext(ctx, make(chan struct{}))
			return nil
		}

		meta, ok := loggermeta.FromContext(ctx)
		if ok {
			meta.KeyVals["function"] = "GetCurrentState"
			defer delete(meta.KeyVals, "function")
		}
		currentState, err = r.ops.GetCurrentState(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var desiredState interface{}
	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
			ctx = resourcecanceledcontext.NewContext(ctx, make(chan struct{}))
			return nil
		}

		meta, ok := loggermeta.FromContext(ctx)
		if ok {
			meta.KeyVals["function"] = "GetDesiredState"
			defer delete(meta.KeyVals, "function")
		}
		desiredState, err = r.ops.GetDesiredState(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var patch *Patch
	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
			ctx = resourcecanceledcontext.NewContext(ctx, make(chan struct{}))
			return nil
		}

		meta, ok := loggermeta.FromContext(ctx)
		if ok {
			meta.KeyVals["function"] = "NewDeletePatch"
			defer delete(meta.KeyVals, "function")
		}
		patch, err = r.ops.NewDeletePatch(ctx, obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
			ctx = resourcecanceledcontext.NewContext(ctx, make(chan struct{}))
			return nil
		}

		if patch != nil {
			createChange, ok := patch.getCreateChange()
			if ok {
				meta, ok := loggermeta.FromContext(ctx)
				if ok {
					meta.KeyVals["function"] = "ApplyCreateChange"
					defer delete(meta.KeyVals, "function")
				}
				err := r.ops.ApplyCreateChange(ctx, obj, createChange)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}
	}

	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
			ctx = resourcecanceledcontext.NewContext(ctx, make(chan struct{}))
			return nil
		}

		if patch != nil {
			deleteChange, ok := patch.getDeleteChange()
			if ok {
				meta, ok := loggermeta.FromContext(ctx)
				if ok {
					meta.KeyVals["function"] = "ApplyDeleteChange"
					defer delete(meta.KeyVals, "function")
				}
				err := r.ops.ApplyDeleteChange(ctx, obj, deleteChange)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}
	}

	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
			ctx = resourcecanceledcontext.NewContext(ctx, make(chan struct{}))
			return nil
		}

		if patch != nil {
			updateChange, ok := patch.getUpdateChange()
			if ok {
				meta, ok := loggermeta.FromContext(ctx)
				if ok {
					meta.KeyVals["function"] = "ApplyUpdateChange"
					defer delete(meta.KeyVals, "function")
				}
				err := r.ops.ApplyUpdateChange(ctx, obj, updateChange)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}
	}

	return nil
}
