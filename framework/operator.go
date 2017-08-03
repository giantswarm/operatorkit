package framework

import (
	"fmt"
	"time"

	"github.com/cenk/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

// Config represents the configuration used to create a new operator framework.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger
}

// DefaultConfig provides a default configuration to create a new operator
// framework by best effort.
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
		Logger: newLogger,
	}
}

// New creates a new configured operator framework.
func New(config Config) (*Framework, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	newFramework := &Framework{
		// Dependencies.
		logger: config.Logger,
	}

	return newFramework, nil
}

type Framework struct {
	// Dependencies.
	logger micrologger.Logger
}

// ProcessCreate is a drop-in for an informer's AddFunc. It receives the custom
// object observed during TPR watches and anything that implements Resource.
// ProcessCreate takes care about all necessary reconciliation logic for create
// events.
//
//     func addFunc(obj interface{}) {
//         err := f.ProcessCreate(obj, resources)
//         if err != nil {
//             // error handling here
//         }
//     }
//
//     newResourceEventHandler := &cache.ResourceEventHandlerFuncs{
//         AddFunc:    addFunc,
//     }
//
func (f *Framework) ProcessCreate(obj interface{}, resources []Resource) error {
	if len(resources) == 0 {
		return microerror.Maskf(executionFailedError, "resources must not be empty")
	}

	for _, r := range resources {
		currentState, err := r.GetCurrentState(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		desiredState, err := r.GetDesiredState(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		createState, err := r.GetCreateState(obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}

		err = r.ProcessCreateState(obj, createState)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

// ProcessCreateWithBackoff is the same as ProcessCreate but takes an additional
// backoff implementation to retry the creation on error. Errors are logged
// eventually in case an appropriate logger is configured.
func (f *Framework) ProcessCreateWithBackoff(obj interface{}, resources []Resource, b backoff.BackOff) error {
	o := func() error {
		err := f.ProcessCreate(obj, resources)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		f.logger.Log("debug", fmt.Sprintf("retrying creation due to error (%s)", err.Error()))
	}

	err := backoff.RetryNotify(o, b, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

// ProcessDelete is a drop-in for an informer's DeleteFunc. It receives the
// custom object observed during TPR watches and anything that implements
// Resource. ProcessDelete takes care about all necessary reconciliation logic
// for delete events.
//
//     func deleteFunc(obj interface{}) {
//         err := f.ProcessDelete(obj, resources)
//         if err != nil {
//             // error handling here
//         }
//     }
//
//     newResourceEventHandler := &cache.ResourceEventHandlerFuncs{
//         DeleteFunc:    deleteFunc,
//     }
//
func (f *Framework) ProcessDelete(obj interface{}, resources []Resource) error {
	if len(resources) == 0 {
		return microerror.Maskf(executionFailedError, "resources must not be empty")
	}

	for _, r := range resources {
		currentState, err := r.GetCurrentState(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		desiredState, err := r.GetDesiredState(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		deleteState, err := r.GetDeleteState(obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}

		err = r.ProcessDeleteState(obj, deleteState)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

// ProcessDeleteWithBackoff is the same as ProcessDelete but takes an additional
// backoff implementation to retry the deletion on error. Errors are logged
// eventually in case an appropriate logger is configured.
func (f *Framework) ProcessDeleteWithBackoff(obj interface{}, resources []Resource, b backoff.BackOff) error {
	o := func() error {
		err := f.ProcessDelete(obj, resources)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	n := func(err error, dur time.Duration) {
		f.logger.Log("debug", fmt.Sprintf("retrying deletion due to error (%s)", err.Error()))
	}

	err := backoff.RetryNotify(o, b, n)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
