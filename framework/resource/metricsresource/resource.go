package metricsresource

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/resource/internal"
)

type Config struct {
	Resource framework.Resource

	// Name is name of the service using the reconciler framework. This may be the
	// name of the executing operator or controller. The service name will be used
	// to label metrics.
	Name string
}

type Resource struct {
	resource framework.Resource

	name string
}

func New(config Config) (*Resource, error) {
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Resource must not be empty")
	}

	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Name must not be empty")
	}

	newResource := &Resource{
		resource: config.Resource,

		name: toCamelCase(config.Name),
	}

	return newResource, nil
}

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	sl := r.name
	rl := r.resource.Name()
	ol := "GetCurrentState"

	operationCounter.WithLabelValues(sl, rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(sl, rl, ol))
	defer t.ObserveDuration()

	v, err := r.resource.GetCurrentState(ctx, obj)
	if err != nil {
		operationErrorCounter.WithLabelValues(sl, rl, ol).Inc()
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	sl := r.name
	rl := r.resource.Name()
	ol := "GetDesiredState"

	operationCounter.WithLabelValues(sl, rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(sl, rl, ol))
	defer t.ObserveDuration()

	v, err := r.resource.GetDesiredState(ctx, obj)
	if err != nil {
		operationErrorCounter.WithLabelValues(sl, rl, ol).Inc()
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *Resource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	sl := r.name
	rl := r.resource.Name()
	ol := "NewUpdatePatch"

	operationCounter.WithLabelValues(sl, rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(sl, rl, ol))
	defer t.ObserveDuration()

	v, err := r.resource.NewUpdatePatch(ctx, obj, currentState, desiredState)
	if err != nil {
		operationErrorCounter.WithLabelValues(sl, rl, ol).Inc()
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *Resource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	sl := r.name
	rl := r.resource.Name()
	ol := "NewDeletePatch"

	operationCounter.WithLabelValues(sl, rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(sl, rl, ol))
	defer t.ObserveDuration()

	v, err := r.resource.NewDeletePatch(ctx, obj, currentState, desiredState)
	if err != nil {
		operationErrorCounter.WithLabelValues(sl, rl, ol).Inc()
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *Resource) Name() string {
	return r.Underlying().Name()
}

func (r *Resource) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	sl := r.name
	rl := r.resource.Name()
	ol := "ApplyCreatePatch"

	operationCounter.WithLabelValues(sl, rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(sl, rl, ol))
	defer t.ObserveDuration()

	err := r.resource.ApplyCreateChange(ctx, obj, createState)
	if err != nil {
		operationErrorCounter.WithLabelValues(sl, rl, ol).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	sl := r.name
	rl := r.resource.Name()
	ol := "ApplyDeletePatch"

	operationCounter.WithLabelValues(sl, rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(sl, rl, ol))
	defer t.ObserveDuration()

	err := r.resource.ApplyDeleteChange(ctx, obj, deleteState)
	if err != nil {
		operationErrorCounter.WithLabelValues(sl, rl, ol).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	sl := r.name
	rl := r.resource.Name()
	ol := "ApplyUpdatePatch"

	operationCounter.WithLabelValues(sl, rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(sl, rl, ol))
	defer t.ObserveDuration()

	err := r.resource.ApplyUpdateChange(ctx, obj, updateState)
	if err != nil {
		operationErrorCounter.WithLabelValues(sl, rl, ol).Inc()
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) Underlying() framework.Resource {
	underlying := r.resource
	for {
		wrapper, ok := underlying.(internal.Wrapper)
		if ok {
			underlying = wrapper.Underlying()
		} else {
			return underlying
		}
	}
}
