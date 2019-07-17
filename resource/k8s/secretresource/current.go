package secretresource

import (
	"context"

	"github.com/giantswarm/microerror"
)

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	state, err := r.stateGetter.GetCurrentState(ctx, obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return state, nil
}
