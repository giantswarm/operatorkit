package framework

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/micrologger/loggermeta"
	"github.com/giantswarm/operatorkit/framework/context/reconciliationcanceledcontext"
	"github.com/giantswarm/operatorkit/framework/context/resourcecanceledcontext"
)

type CRUDResourceConfig struct {
	Logger micrologger.Logger
	// Ops is a set of operations used by CRUDResource to implement the
	// Resource interface.
	Ops CRUDResourceOps

	// Name is the resource's name used for identification.
	Name string
}

// CRUDResource allows implementing complex CRUD Resrouces in structured way.
// Besides that is implements various context features defined in subpackages
// of the context package.
type CRUDResource struct {
	CRUDResourceOps

	logger micrologger.Logger

	name string
}

func NewCRUDResource(config CRUDResourceConfig) (*CRUDResource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Ops == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Ops must not be empty")
	}

	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Name must not be empty")
	}

	r := &CRUDResource{
		CRUDResourceOps: config.Ops,

		logger: config.Logger,

		name: config.Name,
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
		currentState, err = r.GetCurrentState(ctx, obj)
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
		desiredState, err = r.GetDesiredState(ctx, obj)
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
		patch, err = r.NewUpdatePatch(ctx, obj, currentState, desiredState)
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
				err := r.ApplyCreateChange(ctx, obj, createState)
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
				err := r.ApplyDeleteChange(ctx, obj, deleteState)
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
				err := r.ApplyUpdateChange(ctx, obj, updateState)
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
		currentState, err = r.GetCurrentState(ctx, obj)
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
		desiredState, err = r.GetDesiredState(ctx, obj)
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
		patch, err = r.NewDeletePatch(ctx, obj, currentState, desiredState)
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
				err := r.ApplyCreateChange(ctx, obj, createChange)
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
				err := r.ApplyDeleteChange(ctx, obj, deleteChange)
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
				err := r.ApplyUpdateChange(ctx, obj, updateChange)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}
	}

	return nil
}

func (r *CRUDResource) Name() string {
	return r.name
}
