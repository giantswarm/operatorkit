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
	"github.com/prometheus/client_golang/prometheus"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/giantswarm/operatorkit/client/k8scrdclient"
	"github.com/giantswarm/operatorkit/framework/context/reconciliationcanceledcontext"
	"github.com/giantswarm/operatorkit/informer"
)

// Config represents the configuration used to create a new operator framework.
type Config struct {
	CRD       *apiextensionsv1beta1.CustomResourceDefinition
	CRDClient *k8scrdclient.CRDClient
	Informer  informer.Interface
	Logger    micrologger.Logger
	// ResourceRouter determines which resource set to use on reconciliation based
	// on its own implementation. A resource router is to decide which resource
	// set to execute. A resource set provides a specific function to initialize
	// the request context and a list of resources to be executed for a
	// reconciliation loop. That way each runtime object being reconciled is
	// executed against a desired list of resources. Since runtime objects may
	// differ in version and/or structure the resource router enables custom
	// inspection before each reconciliation loop. That way the complete list of
	// resources being executed for the received runtime object can be versioned
	// and different resources can be executed depending on the runtime object
	// being reconciled.
	ResourceRouter *ResourceRouter

	BackOffFactory func() backoff.BackOff
}

type Framework struct {
	crd            *apiextensionsv1beta1.CustomResourceDefinition
	crdClient      *k8scrdclient.CRDClient
	informer       informer.Interface
	logger         micrologger.Logger
	resourceRouter *ResourceRouter

	bootOnce sync.Once
	mutex    sync.Mutex

	backOffFactory func() backoff.BackOff
}

