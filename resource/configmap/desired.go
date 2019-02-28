package configmap

import (
	"context"

	"github.com/giantswarm/microerror"
)

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	state, err := r.stateGetter.GetDesiredState(ctx, obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return state, nil
}
