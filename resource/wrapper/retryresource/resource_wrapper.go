package retryresource

import (
	"context"
	"fmt"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/resource"
)

type resourceWrapper struct {
	backOff  backoff.Interface
	logger   micrologger.Logger
	resource resource.Interface
}

func newResourceWrapper(config Config) (*resourceWrapper, error) {
	if config.BackOff == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.BackOff must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
	}

	r := &resourceWrapper{
		backOff:  config.BackOff,
		logger:   config.Logger,
		resource: config.Resource,
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
	return r.resource.Name()
}

// Wrapped implements internal.Wrapper interface.
func (r *resourceWrapper) Wrapped() resource.Interface {
	return r.resource
}
