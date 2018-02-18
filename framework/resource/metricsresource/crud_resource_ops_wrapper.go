package metricsresource

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/resource/internal"
)

type crudResourceOpsWrapperConfig struct {
	Ops framework.CRUDResourceOps

	Name string
}

type crudResourceWrapperOps struct {
	underlying framework.CRUDResourceOps

	name string
}

func newCRUDResourceWrapperOps(config crudResourceOpsWrapperConfig) (*crudResourceWrapperOps, error) {
	if config.Ops == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Ops must not be empty")
	}

	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Name must not be empty")
	}

	o := &crudResourceWrapperOps{
		underlying: config.Ops,

		name: toCamelCase(config.Name),
	}

	return o, nil
}

func (o *crudResourceWrapperOps) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	sl := o.name
	rl := o.underlying.Name()
	ol := "GetCurrentState"

	operationCounter.WithLabelValues(sl, rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(sl, rl, ol))
	defer t.ObserveDuration()

	v, err := o.underlying.GetCurrentState(ctx, obj)
	if err != nil {
		operationErrorCounter.WithLabelValues(sl, rl, ol).Inc()
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (o *crudResourceWrapperOps) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	sl := o.name
	rl := o.underlying.Name()
	ol := "GetDesiredState"

	operationCounter.WithLabelValues(sl, rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(sl, rl, ol))
	defer t.ObserveDuration()

	v, err := o.underlying.GetDesiredState(ctx, obj)
	if err != nil {
		operationErrorCounter.WithLabelValues(sl, rl, ol).Inc()
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (o *crudResourceWrapperOps) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	sl := o.name
	rl := o.underlying.Name()
	ol := "NewUpdatePatch"

	operationCounter.WithLabelValues(sl, rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(sl, rl, ol))
	defer t.ObserveDuration()

	v, err := o.underlying.NewUpdatePatch(ctx, obj, currentState, desiredState)
	if err != nil {
		operationErrorCounter.WithLabelValues(sl, rl, ol).Inc()
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (o *crudResourceWrapperOps) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	sl := o.name
	rl := o.underlying.Name()
	ol := "NewDeletePatch"

	operationCounter.WithLabelValues(sl, rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(sl, rl, ol))
	defer t.ObserveDuration()

	v, err := o.underlying.NewDeletePatch(ctx, obj, currentState, desiredState)
	if err != nil {
		operationErrorCounter.WithLabelValues(sl, rl, ol).Inc()
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (o *crudResourceWrapperOps) Name() string {
	return internal.OldUnderlying(o).Name()
}

func (o *crudResourceWrapperOps) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	sl := o.name
	rl := o.underlying.Name()
	ol := "ApplyCreatePatch"

	operationCounter.WithLabelValues(sl, rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(sl, rl, ol))
	defer t.ObserveDuration()

	err := o.underlying.ApplyCreateChange(ctx, obj, createState)
	if err != nil {
		operationErrorCounter.WithLabelValues(sl, rl, ol).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (o *crudResourceWrapperOps) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	sl := o.name
	rl := o.underlying.Name()
	ol := "ApplyDeletePatch"

	operationCounter.WithLabelValues(sl, rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(sl, rl, ol))
	defer t.ObserveDuration()

	err := o.underlying.ApplyDeleteChange(ctx, obj, deleteState)
	if err != nil {
		operationErrorCounter.WithLabelValues(sl, rl, ol).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (o *crudResourceWrapperOps) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	sl := o.name
	rl := o.underlying.Name()
	ol := "ApplyUpdatePatch"

	operationCounter.WithLabelValues(sl, rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(sl, rl, ol))
	defer t.ObserveDuration()

	err := o.underlying.ApplyUpdateChange(ctx, obj, updateState)
	if err != nil {
		operationErrorCounter.WithLabelValues(sl, rl, ol).Inc()
		return microerror.Mask(err)
	}

	return nil
}
