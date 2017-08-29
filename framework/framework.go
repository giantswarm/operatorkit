package framework

import (
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

// ProcessUpdate is a drop-in for an informer's UpdateFunc. It receives the new
// custom object observed during TPR watches and anything that implements
// Resource. ProcessUpdate takes care about all necessary reconciliation logic
// for update events. For complex resources this means state has to be created,
// deleted and updated eventually, in this order.
//
//     func updateFunc(oldObj, newObj interface{}) {
//         err := f.ProcessUpdate(newObj, resources)
//         if err != nil {
//             // error handling here
//         }
//     }
//
//     newResourceEventHandler := &cache.ResourceEventHandlerFuncs{
//         UpdateFunc:    updateFunc,
//     }
//
func (f *Framework) ProcessUpdate(newObj interface{}, resources []Resource) error {
	if len(resources) == 0 {
		return microerror.Maskf(executionFailedError, "resources must not be empty")
	}

	for _, r := range resources {
		currentState, err := r.GetCurrentState(newObj)
		if err != nil {
			return microerror.Mask(err)
		}

		desiredState, err := r.GetDesiredState(newObj)
		if err != nil {
			return microerror.Mask(err)
		}

		createState, deleteState, updateState, err := r.GetUpdateState(newObj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}

		err = r.ProcessCreateState(newObj, createState)
		if err != nil {
			return microerror.Mask(err)
		}

		err = r.ProcessDeleteState(newObj, deleteState)
		if err != nil {
			return microerror.Mask(err)
		}

		err = r.ProcessUpdateState(newObj, updateState)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
