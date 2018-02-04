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

const (
	// Name is the identifier of the resource.
	Name = "retry"
)

type crudResourceConfig struct {
	Logger   micrologger.Logger
	Resource framework.CRUDResourceOps

	BackOff backoff.BackOff
}

type crudResource struct {
	logger   micrologger.Logger
	resource framework.CRUDResourceOps

	backOff backoff.BackOff
}

func newCRUDResource(config crudResourceConfig) (*crudResource, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Resource must not be empty")
	}

	if config.BackOff == nil {
		config.BackOff = backoff.NewExponentialBackOff()
	}

	r := &crudResource{
		logger: config.Logger.With(
			"underlyingResource", config.Resource.Underlying().Name(),
		),
		resource: config.Resource,

		backOff: config.BackOff,
	}

	return r, nil
}

func (r *crudResource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	var err error

	var v interface{}
	o := func() error {
		v, err = r.resource.GetCurrentState(ctx, obj)
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

func (r *crudResource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	var err error

	var v interface{}
	o := func() error {
		v, err = r.resource.GetDesiredState(ctx, obj)
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

func (r *crudResource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	var err error

	var v *framework.Patch
	o := func() error {
		v, err = r.resource.NewUpdatePatch(ctx, obj, currentState, desiredState)
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

func (r *crudResource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	var err error

	var v *framework.Patch
	o := func() error {
		v, err = r.resource.NewDeletePatch(ctx, obj, currentState, desiredState)
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

func (r *crudResource) Name() string {
	return Name
}

func (r *crudResource) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	o := func() error {
		err := r.resource.ApplyCreateChange(ctx, obj, createState)
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

func (r *crudResource) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	o := func() error {
		err := r.resource.ApplyDeleteChange(ctx, obj, deleteState)
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

func (r *crudResource) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	o := func() error {
		err := r.resource.ApplyUpdateChange(ctx, obj, updateState)
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

func (r *crudResource) Underlying() framework.Resource {
	return r.resource.Underlying()
}
