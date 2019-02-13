package metricsresource

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/giantswarm/operatorkit/controller"
)

type crudResourceOpsWrapperConfig struct {
	Ops controller.CRUDResourceOps

	ResourceName string
}

type crudResourceWrapperOps struct {
	underlying controller.CRUDResourceOps

	resourceName string
}

func newCRUDResourceWrapperOps(config crudResourceOpsWrapperConfig) (*crudResourceWrapperOps, error) {
	if config.Ops == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Ops must not be empty", config)
	}

	if config.ResourceName == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.ResourceName must not be empty", config)
	}

	o := &crudResourceWrapperOps{
		underlying: config.Ops,

		resourceName: config.ResourceName,
	}

	return o, nil
}

func (o *crudResourceWrapperOps) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	rl := o.resourceName
	ol := "GetCurrentState"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	v, err := o.underlying.GetCurrentState(ctx, obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (o *crudResourceWrapperOps) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	rl := o.resourceName
	ol := "GetDesiredState"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	v, err := o.underlying.GetDesiredState(ctx, obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (o *crudResourceWrapperOps) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*controller.Patch, error) {
	rl := o.resourceName
	ol := "NewUpdatePatch"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	v, err := o.underlying.NewUpdatePatch(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (o *crudResourceWrapperOps) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*controller.Patch, error) {
	rl := o.resourceName
	ol := "NewDeletePatch"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	v, err := o.underlying.NewDeletePatch(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (o *crudResourceWrapperOps) Name() string {
	return o.resourceName
}

func (o *crudResourceWrapperOps) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	rl := o.resourceName
	ol := "ApplyCreatePatch"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	err := o.underlying.ApplyCreateChange(ctx, obj, createState)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (o *crudResourceWrapperOps) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	rl := o.resourceName
	ol := "ApplyDeletePatch"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	err := o.underlying.ApplyDeleteChange(ctx, obj, deleteState)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (o *crudResourceWrapperOps) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	rl := o.resourceName
	ol := "ApplyUpdatePatch"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	err := o.underlying.ApplyUpdateChange(ctx, obj, updateState)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
