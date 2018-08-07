package retryresource

import (
	"context"
	"fmt"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/controller/resource/internal"
)

type resourceWrapper struct {
	logger   micrologger.Logger
	resource controller.Resource

	backOff backoff.Interface

	name string
}

func newResourceWrapper(config Config) (*resourceWrapper, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Resource must not be empty")
	}

	if config.BackOff == nil {
		config.BackOff = backoff.NewExponential(2*time.Minute, 10*time.Second)
	}

	var name string
	{
		u, err := internal.Underlying(config.Resource)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		name = u.Name()
	}

	r := &resourceWrapper{
		logger: config.Logger.With(
			"underlyingResource", name,
		),
		resource: config.Resource,

		backOff: config.BackOff,

		name: name,
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
		r.logger.LogCtx(ctx, "level", "warning", "message", "retrying due to error", "stack", fmt.Sprintf("%#v", err))
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
		r.logger.LogCtx(ctx, "level", "warning", "message", "retrying due to error", "stack", fmt.Sprintf("%#v", err))
	}

	err = backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *resourceWrapper) Name() string {
	return r.name
}

// Wrapped implements internal.Wrapper interface.
func (r *resourceWrapper) Wrapped() controller.Resource {
	return r.resource
}
