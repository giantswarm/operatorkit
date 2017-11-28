package framework

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/micrologger/loggermeta"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/giantswarm/operatorkit/client/k8scrdclient"
	"github.com/giantswarm/operatorkit/crd"
	"github.com/giantswarm/operatorkit/framework/context/canceledcontext"
	"github.com/giantswarm/operatorkit/informer"
	"github.com/giantswarm/operatorkit/tpr"
)

// Config represents the configuration used to create a new operator framework.
type Config struct {
	// Dependencies.

	CRD       *crd.CRD
	CRDClient *k8scrdclient.CRDClient
	Informer  informer.Interface
	// InitCtxFunc is to prepare the given context for a single reconciliation
	// loop. Operators can implement common context packages to enable
	// communication between resources. These context packages can be set up
	// within the context initializer function. InitCtxFunc receives the custom
	// object being reconciled as second argument. Information provided by the
	// custom object can be used to initialize the context.
	InitCtxFunc func(ctx context.Context, obj interface{}) (context.Context, error)
	Logger      micrologger.Logger
	// ResourceRouter is to decide which resources to execute. Each custom object
	// being reconciled is executed against a list of resources. Since custom
	// objects may differ in version and/or structure the resource router enables
	// custom inspection before each reconciliation loop. That way whole resources
	// can be versioned and different resources can be executed depending on the
	// custom object being reconciled.
	ResourceRouter func(ctx context.Context, obj interface{}) ([]Resource, error)
	// TPR can be provided to ensure a third party resource exists. When given the
	// boot process of the framework ensures the TPR is created before starting
	// the conciliation.
	//
	// NOTE this is deprecated since the CRD concept is the successor of the TPR.
	TPR tpr.Interface

	// Settings.
	BackOffFactory func() backoff.BackOff
}

// DefaultConfig provides a default configuration to create a new operator
// framework by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		CRD:            nil,
		CRDClient:      nil,
		Informer:       nil,
		Logger:         nil,
		ResourceRouter: nil,
		TPR:            nil,

		// Settings.
		BackOffFactory: func() backoff.BackOff {
			b := backoff.NewExponentialBackOff()
			b.MaxElapsedTime = 0
			return backoff.WithMaxTries(b, 7)
		},
		InitCtxFunc: func(ctx context.Context, obj interface{}) (context.Context, error) {
			return ctx, nil
		},
	}
}

type Framework struct {
	// Dependencies.
	crd            *crd.CRD
	crdClient      *k8scrdclient.CRDClient
	informer       informer.Interface
	logger         micrologger.Logger
	resourceRouter func(ctx context.Context, obj interface{}) ([]Resource, error)
	tpr            tpr.Interface

	// Settings.
	backOffFactory func() backoff.BackOff
	initCtxFunc    func(ctx context.Context, obj interface{}) (context.Context, error)

	// Internals.
	bootOnce sync.Once
	mutex    sync.Mutex
}

// New creates a new configured operator framework.
func New(config Config) (*Framework, error) {
	// Dependencies.
	if config.CRD != nil && config.CRDClient == nil || config.CRD == nil && config.CRDClient != nil {
		return nil, microerror.Maskf(invalidConfigError, "config.CRD and config.CRDClient must not be empty when either given")
	}
	if config.Informer == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Informer must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.ResourceRouter == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.ResourceRouter must not be empty")
	}

	// Settings.
	if config.BackOffFactory == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.BackOffFactory must not be empty")
	}
	if config.InitCtxFunc == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.InitCtxFunc must not be empty")
	}

	initCtxFunc := func(ctx context.Context, obj interface{}) (context.Context, error) {
		ctx = canceledcontext.NewContext(ctx, make(chan struct{}))

		ctx, err := config.InitCtxFunc(ctx, obj)
		if err != nil {
			return nil, microerror.Maskf(err, "initializing context")
		}

		accessor, err := meta.Accessor(obj)
		if err != nil {
			config.Logger.Log("warning", fmt.Sprintf("cannot create accessor for object %#v", obj))
		} else {
			meta, ok := loggermeta.FromContext(ctx)
			if !ok {
				meta = loggermeta.New()
			}
			meta.KeyVals["object"] = accessor.GetSelfLink()

			ctx = loggermeta.NewContext(ctx, meta)
		}

		return ctx, nil
	}

	f := &Framework{
		// Dependencies.
		crd:            config.CRD,
		crdClient:      config.CRDClient,
		informer:       config.Informer,
		logger:         config.Logger,
		resourceRouter: config.ResourceRouter,
		tpr:            config.TPR,

		// Settings.
		backOffFactory: config.BackOffFactory,
		initCtxFunc:    initCtxFunc,

		// Internals.
		bootOnce: sync.Once{},
		mutex:    sync.Mutex{},
	}

	return f, nil
}

