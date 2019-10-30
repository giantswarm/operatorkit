package retryresource

import (
	"context"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/resource/crud"
)

type crudResourceConfig struct {
	BackOff backoff.Interface
	CRUD    crud.Interface
	Logger  micrologger.Logger
}

type crudResource struct {
	backOff backoff.Interface
	crud    crud.Interface
	logger  micrologger.Logger
}

func newCRUDResource(config crudResourceConfig) (*crudResource, error) {
	if config.BackOff == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.BackOff must not be empty", config)
	}
	if config.CRUD == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.CRUD must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	r := &crudResource{
		backOff: config.BackOff,
		crud:    config.CRUD,
		logger:  config.Logger,
	}

	return r, nil
}

func (r *crudResource) Name() string {
	return r.crud.Name()
}

func (r *crudResource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	var err error

	var v interface{}
	o := func() error {
		v, err = r.crud.GetCurrentState(ctx, obj)
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
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *crudResource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	var err error

	var v interface{}
	o := func() error {
		v, err = r.crud.GetDesiredState(ctx, obj)
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
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *crudResource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*crud.Patch, error) {
	var err error

	var v *crud.Patch
	o := func() error {
		v, err = r.crud.NewUpdatePatch(ctx, obj, currentState, desiredState)
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
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *crudResource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*crud.Patch, error) {
	var err error

	var v *crud.Patch
	o := func() error {
		v, err = r.crud.NewDeletePatch(ctx, obj, currentState, desiredState)
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
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *crudResource) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	o := func() error {
		err := r.crud.ApplyCreateChange(ctx, obj, createState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.LogCtx(ctx, "level", "warning", "message", "retrying due to error", "stack", microerror.Stack(err))
	}

	err := backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *crudResource) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	o := func() error {
		err := r.crud.ApplyDeleteChange(ctx, obj, deleteState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.LogCtx(ctx, "level", "warning", "message", "retrying due to error", "stack", microerror.Stack(err))
	}

	err := backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *crudResource) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	o := func() error {
		err := r.crud.ApplyUpdateChange(ctx, obj, updateState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.LogCtx(ctx, "level", "warning", "message", "retrying due to error", "stack", microerror.Stack(err))
	}

	err := backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
