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

type crudResourceOpsConfig struct {
	Logger micrologger.Logger
	Ops    framework.CRUDResourceOps

	BackOff backoff.BackOff
}

type crudResourceOps struct {
	logger     micrologger.Logger
	underlying framework.CRUDResourceOps

	backOff backoff.BackOff
}

func newCRUDResourceOps(config crudResourceOpsConfig) (*crudResourceOps, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Ops == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Ops must not be empty")
	}

	if config.BackOff == nil {
		config.BackOff = backoff.NewExponentialBackOff()
	}

	r := &crudResourceOps{
		logger: config.Logger.With(
			"underlyingResource", config.Ops.Underlying().Name(),
		),
		underlying: config.Ops,

		backOff: config.BackOff,
	}

	return r, nil
}

func (r *crudResourceOps) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	var err error

	var v interface{}
	o := func() error {
		v, err = r.underlying.GetCurrentState(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.LogCtx(ctx, "warning", fmt.Sprintf("retrying 'GetCurrentState' due to error (%s)", err.Error()))
	}

	err = backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *crudResourceOps) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	var err error

	var v interface{}
	o := func() error {
		v, err = r.underlying.GetDesiredState(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.LogCtx(ctx, "warning", fmt.Sprintf("retrying 'GetDesiredState' due to error (%s)", err.Error()))
	}

	err = backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *crudResourceOps) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	var err error

	var v *framework.Patch
	o := func() error {
		v, err = r.underlying.NewUpdatePatch(ctx, obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.LogCtx(ctx, "warning", fmt.Sprintf("retrying 'NewUpdatePatch' due to error (%s)", err.Error()))
	}

	err = backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *crudResourceOps) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	var err error

	var v *framework.Patch
	o := func() error {
		v, err = r.underlying.NewDeletePatch(ctx, obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.LogCtx(ctx, "warning", fmt.Sprintf("retrying 'NewDeletePatch' due to error (%s)", err.Error()))
	}

	err = backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *crudResourceOps) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	o := func() error {
		err := r.underlying.ApplyCreateChange(ctx, obj, createState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.LogCtx(ctx, "warning", fmt.Sprintf("retrying 'ApplyCreatePatch' due to error (%s)", err.Error()))
	}

	err := backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *crudResourceOps) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	o := func() error {
		err := r.underlying.ApplyDeleteChange(ctx, obj, deleteState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.LogCtx(ctx, "warning", fmt.Sprintf("retrying 'ApplyDeletePatch' due to error (%s)", err.Error()))
	}

	err := backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *crudResourceOps) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	o := func() error {
		err := r.underlying.ApplyUpdateChange(ctx, obj, updateState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.LogCtx(ctx, "warning", fmt.Sprintf("retrying 'ApplyUpdatePatch' due to error (%s)", err.Error()))
	}

	err := backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
