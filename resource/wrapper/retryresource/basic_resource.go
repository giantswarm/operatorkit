package retryresource

import (
	"context"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/resource"
)

type basicResourceConfig struct {
	BackOff  backoff.Interface
	Logger   micrologger.Logger
	Resource resource.Interface
}

type basicResource struct {
	backOff  backoff.Interface
	logger   micrologger.Logger
	resource resource.Interface
}

func newBasicResource(config basicResourceConfig) (*basicResource, error) {
	if config.BackOff == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.BackOff must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
	}

	r := &basicResource{
		backOff:  config.BackOff,
		logger:   config.Logger,
		resource: config.Resource,
	}

	return r, nil
}

func (r *basicResource) EnsureCreated(ctx context.Context, obj interface{}) error {
	var err error

	o := func() error {
		err = r.resource.EnsureCreated(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.LogCtx(ctx, "level", "warning", "message", "retrying due to error", "stack", microerror.Stack(err))
	}

	err = backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *basicResource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	var err error

	o := func() error {
		err = r.resource.EnsureDeleted(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.LogCtx(ctx, "level", "warning", "message", "retrying due to error", "stack", microerror.Stack(err))
	}

	err = backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *basicResource) Name() string {
	return r.resource.Name()
}
