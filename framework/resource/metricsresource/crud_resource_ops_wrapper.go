package metricsresource

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/giantswarm/operatorkit/framework"
)

type crudResourceOpsWrapperConfig struct {
	Ops framework.CRUDResourceOps

	ServiceName  string
	ResourceName string
}

type crudResourceWrapperOps struct {
	underlying framework.CRUDResourceOps

	serviceName  string
	resourceName string
}

func newCRUDResourceWrapperOps(config crudResourceOpsWrapperConfig) (*crudResourceWrapperOps, error) {
	if config.Ops == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Ops must not be empty")
	}

	if config.ServiceName == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.ServiceName must not be empty")
	}
	if config.ResourceName == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.ResourceName must not be empty")
	}

	o := &crudResourceWrapperOps{
		underlying: config.Ops,

		serviceName:  toCamelCase(config.ServiceName),
		resourceName: config.ResourceName,
	}

	return o, nil
}

func (o *crudResourceWrapperOps) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	sl := o.serviceName
	rl := o.resourceName
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
	sl := o.serviceName
	rl := o.resourceName
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
	sl := o.serviceName
	rl := o.resourceName
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
	sl := o.serviceName
	rl := o.resourceName
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
	return o.resourceName
}

func (o *crudResourceWrapperOps) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	sl := o.serviceName
	rl := o.resourceName
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
	sl := o.serviceName
	rl := o.resourceName
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
	sl := o.serviceName
	rl := o.resourceName
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