// AddFunc executes the framework's ProcessCreate function.
func (f *Framework) AddFunc(obj interface{}) {
	// AddFunc/DeleteFunc/UpdateFunc is synchronized to make sure only one
	// of them is executed at a time. AddFunc/DeleteFunc/UpdateFunc is not
	// thread safe. This is important because the source of truth for an
	// operator are the reconciled resources. In case we would run the
	// operator logic in parallel, we would run into race conditions.
	f.mutex.Lock()
	defer f.mutex.Unlock()

	ctx := context.Background()
	ctx, err := f.initCtxFunc(ctx, obj)
	if err != nil {
		f.logger.LogCtx(ctx, "error", fmt.Sprintf("%#v", err), "event", "create")
		return
	}

	rs, err := f.resourceRouter(ctx, obj)
	if err != nil {
		f.logger.LogCtx(ctx, "error", fmt.Sprintf("%#v", err), "event", "create")
		return
	}

	f.logger.LogCtx(ctx, "action", "start", "component", "operatorkit", "function", "ProcessCreate")

	err = ProcessCreate(ctx, obj, rs)
	if err != nil {
		f.logger.LogCtx(ctx, "error", fmt.Sprintf("%#v", err), "event", "create")
		return
	}

	f.logger.LogCtx(ctx, "action", "end", "component", "operatorkit", "function", "ProcessCreate")
}

func (f *Framework) Boot() {
	ctx := context.TODO()

	f.bootOnce.Do(func() {
		operation := func() error {
			err := f.bootWithError(ctx)
			if err != nil {
				return microerror.Mask(err)
			}

			return nil
		}

		notifier := func(err error, d time.Duration) {
			f.logger.LogCtx(ctx, "warning", fmt.Sprintf("retrying operator boot due to error: %#v", microerror.Mask(err)))
		}

		err := backoff.RetryNotify(operation, f.backOffFactory(), notifier)
		if err != nil {
			f.logger.LogCtx(ctx, "error", fmt.Sprintf("stop operator boot retries due to too many errors: %#v", microerror.Mask(err)))
			os.Exit(1)
		}
	})
}

