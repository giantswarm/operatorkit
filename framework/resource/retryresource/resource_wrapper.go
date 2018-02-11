package retryresource

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/framework/resource/internal/framework"
)

type resourceWrapperConfig struct {
	Logger micrologger.Logger
	// TODO make Resource framework.Resource
	Resource interface {
		Name() string
		Underlying() framework.Resource

		EnsureCreated(ctx context.Context, obj interface{}) error
		EnsureDeleted(ctx context.Context, obj interface{}) error
	}

	BackOff backoff.BackOff
}

type resourceWrapper struct {
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

func newResourceWrapper(config resourceWrapperConfig) (*resourceWrapper, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Resource must not be empty")
	}

	if config.BackOff == nil {
		config.BackOff = backoff.NewExponentialBackOff()
	}

	r := &resourceWrapper{
		logger: config.Logger.With(
		// TODO Uncomment when Resource interface is updated.
		//"underlyingResource", internal.Underlying(config.Resource).Name(),
		),
		resource: config.Resource,

		backOff: config.BackOff,
	}

	return r, nil
}

func (r *resourceWrapper) EnsureCreated(ctx context.Context, obj interface{}) error {
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

func (r *resourceWrapper) EnsureDeleted(ctx context.Context, obj interface{}) error {
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

func (r *resourceWrapper) Name() string {
	return r.resource.Name()
}

// Wrapped implements internal.Wrapper interface.
func (r *resourceWrapper) Wrapped() framework.Resource {
	return r.resource
}
