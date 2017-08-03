package framework

import (
	"fmt"
	"time"

	"github.com/cenk/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

// RetryResourceConfig represents the configuration used to create a new retry
// resource.
type RetryResourceConfig struct {
	// Dependencies.
	BackOff  backoff.BackOff
	Logger   micrologger.Logger
	Resource Resource
}

// DefaultRetryResourceConfig provides a default configuration to create a new
// retry resource by best effort.
func DefaultRetryResourceConfig() RetryResourceConfig {
	var err error

	var newLogger micrologger.Logger
	{
		config := micrologger.DefaultConfig()
		newLogger, err = micrologger.New(config)
		if err != nil {
			panic(err)
		}
	}

	return RetryResourceConfig{
		// Dependencies.
		BackOff:  backoff.NewExponentialBackOff(),
		Logger:   newLogger,
		Resource: nil,
	}
}

// NewRetryResource creates a new configured retry resource.
func NewRetryResource(config RetryResourceConfig) (*RetryResource, error) {
	// Dependencies.
	if config.BackOff == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.BackOff must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Resource must not be empty")
	}

	newRetryResource := &RetryResource{
		// Dependencies.
		backOff:  config.BackOff,
		logger:   config.Logger,
		resource: config.Resource,
	}

	return newRetryResource, nil
}

type RetryResource struct {
	// Dependencies.
	backOff  backoff.BackOff
	logger   micrologger.Logger
	resource Resource
}

func (r *RetryResource) GetCurrentState(obj interface{}) (interface{}, error) {
	var err error

	var v interface{}
	o := func() error {
		v, err = r.resource.GetCurrentState(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.Log("warning", fmt.Sprintf("retrying 'GetCurrentState' due to error (%s)", err.Error()))
	}

	err = backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *RetryResource) GetDesiredState(obj interface{}) (interface{}, error) {
	var err error

	var v interface{}
	o := func() error {
		v, err = r.resource.GetDesiredState(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.Log("warning", fmt.Sprintf("retrying 'GetDesiredState' due to error (%s)", err.Error()))
	}

	err = backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *RetryResource) GetCreateState(obj, currentState, desiredState interface{}) (interface{}, error) {
	var err error

	var v interface{}
	o := func() error {
		v, err = r.resource.GetCreateState(obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.Log("warning", fmt.Sprintf("retrying 'GetCreateState' due to error (%s)", err.Error()))
	}

	err = backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *RetryResource) GetDeleteState(obj, currentState, desiredState interface{}) (interface{}, error) {
	var err error

	var v interface{}
	o := func() error {
		v, err = r.resource.GetDeleteState(obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.Log("warning", fmt.Sprintf("retrying 'GetDeleteState' due to error (%s)", err.Error()))
	}

	err = backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *RetryResource) ProcessCreateState(obj, createState interface{}) error {
	o := func() error {
		err := r.resource.ProcessCreateState(obj, createState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.Log("warning", fmt.Sprintf("retrying 'ProcessCreateState' due to error (%s)", err.Error()))
	}

	err := backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *RetryResource) ProcessDeleteState(obj, deleteState interface{}) error {
	o := func() error {
		err := r.resource.ProcessDeleteState(obj, deleteState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.Log("warning", fmt.Sprintf("retrying 'ProcessDeleteState' due to error (%s)", err.Error()))
	}

	err := backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
