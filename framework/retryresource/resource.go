package retryresource

import (
	"fmt"
	"time"

	"github.com/cenk/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/framework"
)

const (
	// Name is the identifier of the resource.
	Name = "retry"
)

// Config represents the configuration used to create a new retry resource.
type Config struct {
	// Dependencies.
	BackOff  backoff.BackOff
	Logger   micrologger.Logger
	Resource framework.Resource
}

// DefaultConfig provides a default configuration to create a new retry resource
// by best effort.
func DefaultConfig() Config {
	var err error

	var newLogger micrologger.Logger
	{
		config := micrologger.DefaultConfig()
		newLogger, err = micrologger.New(config)
		if err != nil {
			panic(err)
		}
	}

	return Config{
		// Dependencies.
		BackOff:  backoff.NewExponentialBackOff(),
		Logger:   newLogger,
		Resource: nil,
	}
}

// New creates a new configured retry resource.
func New(config Config) (*Resource, error) {
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

	newResource := &Resource{
		// Dependencies.
		backOff: config.BackOff,
		logger: config.Logger.With(
			"underlyingResource", config.Resource.Underlying().Name(),
		),
		resource: config.Resource,
	}

	return newResource, nil
}

type Resource struct {
	// Dependencies.
	backOff  backoff.BackOff
	logger   micrologger.Logger
	resource framework.Resource
}

func (r *Resource) GetCurrentState(obj interface{}) (interface{}, error) {
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

func (r *Resource) GetDesiredState(obj interface{}) (interface{}, error) {
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

func (r *Resource) GetCreateState(obj, currentState, desiredState interface{}) (interface{}, error) {
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

func (r *Resource) GetDeleteState(obj, currentState, desiredState interface{}) (interface{}, error) {
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

func (r *Resource) GetUpdateState(obj, currentState, desiredState interface{}) (interface{}, interface{}, interface{}, error) {
	var err error

	var createState interface{}
	var deleteState interface{}
	var updateState interface{}

	o := func() error {
		createState, deleteState, updateState, err = r.resource.GetUpdateState(obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.Log("warning", fmt.Sprintf("retrying 'GetUpdateState' due to error (%s)", err.Error()))
	}

	err = backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return nil, nil, nil, microerror.Mask(err)
	}

	return createState, deleteState, updateState, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) ProcessCreateState(obj, createState interface{}) error {
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

func (r *Resource) ProcessDeleteState(obj, deleteState interface{}) error {
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

func (r *Resource) ProcessUpdateState(obj, updateState interface{}) error {
	o := func() error {
		err := r.resource.ProcessUpdateState(obj, updateState)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		r.logger.Log("warning", fmt.Sprintf("retrying 'ProcessUpdateState' due to error (%s)", err.Error()))
	}

	err := backoff.RetryNotify(o, r.backOff, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) Underlying() framework.Resource {
	return r.resource.Underlying()
}
