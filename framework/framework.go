package framework

import (
	"context"
	"fmt"
	"sync"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/tools/cache"

	"github.com/giantswarm/operatorkit/framework/canceledcontext"
)

// Config represents the configuration used to create a new operator framework.
type Config struct {
	// Dependencies.
	Initializer func(ctx context.Context, obj interface{}) (context.Context, error)
	Logger      micrologger.Logger
	Resources   []Resource
}

// DefaultConfig provides a default configuration to create a new operator
// framework by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Initializer: nil,
		Logger:      nil,
		Resources:   nil,
	}
}

type Framework struct {
	// Dependencies.
	initializer func(ctx context.Context, obj interface{}) (context.Context, error)
	logger      micrologger.Logger
	resources   []Resource

	// Internals.
	mutex sync.Mutex
}

// New creates a new configured operator framework.
func New(config Config) (*Framework, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if len(config.Resources) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "config.Resources must not be empty")
	}

	newFramework := &Framework{
		// Dependencies.
		initializer: config.Initializer,
		logger:      config.Logger,
		resources:   config.Resources,

		// Internals.
		mutex: sync.Mutex{},
	}

	return newFramework, nil
}

// AddFunc executes the framework's ProcessCreate and ProcessUpdate functions,
// in this order. This guarantees resource creation is always done before
// resource updates.
func (f *Framework) AddFunc(obj interface{}) {
	// We lock the AddFunc/DeleteFunc to make sure only one AddFunc/DeleteFunc is
	// executed at a time. AddFunc/DeleteFunc is not thread safe. This is
	// important because the source of truth for an operator are the reconciled
	// resources. In case we would run the operator logic in parallel, we would
	// run into race conditions.
	f.mutex.Lock()
	defer f.mutex.Unlock()

	ctx := context.Background()
	ctx = canceledcontext.NewContext(ctx, make(chan struct{}))

	if f.initializer != nil {
		var err error
		ctx, err = f.initializer(ctx, obj)
		if err != nil {
			f.logger.Log("error", fmt.Sprintf("%#v", err), "event", "create")
			return
		}
	}

	f.logger.Log("action", "start", "component", "operatorkit", "function", "ProcessCreate")

	err := ProcessCreate(ctx, obj, f.resources)
	if err != nil {
		f.logger.Log("error", fmt.Sprintf("%#v", err), "event", "create")
		return
	}

	f.logger.Log("action", "end", "component", "operatorkit", "function", "ProcessCreate")

	f.logger.Log("action", "start", "component", "operatorkit", "function", "ProcessUpdate")

	err = ProcessUpdate(ctx, obj, f.resources)
	if err != nil {
		f.logger.Log("error", fmt.Sprintf("%#v", err), "event", "update")
		return
	}

	f.logger.Log("action", "end", "component", "operatorkit", "function", "ProcessUpdate")
}

// DeleteFunc executes the framework's ProcessDelete function.
func (f *Framework) DeleteFunc(obj interface{}) {
	// We lock the AddFunc/DeleteFunc to make sure only one AddFunc/DeleteFunc is
	// executed at a time. AddFunc/DeleteFunc is not thread safe. This is
	// important because the source of truth for an operator are the reconciled
	// resources. In case we would run the operator logic in parallel, we would
	// run into race conditions.
	f.mutex.Lock()
	defer f.mutex.Unlock()

	ctx := context.Background()
	ctx = canceledcontext.NewContext(ctx, make(chan struct{}))

	if f.initializer != nil {
		var err error
		ctx, err = f.initializer(ctx, obj)
		if err != nil {
			f.logger.Log("error", fmt.Sprintf("%#v", err), "event", "delete")
			return
		}
	}

	f.logger.Log("action", "start", "component", "operatorkit", "function", "ProcessDelete")

	err := ProcessDelete(ctx, obj, f.resources)
	if err != nil {
		f.logger.Log("error", fmt.Sprintf("%#v", err), "event", "delete")
		return
	}

	f.logger.Log("action", "end", "component", "operatorkit", "function", "ProcessDelete")
}

// NewCacheResourceEventHandler returns the framework's event handler for the
// k8s client's cache informer implementation. The event handler has functions
// registered for the k8s client's add, delete and update events.
func (f *Framework) NewCacheResourceEventHandler() *cache.ResourceEventHandlerFuncs {
	newHandler := &cache.ResourceEventHandlerFuncs{
		AddFunc:    f.AddFunc,
		DeleteFunc: f.DeleteFunc,
		UpdateFunc: f.UpdateFunc,
	}

	return newHandler
}

// UpdateFunc only redirects to AddFunc and only dispatches the new custom
// object received.
func (f *Framework) UpdateFunc(oldObj, newObj interface{}) {
	f.AddFunc(newObj)
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
func ProcessCreate(ctx context.Context, obj interface{}, resources []Resource) error {
	if len(resources) == 0 {
		return microerror.Maskf(executionFailedError, "resources must not be empty")
	}

	for _, r := range resources {
		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
		currentState, err := r.GetCurrentState(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
		desiredState, err := r.GetDesiredState(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
		createState, err := r.GetCreateState(ctx, obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}

		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
		err = r.ProcessCreateState(ctx, obj, createState)
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
func ProcessDelete(ctx context.Context, obj interface{}, resources []Resource) error {
	if len(resources) == 0 {
		return microerror.Maskf(executionFailedError, "resources must not be empty")
	}

	for _, r := range resources {
		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
		currentState, err := r.GetCurrentState(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
		desiredState, err := r.GetDesiredState(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
		deleteState, err := r.GetDeleteState(ctx, obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}

		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
		err = r.ProcessDeleteState(ctx, obj, deleteState)
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
func ProcessUpdate(ctx context.Context, obj interface{}, resources []Resource) error {
	if len(resources) == 0 {
		return microerror.Maskf(executionFailedError, "resources must not be empty")
	}

	for _, r := range resources {
		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
		currentState, err := r.GetCurrentState(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
		desiredState, err := r.GetDesiredState(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
		createState, deleteState, updateState, err := r.GetUpdateState(ctx, obj, currentState, desiredState)
		if err != nil {
			return microerror.Mask(err)
		}

		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
		err = r.ProcessCreateState(ctx, obj, createState)
		if err != nil {
			return microerror.Mask(err)
		}

		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
		err = r.ProcessDeleteState(ctx, obj, deleteState)
		if err != nil {
			return microerror.Mask(err)
		}

		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
		err = r.ProcessUpdateState(ctx, obj, updateState)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