// DeleteFunc executes the framework's ProcessDelete function.
func (f *Framework) DeleteFunc(obj interface{}) {
	// AddFunc/DeleteFunc/UpdateFunc is synchronized to make sure only one
	// of them is executed at a time. AddFunc/DeleteFunc/UpdateFunc is not
	// thread safe. This is important because the source of truth for an
	// operator are the reconciled resources. In case we would run the
	// operator logic in parallel, we would run into race conditions.
	f.mutex.Lock()
	defer f.mutex.Unlock()

	ctx := context.Background()
	ctx, err := f.initCtxFunc(ctx, obj)
	if err != nil {
		f.logger.LogCtx(ctx, "error", fmt.Sprintf("%#v", err), "event", "delete")
		return
	}

	rs, err := f.resourceRouter(ctx, obj)
	if err != nil {
		f.logger.LogCtx(ctx, "error", fmt.Sprintf("%#v", err), "event", "delete")
		return
	}

	f.logger.LogCtx(ctx, "action", "start", "component", "operatorkit", "function", "ProcessDelete")

	err = ProcessDelete(ctx, obj, rs)
	if err != nil {
		f.logger.LogCtx(ctx, "error", fmt.Sprintf("%#v", err), "event", "delete")
		return
	}

	f.logger.LogCtx(ctx, "action", "end", "component", "operatorkit", "function", "ProcessDelete")
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

// UpdateFunc executes the framework's ProcessUpdate function.
func (f *Framework) UpdateFunc(oldObj, newObj interface{}) {
	obj := newObj

	// AddFunc/DeleteFunc/UpdateFunc is synchronized to make sure only one
	// of them is executed at a time. AddFunc/DeleteFunc/UpdateFunc is not
	// thread safe. This is important because the source of truth for an
	// operator are the reconciled resources. In case we would run the
	// operator logic in parallel, we would run into race conditions.
	f.mutex.Lock()
	defer f.mutex.Unlock()

	ctx := context.Background()
	ctx, err := f.initCtxFunc(ctx, obj)
	if err != nil {
		f.logger.LogCtx(ctx, "error", fmt.Sprintf("%#v", err), "event", "update")
		return
	}

	rs, err := f.resourceRouter(ctx, obj)
	if err != nil {
		f.logger.LogCtx(ctx, "error", fmt.Sprintf("%#v", err), "event", "update")
		return
	}

	f.logger.LogCtx(ctx, "action", "start", "component", "operatorkit", "function", "ProcessUpdate")

	err = ProcessUpdate(ctx, obj, rs)
	if err != nil {
		f.logger.LogCtx(ctx, "error", fmt.Sprintf("%#v", err), "event", "update")
		return
	}

	f.logger.LogCtx(ctx, "action", "end", "component", "operatorkit", "function", "ProcessUpdate")
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
	err := ProcessUpdate(ctx, obj, resources)
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
func ProcessDelete(ctx context.Context, obj interface{}, resources []Resource) error {
	if len(resources) == 0 {
		return microerror.Maskf(executionFailedError, "resources must not be empty")
	}

	for _, r := range resources {
		var err error

		var currentState interface{}
		{
			if canceledcontext.IsCanceled(ctx) {
				return nil
			}

			meta, ok := loggermeta.FromContext(ctx)
			if ok {
				meta.KeyVals["function"] = "GetCurrentState"
				defer delete(meta.KeyVals, "function")
			}
			currentState, err = r.GetCurrentState(ctx, obj)
			if err != nil {
				return microerror.Mask(err)
			}
		}

		var desiredState interface{}
		{
			if canceledcontext.IsCanceled(ctx) {
				return nil
			}

			meta, ok := loggermeta.FromContext(ctx)
			if ok {
				meta.KeyVals["function"] = "GetDesiredState"
				defer delete(meta.KeyVals, "function")
			}
			desiredState, err = r.GetDesiredState(ctx, obj)
			if err != nil {
				return microerror.Mask(err)
			}
		}

		var patch *Patch
		{
			if canceledcontext.IsCanceled(ctx) {
				return nil
			}

			meta, ok := loggermeta.FromContext(ctx)
			if ok {
				meta.KeyVals["function"] = "NewDeletePatch"
				defer delete(meta.KeyVals, "function")
			}
			patch, err = r.NewDeletePatch(ctx, obj, currentState, desiredState)
			if err != nil {
				return microerror.Mask(err)
			}

			if patch == nil {
				return microerror.Maskf(executionFailedError, "patch must not be nil")
			}
		}

		{
			if canceledcontext.IsCanceled(ctx) {
				return nil
			}

			createChange, ok := patch.getCreateChange()
			if ok {
				meta, ok := loggermeta.FromContext(ctx)
				if ok {
					meta.KeyVals["function"] = "ApplyCreateChange"
					defer delete(meta.KeyVals, "function")
				}
				err := r.ApplyCreateChange(ctx, obj, createChange)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}

		{
			if canceledcontext.IsCanceled(ctx) {
				return nil
			}

			deleteChange, ok := patch.getDeleteChange()
			if ok {
				meta, ok := loggermeta.FromContext(ctx)
				if ok {
					meta.KeyVals["function"] = "ApplyDeleteChange"
					defer delete(meta.KeyVals, "function")
				}
				err := r.ApplyDeleteChange(ctx, obj, deleteChange)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}

		{
			if canceledcontext.IsCanceled(ctx) {
				return nil
			}

			updateChange, ok := patch.getUpdateChange()
			if ok {
				meta, ok := loggermeta.FromContext(ctx)
				if ok {
					meta.KeyVals["function"] = "ApplyUpdateChange"
					defer delete(meta.KeyVals, "function")
				}
				err := r.ApplyUpdateChange(ctx, obj, updateChange)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}
	}

	return nil
}

// ProcessEvents takes the event channels created by the operatorkit informer
// and executes the framework's event functions accordingly.
func (f *Framework) ProcessEvents(ctx context.Context, deleteChan chan watch.Event, updateChan chan watch.Event, errChan chan error) {
	operation := func() error {
		for {
			select {
			case e := <-deleteChan:
				t := prometheus.NewTimer(frameworkHistogram.WithLabelValues("delete"))
				f.DeleteFunc(e.Object)
				t.ObserveDuration()
			case e := <-updateChan:
				t := prometheus.NewTimer(frameworkHistogram.WithLabelValues("update"))
				f.UpdateFunc(nil, e.Object)
				t.ObserveDuration()
			case err := <-errChan:
				return microerror.Mask(err)
			case <-ctx.Done():
				return nil
			}
		}
	}

	notifier := func(err error, d time.Duration) {
		f.logger.LogCtx(ctx, "warning", fmt.Sprintf("retrying operator event processing due to error: %#v", microerror.Mask(err)))
	}

	err := backoff.RetryNotify(operation, f.backOffFactory(), notifier)
	if err != nil {
		f.logger.LogCtx(ctx, "error", fmt.Sprintf("stop operator event processing retries due to too many errors: %#v", microerror.Mask(err)))
		os.Exit(1)
	}
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
		var err error

		var currentState interface{}
		{
			if canceledcontext.IsCanceled(ctx) {
				return nil
			}

			meta, ok := loggermeta.FromContext(ctx)
			if ok {
				meta.KeyVals["function"] = "GetCurrentState"
				defer delete(meta.KeyVals, "function")
			}
			currentState, err = r.GetCurrentState(ctx, obj)
			if err != nil {
				return microerror.Mask(err)
			}
		}

		var desiredState interface{}
		{
			if canceledcontext.IsCanceled(ctx) {
				return nil
			}

			meta, ok := loggermeta.FromContext(ctx)
			if ok {
				meta.KeyVals["function"] = "GetDesiredState"
				defer delete(meta.KeyVals, "function")
			}
			desiredState, err = r.GetDesiredState(ctx, obj)
			if err != nil {
				return microerror.Mask(err)
			}
		}

		var patch *Patch
		{
			if canceledcontext.IsCanceled(ctx) {
				return nil
			}

			meta, ok := loggermeta.FromContext(ctx)
			if ok {
				meta.KeyVals["function"] = "NewUpdatePatch"
				defer delete(meta.KeyVals, "function")
			}
			patch, err = r.NewUpdatePatch(ctx, obj, currentState, desiredState)
			if err != nil {
				return microerror.Mask(err)
			}

			if patch == nil {
				return microerror.Maskf(executionFailedError, "patch must not be nil")
			}
		}

		{
			if canceledcontext.IsCanceled(ctx) {
				return nil
			}

			createState, ok := patch.getCreateChange()
			if ok {
				meta, ok := loggermeta.FromContext(ctx)
				if ok {
					meta.KeyVals["function"] = "ApplyCreateChange"
					defer delete(meta.KeyVals, "function")
				}
				err := r.ApplyCreateChange(ctx, obj, createState)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}

		{
			if canceledcontext.IsCanceled(ctx) {
				return nil
			}

			deleteState, ok := patch.getDeleteChange()
			if ok {
				meta, ok := loggermeta.FromContext(ctx)
				if ok {
					meta.KeyVals["function"] = "ApplyDeleteChange"
					defer delete(meta.KeyVals, "function")
				}
				err := r.ApplyDeleteChange(ctx, obj, deleteState)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}

		{
			if canceledcontext.IsCanceled(ctx) {
				return nil
			}

			updateState, ok := patch.getUpdateChange()
			if ok {
				meta, ok := loggermeta.FromContext(ctx)
				if ok {
					meta.KeyVals["function"] = "ApplyUpdateChange"
					defer delete(meta.KeyVals, "function")
				}
				err := r.ApplyUpdateChange(ctx, obj, updateState)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}
	}

	return nil
}

func (f *Framework) bootWithError(ctx context.Context) error {
	if f.tpr != nil {
		f.logger.LogCtx(ctx, "debug", "ensuring third party resource exists")

		err := f.tpr.CreateAndWait()
		if tpr.IsAlreadyExists(err) {
			f.logger.LogCtx(ctx, "debug", "third party resource already exists")
		} else if err != nil {
			return microerror.Mask(err)
		} else {
			f.logger.LogCtx(ctx, "debug", "created third party resource")
		}

		f.tpr.CollectMetrics(context.TODO())
	}

	if f.crd != nil {
		f.logger.LogCtx(ctx, "debug", "ensuring custom resource definition exists")

		err := f.crdClient.Ensure(ctx, f.crd, f.backOffFactory())
		if err != nil {
			return microerror.Mask(err)
		}

		f.logger.LogCtx(ctx, "debug", "ensured custom resource definition")

		// TODO collect metrics
	}

	f.logger.LogCtx(ctx, "debug", "starting list/watch")

	deleteChan, updateChan, errChan := f.informer.Watch(ctx)
	f.ProcessEvents(ctx, deleteChan, updateChan, errChan)

	return nil
}
