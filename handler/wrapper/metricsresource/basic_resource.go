package metricsresource

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/giantswarm/operatorkit/resource"
)

type basicResourceConfig struct {
	Resource resource.Interface
}

type basicResource struct {
	resource resource.Interface
}

func newBasicResource(config basicResourceConfig) (*basicResource, error) {
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
	}

	r := &basicResource{
		resource: config.Resource,
	}

	return r, nil
}

func (r *basicResource) EnsureCreated(ctx context.Context, obj interface{}) error {
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

func (r *basicResource) EnsureDeleted(ctx context.Context, obj interface{}) error {
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

func (r *basicResource) Name() string {
	return r.resource.Name()
}
