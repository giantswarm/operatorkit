package metricsresource

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/giantswarm/operatorkit/v4/pkg/resource"
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

	reportLastReconciled(obj)
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

func reportLastReconciled(obj interface{}) {
	var kind string
	{
		obj, ok := obj.(runtime.Object)
		if !ok {
			return
		}

		kind = obj.GetObjectKind().GroupVersionKind().Kind
	}

	var name, namespace string
	{
		obj, ok := obj.(metav1.Object)
		if !ok {
			return
		}

		name = obj.GetName()
		namespace = obj.GetNamespace()
	}

	lastReconciledGauge.WithLabelValues(
		kind,
		name,
		namespace,
	).SetToCurrentTime()
}
