package metricsresource

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/giantswarm/operatorkit/resource/crud"
)

type crudResourceConfig struct {
	CRUD crud.Interface
}

type crudResource struct {
	crud crud.Interface
}

func newCRUDResource(config crudResourceConfig) (*crudResource, error) {
	if config.CRUD == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.CRUD must not be empty", config)
	}

	r := &crudResource{
		crud: config.CRUD,
	}

	return r, nil
}

func (r *crudResource) Name() string {
	return r.crud.Name()
}

func (r *crudResource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	rl := r.crud.Name()
	ol := "GetCurrentState"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	v, err := r.crud.GetCurrentState(ctx, obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *crudResource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	rl := r.crud.Name()
	ol := "GetDesiredState"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	v, err := r.crud.GetDesiredState(ctx, obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *crudResource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*crud.Patch, error) {
	rl := r.crud.Name()
	ol := "NewUpdatePatch"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	v, err := r.crud.NewUpdatePatch(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *crudResource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*crud.Patch, error) {
	rl := r.crud.Name()
	ol := "NewDeletePatch"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	v, err := r.crud.NewDeletePatch(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *crudResource) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	rl := r.crud.Name()
	ol := "ApplyCreatePatch"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	err := r.crud.ApplyCreateChange(ctx, obj, createState)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *crudResource) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	rl := r.crud.Name()
	ol := "ApplyDeletePatch"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	err := r.crud.ApplyDeleteChange(ctx, obj, deleteState)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *crudResource) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	rl := r.crud.Name()
	ol := "ApplyUpdatePatch"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	err := r.crud.ApplyUpdateChange(ctx, obj, updateState)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