// New creates a new configured operator framework.
func New(config Config) (*Framework, error) {
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

	if config.BackOffFactory == nil {
		config.BackOffFactory = DefaultBackOffFactory()
	}

	f := &Framework{
		crd:            config.CRD,
		crdClient:      config.CRDClient,
		informer:       config.Informer,
		logger:         config.Logger,
		resourceRouter: config.ResourceRouter,

		bootOnce: sync.Once{},
		mutex:    sync.Mutex{},

		backOffFactory: config.BackOffFactory,
	}

	return f, nil
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

// deleteFunc executes the framework's processDelete function.
func (f *Framework) deleteFunc(obj interface{}) {
	// deleteFunc/updateFunc is synchronized to make sure only one of them is
	// executed at a time. deleteFunc/updateFunc is not thread safe. This is
	// important because the source of truth for an operator are the reconciled
	// resources. In case we would run the operator logic in parallel, we would
	// run into race conditions.
	f.mutex.Lock()
	defer f.mutex.Unlock()

	resourceSet, err := f.resourceRouter.ResourceSet(obj)
	if err != nil {
		f.logger.Log("error", fmt.Sprintf("%#v", err), "event", "delete")
		return
	}

	ctx, err := resourceSet.InitCtx(context.Background(), obj)
	if err != nil {
		f.logger.Log("error", fmt.Sprintf("%#v", err), "event", "delete")
		return
	}

	f.logger.LogCtx(ctx, "action", "start", "component", "operatorkit", "function", "processDelete")

	err = processDelete(ctx, obj, resourceSet.Resources())
	if err != nil {
		f.logger.LogCtx(ctx, "error", fmt.Sprintf("%#v", err), "event", "delete")
		return
	}

	f.logger.LogCtx(ctx, "action", "end", "component", "operatorkit", "function", "processDelete")
}

// updateFunc executes the framework's processUpdate function.
func (f *Framework) updateFunc(oldObj, newObj interface{}) {
	obj := newObj

	// deleteFunc/updateFunc is synchronized to make sure only one of them is
	// executed at a time. deleteFunc/updateFunc is not thread safe. This is
	// important because the source of truth for an operator are the reconciled
	// resources. In case we would run the operator logic in parallel, we would
	// run into race conditions.
	f.mutex.Lock()
	defer f.mutex.Unlock()

	resourceSet, err := f.resourceRouter.ResourceSet(obj)
	if err != nil {
		f.logger.Log("error", fmt.Sprintf("%#v", err), "event", "update")
		return
	}

	ctx, err := resourceSet.InitCtx(context.Background(), obj)
	if err != nil {
		f.logger.Log("error", fmt.Sprintf("%#v", err), "event", "update")
		return
	}

	f.logger.LogCtx(ctx, "action", "start", "component", "operatorkit", "function", "processUpdate")

	err = processUpdate(ctx, obj, resourceSet.Resources())
	if err != nil {
		f.logger.LogCtx(ctx, "error", fmt.Sprintf("%#v", err), "event", "update")
		return
	}

	f.logger.LogCtx(ctx, "action", "end", "component", "operatorkit", "function", "processUpdate")
}

// processEvents takes the event channels created by the operatorkit informer
// and executes the framework's event functions accordingly.
func (f *Framework) processEvents(ctx context.Context, deleteChan chan watch.Event, updateChan chan watch.Event, errChan chan error) {
	operation := func() error {
		for {
			select {
			case e := <-deleteChan:
				t := prometheus.NewTimer(frameworkHistogram.WithLabelValues("delete"))
				f.deleteFunc(e.Object)
				t.ObserveDuration()
			case e := <-updateChan:
				t := prometheus.NewTimer(frameworkHistogram.WithLabelValues("update"))
				f.updateFunc(nil, e.Object)
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

func (f *Framework) bootWithError(ctx context.Context) error {
	if f.crd != nil {
		f.logger.LogCtx(ctx, "debug", "ensuring custom resource definition exists")

		err := f.crdClient.EnsureCreated(ctx, f.crd, f.backOffFactory())
		if err != nil {
			return microerror.Mask(err)
		}

		f.logger.LogCtx(ctx, "debug", "ensured custom resource definition")

		// TODO collect metrics
	}

	f.logger.LogCtx(ctx, "debug", "starting list/watch")

	deleteChan, updateChan, errChan := f.informer.Watch(ctx)
	f.processEvents(ctx, deleteChan, updateChan, errChan)

	return nil
}

// processDelete is a drop-in for an informer's DeleteFunc.
//
//	func deleteFunc(obj interface{}) {
//		err := processDelete(obj, resources)
//		if err != nil {
//			// error handling here
//		}
//	}
//
//	&cache.ResourceEventHandlerFuncs{
//		DeleteFunc:    deleteFunc,
//	}
//
func processDelete(ctx context.Context, obj interface{}, resources []Resource) error {
	if len(resources) == 0 {
		return microerror.Maskf(executionFailedError, "resources must not be empty")
	}

	ctx = reconciliationcanceledcontext.NewContext(ctx, make(chan struct{}))
	// Garbage collect. TODO use bool to not have to do so.
	defer reconciliationcanceledcontext.SetCanceled(ctx)

	for _, r := range resources {
		err := r.EnsureDeleted(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
	}

	return nil
}

// processUpdate is a drop-in for an informer's AddFunc and updateFunc.
//
//	func addFunc(obj interface{}) {
//		err := processUpdate(obj, resources)
//		if err != nil {
//			// error handling here
//		}
//	}
//
//	func updateFunc(oldObj, newObj interface{}) {
//		err := processUpdate(newObj, resources)
//		if err != nil {
//			// error handling here
//		}
//	}
//
//	&cache.ResourceEventHandlerFuncs{
//		AddFunc:       addFunc,
//		UpdateFunc:    updateFunc,
//	}
//
func processUpdate(ctx context.Context, obj interface{}, resources []Resource) error {
	if len(resources) == 0 {
		return microerror.Maskf(executionFailedError, "resources must not be empty")
	}

	ctx = reconciliationcanceledcontext.NewContext(ctx, make(chan struct{}))
	// Garbage collect. TODO use bool to not have to do so.
	defer reconciliationcanceledcontext.SetCanceled(ctx)

	for _, r := range resources {
		err := r.EnsureCreated(ctx, obj)
		if err != nil {
			return microerror.Mask(err)
		}

		if reconciliationcanceledcontext.IsCanceled(ctx) {
			return nil
		}
	}

	return nil
}
