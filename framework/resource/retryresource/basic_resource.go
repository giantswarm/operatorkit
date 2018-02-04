package retryresource

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
)

type basicResourceConfig struct {
	Logger micrologger.Logger
	//Resource framework.Resource
	Resource interface {
		Name() string
		Underlying() framework.Resource

		EnsureCreated(ctx context.Context, obj interface{}) error
		EnsureDeleted(ctx context.Context, obj interface{}) error
	}

	BackOff backoff.BackOff
}

type basicResource struct {
	logger micrologger.Logger
	//resource framework.Resource
	resource interface {
		Name() string
		Underlying() framework.Resource

		EnsureCreated(ctx context.Context, obj interface{}) error
		EnsureDeleted(ctx context.Context, obj interface{}) error
	}

	backOff backoff.BackOff
}

func newBasicResource(config basicResourceConfig) (*basicResource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Resource must not be empty")
	}

	if config.BackOff == nil {
		config.BackOff = backoff.NewExponentialBackOff()
	}

	r := &basicResource{
		logger: config.Logger.With(
			"underlyingResource", config.Resource.Underlying().Name(),
		),
		resource: config.Resource,

		backOff: config.BackOff,
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
		r.logger.LogCtx(ctx, "warning", fmt.Sprintf("retrying 'EnsureCreated' due to error (%s)", err.Error()))
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
		r.logger.LogCtx(ctx, "warning", fmt.Sprintf("retrying 'EnsureDeleted' due to error (%s)", err.Error()))
	}

	err = backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *basicResource) Underlying() framework.Resource {
	return r.resource.Underlying()
}
