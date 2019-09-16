package metricsresource

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/giantswarm/operatorkit/resource"
)

type resourceWrapper struct {
	resource resource.Interface
}

func newResourceWrapper(config Config) (*resourceWrapper, error) {
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
	}

	r := &resourceWrapper{
		resource: config.Resource,
	}

	return r, nil
}

func (r *resourceWrapper) EnsureCreated(ctx context.Context, obj interface{}) error {
	rl := r.resource.Name()
	ol := "EnsureCreated"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	err := r.resource.EnsureCreated(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *resourceWrapper) EnsureDeleted(ctx context.Context, obj interface{}) error {
	rl := r.resource.Name()
	ol := "EnsureDeleted"

	operationCounter.WithLabelValues(rl, ol).Inc()

	t := prometheus.NewTimer(operationHistogram.WithLabelValues(rl, ol))
	defer t.ObserveDuration()

	err := r.resource.EnsureDeleted(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *resourceWrapper) Name() string {
	return r.resource.Name()
}

// Wrapped implements internal.Wrapper interface.
func (r *resourceWrapper) Wrapped() resource.Interface {
	return r.resource
}
