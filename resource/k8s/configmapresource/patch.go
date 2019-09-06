package configmapresource

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/operatorkit/controller"
)

func (r *Resource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*controller.Patch, error) {
	delete, err := r.newDeleteChangeForDeletePatch(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := controller.NewPatch()
	patch.SetDeleteChange(delete)

	return patch, nil
}

func (r *Resource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*controller.Patch, error) {
	create, err := r.newCreateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	delete, err := r.newDeleteChangeForUpdatePatch(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	update, err := r.newUpdateChange(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	patch := controller.NewPatch()
	patch.SetCreateChange(create)
	patch.SetDeleteChange(delete)
	patch.SetUpdateChange(update)

	return patch, nil
}
