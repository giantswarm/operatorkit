package retryresource

import (
	"context"
	"fmt"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/controller"
)

type crudResourceOpsWrapperConfig struct {
	Logger micrologger.Logger
	Ops    controller.CRUDResourceOps

	BackOff backoff.Interface
}

type crudResourceWrapperOps struct {
	logger     micrologger.Logger
	underlying controller.CRUDResourceOps

	backOff backoff.Interface
}

func newCRUDResourceWrapperOps(config crudResourceOpsWrapperConfig) (*crudResourceWrapperOps, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Ops == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Ops must not be empty")
	}

	if config.BackOff == nil {
		config.BackOff = backoff.NewExponential(2*time.Minute, 10*time.Second)
	}

	o := &crudResourceWrapperOps{
		logger: config.Logger.With(
			"underlyingResource", config.Ops.Name(),
		),
		underlying: config.Ops,

		backOff: config.BackOff,
	}

	return o, nil
}

func (o *crudResourceWrapperOps) Name() string {
	return o.underlying.Name()
}

func (o *crudResourceWrapperOps) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	var err error

	var v interface{}
	op := func() error {
		v, err = o.underlying.GetCurrentState(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		o.logger.LogCtx(ctx, "level", "warning", "message", "retrying due to error", "stack", fmt.Sprintf("%#v", err))
	}

	err = backoff.RetryNotify(op, o.backOff, n)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (o *crudResourceWrapperOps) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	var err error

	var v interface{}
	op := func() error {
		v, err = o.underlying.GetDesiredState(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		o.logger.LogCtx(ctx, "level", "warning", "message", "retrying due to error", "stack", fmt.Sprintf("%#v", err))

	}

	err = backoff.RetryNotify(op, o.backOff, n)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (o *crudResourceWrapperOps) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*controller.Patch, error) {
	var err error

	var v *controller.Patch
	op := func() error {
		v, err = o.underlying.NewUpdatePatch(ctx, obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		o.logger.LogCtx(ctx, "level", "warning", "message", "retrying due to error", "stack", fmt.Sprintf("%#v", err))
	}

	err = backoff.RetryNotify(op, o.backOff, n)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (o *crudResourceWrapperOps) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*controller.Patch, error) {
	var err error

	var v *controller.Patch
	op := func() error {
		v, err = o.underlying.NewDeletePatch(ctx, obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		o.logger.LogCtx(ctx, "level", "warning", "message", "retrying due to error", "stack", fmt.Sprintf("%#v", err))
	}

	err = backoff.RetryNotify(op, o.backOff, n)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (o *crudResourceWrapperOps) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	op := func() error {
		err := o.underlying.ApplyCreateChange(ctx, obj, createState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		o.logger.LogCtx(ctx, "level", "warning", "message", "retrying due to error", "stack", fmt.Sprintf("%#v", err))
	}

	err := backoff.RetryNotify(op, o.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (o *crudResourceWrapperOps) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	op := func() error {
		err := o.underlying.ApplyDeleteChange(ctx, obj, deleteState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		o.logger.LogCtx(ctx, "level", "warning", "message", "retrying due to error", "stack", fmt.Sprintf("%#v", err))
	}

	err := backoff.RetryNotify(op, o.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (o *crudResourceWrapperOps) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	op := func() error {
		err := o.underlying.ApplyUpdateChange(ctx, obj, updateState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		o.logger.LogCtx(ctx, "level", "warning", "message", "retrying due to error", "stack", fmt.Sprintf("%#v", err))
	}

	err := backoff.RetryNotify(op, o.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
