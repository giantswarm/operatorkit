package crud

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/micrologger/loggermeta"

	"github.com/giantswarm/operatorkit/controller/context/reconciliationcanceledcontext"
	"github.com/giantswarm/operatorkit/controller/context/resourcecanceledcontext"
)

type ResourceConfig struct {
	// CRUD is any CRUD Resource implementation wrapped and properly executed by
	// Resource.
	CRUD   Interface
	Logger micrologger.Logger
}

// Resource wraps Interface and allows the implemention of complex CRUD
// Resrouces in a structured way. The usual control flow primitives like
// canceling the resource itself or canceling the whole reconciliation loop are
// available. Further benefits are dedicated metrics for the several CRUD
// Resource steps defined by Interface when wrapped by the metricsresource
// package.
type Resource struct {
	crud   Interface
	logger micrologger.Logger
}

func NewResource(config ResourceConfig) (*Resource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.CRUD == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.CRUD must not be empty", config)
	}
	if config.CRUD.Name() == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.CRUD.Name() must not be empty", config)
	}

	r := &Resource{
		crud:   config.CRUD,
		logger: config.Logger,
	}

	return r, nil
}

// CRUD is not part of the resource.Interface and should not be used. It is
// exposed for wrapping purposes in the wrapper package.
//
// NOTE This method should not be used outside operatorkit.
//
func (r *Resource) CRUD() Interface {
	return r.crud
}

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	var err error

	var currentState interface{}
	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
			return nil
		}

		meta, ok := loggermeta.FromContext(ctx)
		if ok {
			meta.KeyVals["function"] = "GetCurrentState"
			defer delete(meta.KeyVals, "function")
		}
		currentState, err = r.crud.GetCurrentState(ctx, obj)
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
			return nil
		}

		meta, ok := loggermeta.FromContext(ctx)
		if ok {
			meta.KeyVals["function"] = "GetDesiredState"
			defer delete(meta.KeyVals, "function")
		}
		desiredState, err = r.crud.GetDesiredState(ctx, obj)
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
			return nil
		}

		meta, ok := loggermeta.FromContext(ctx)
		if ok {
			meta.KeyVals["function"] = "NewUpdatePatch"
			defer delete(meta.KeyVals, "function")
		}
		patch, err = r.crud.NewUpdatePatch(ctx, obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
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
				err := r.crud.ApplyCreateChange(ctx, obj, createState)
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
				err := r.crud.ApplyDeleteChange(ctx, obj, deleteState)
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
				err := r.crud.ApplyUpdateChange(ctx, obj, updateState)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}
	}

	return nil
}

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	var err error

	var currentState interface{}
	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
			return nil
		}

		meta, ok := loggermeta.FromContext(ctx)
		if ok {
			meta.KeyVals["function"] = "GetCurrentState"
			defer delete(meta.KeyVals, "function")
		}
		currentState, err = r.crud.GetCurrentState(ctx, obj)
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
			return nil
		}

		meta, ok := loggermeta.FromContext(ctx)
		if ok {
			meta.KeyVals["function"] = "GetDesiredState"
			defer delete(meta.KeyVals, "function")
		}
		desiredState, err = r.crud.GetDesiredState(ctx, obj)
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
			return nil
		}

		meta, ok := loggermeta.FromContext(ctx)
		if ok {
			meta.KeyVals["function"] = "NewDeletePatch"
			defer delete(meta.KeyVals, "function")
		}
		patch, err = r.crud.NewDeletePatch(ctx, obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	{
		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
		if resourcecanceledcontext.IsCanceled(ctx) {
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
				err := r.crud.ApplyCreateChange(ctx, obj, createChange)
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
				err := r.crud.ApplyDeleteChange(ctx, obj, deleteChange)
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
				err := r.crud.ApplyUpdateChange(ctx, obj, updateChange)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}
	}

	return nil
}

func (r *Resource) Name() string {
	return r.crud.Name()
}
