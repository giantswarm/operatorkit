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
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/operatorkit/client/k8scrdclient"
	"github.com/giantswarm/operatorkit/framework/context/reconciliationcanceledcontext"
	"github.com/giantswarm/operatorkit/informer"
)

const (
	loggerResourceKey = "resource"
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
	K8sClient      kubernetes.Interface

	BackOffFactory func() backoff.BackOff
	// Name is the name which the framework uses on finalizers for resources.
	// The name used should be unique in the kubernetes cluster, to ensure that
	// two operators which handle the same resource add two distinct finalizers.
	Name string
}

type Framework struct {
	crd            *apiextensionsv1beta1.CustomResourceDefinition
	crdClient      *k8scrdclient.CRDClient
	informer       informer.Interface
	k8sClient      kubernetes.Interface
	logger         micrologger.Logger
	resourceRouter *ResourceRouter

	bootOnce sync.Once
	mutex    sync.Mutex

	backOffFactory func() backoff.BackOff
	name           string
}

// New creates a new configured operator framework.
func New(config Config) (*Framework, error) {
	if config.CRD != nil && config.CRDClient == nil || config.CRD == nil && config.CRDClient != nil {
		return nil, microerror.Maskf(invalidConfigError, "config.CRD and config.CRDClient must not be empty when either given")
	}
	if config.Informer == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Informer must not be empty")
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Name must not be empty")
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
		k8sClient:      config.K8sClient,
		logger:         config.Logger,
		name:           config.Name,
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
			f.logger.LogCtx(ctx, "function", "Boot", "level", "warning", "message", "retrying framework boot due to error", "stack", fmt.Sprintf("%#v", err))
		}

		err := backoff.RetryNotify(operation, f.backOffFactory(), notifier)
		if err != nil {
			f.logger.LogCtx(ctx, "function", "Boot", "level", "error", "message", "stop framework boot retries due to too many errors", "stack", fmt.Sprintf("%#v", err))
			os.Exit(1)
		}
	})
}

// DeleteFunc executes the framework's ProcessDelete function.
func (f *Framework) DeleteFunc(obj interface{}) {
	// DeleteFunc/UpdateFunc is synchronized to make sure only one of them is
	// executed at a time. DeleteFunc/UpdateFunc is not thread safe. This is
	// important because the source of truth for an operator are the reconciled
	// resources. In case we would run the operator logic in parallel, we would
	// run into race conditions.
	f.mutex.Lock()
	defer f.mutex.Unlock()

	resourceSet, err := f.resourceRouter.ResourceSet(obj)
	if IsNoResourceSet(err) {
		// In case the resource router is not able to find any resource set to
		// handle the reconciled runtime object, we stop here.
		return
	} else if err != nil {
		f.logger.Log("event", "delete", "function", "DeleteFunc", "level", "error", "message", "stop framework reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}

	ctx, err := resourceSet.InitCtx(context.Background(), obj)
	if err != nil {
		f.logger.Log("event", "delete", "function", "DeleteFunc", "level", "error", "message", "stop framework reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}

	err = ProcessDelete(ctx, obj, resourceSet.Resources())
	if err != nil {
		f.logger.LogCtx(ctx, "event", "delete", "function", "DeleteFunc", "level", "error", "message", "stop framework reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}

	err = f.removeFinalizer(ctx, obj)
	if err != nil {
		f.logger.LogCtx(ctx, "event", "delete", "function", "DeleteFunc", "level", "error", "message", "stop framework reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}
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
		f.logger.LogCtx(ctx, "function", "ProcessEvents", "level", "warning", "message", "retrying framework event processing due to error", "stack", fmt.Sprintf("%#v", err))
	}

	err := backoff.RetryNotify(operation, f.backOffFactory(), notifier)
	if err != nil {
		f.logger.LogCtx(ctx, "function", "ProcessEvents", "level", "error", "message", "stop framework event processing retries due to too many errors", "stack", fmt.Sprintf("%#v", err))
		os.Exit(1)
	}
}

// UpdateFunc executes the framework's ProcessUpdate function.
func (f *Framework) UpdateFunc(oldObj, newObj interface{}) {
	obj := newObj

	// DeleteFunc/UpdateFunc is synchronized to make sure only one of them is
	// executed at a time. DeleteFunc/UpdateFunc is not thread safe. This is
	// important because the source of truth for an operator are the reconciled
	// resources. In case we would run the operator logic in parallel, we would
	// run into race conditions.
	f.mutex.Lock()
	defer f.mutex.Unlock()

	ok, err := f.addFinalizer(obj)
	if err != nil {
		f.logger.Log("event", "update", "function", "UpdateFunc", "level", "error", "message", "stop framework reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}
	if ok {
		// A finalizer was added, this causes a new update event, so we stop
		// reconciling here and will pick up the new event.
		f.logger.Log("event", "update", "function", "UpdateFunc", "level", "debug", "message", "stop framework reconciliation due to finalizer added")
		return
	}

	resourceSet, err := f.resourceRouter.ResourceSet(obj)
	if IsNoResourceSet(err) {
		// In case the resource router is not able to find any resource set to
		// handle the reconciled runtime object, we stop here.
		return
	} else if err != nil {
		f.logger.Log("event", "update", "function", "UpdateFunc", "level", "error", "message", "stop framework reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}

	ctx, err := resourceSet.InitCtx(context.Background(), obj)
	if err != nil {
		f.logger.LogCtx(ctx, "event", "update", "function", "UpdateFunc", "level", "error", "message", "stop framework reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}

	err = ProcessUpdate(ctx, obj, resourceSet.Resources())
	if err != nil {
		f.logger.LogCtx(ctx, "event", "update", "function", "UpdateFunc", "level", "error", "message", "stop framework reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}
}

func (f *Framework) bootWithError(ctx context.Context) error {
	if f.crd != nil {
		f.logger.LogCtx(ctx, "function", "bootWithError", "level", "debug", "message", "ensuring custom resource definition exists")

		err := f.crdClient.EnsureCreated(ctx, f.crd, f.backOffFactory())
		if err != nil {
			return microerror.Mask(err)
		}

		f.logger.LogCtx(ctx, "function", "bootWithError", "level", "debug", "message", "ensured custom resource definition exists")

		// TODO collect metrics
	}

	f.logger.LogCtx(ctx, "function", "bootWithError", "level", "debug", "message", "starting list-watch")

	deleteChan, updateChan, errChan := f.informer.Watch(ctx)
	f.ProcessEvents(ctx, deleteChan, updateChan, errChan)

	return nil
}

// ProcessDelete is a drop-in for an informer's DeleteFunc. It receives the
// custom object observed during custom resource watches and anything that
// implements Resource. ProcessDelete takes care about all necessary
// reconciliation logic for delete events.
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

	defer unsetLoggerCtxValue(ctx, loggerResourceKey)

	for _, r := range resources {
		ctx = setLoggerCtxValue(ctx, loggerResourceKey, r.Name())

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

// ProcessUpdate is a drop-in for an informer's UpdateFunc. It receives the new
// custom object observed during custom resource watches and anything that
// implements Resource. ProcessUpdate takes care about all necessary
// reconciliation logic for update events. For complex resources this means
// state has to be created, deleted and updated eventually, in this order.
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

	defer unsetLoggerCtxValue(ctx, loggerResourceKey)

	for _, r := range resources {
		ctx = setLoggerCtxValue(ctx, loggerResourceKey, r.Name())

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

func setLoggerCtxValue(ctx context.Context, key, value string) context.Context {
	m, ok := loggermeta.FromContext(ctx)
	if !ok {
		m = loggermeta.New()
	}
	m.KeyVals[key] = value

	return loggermeta.NewContext(ctx, m)
}

func unsetLoggerCtxValue(ctx context.Context, key string) context.Context {
	m, ok := loggermeta.FromContext(ctx)
	if !ok {
		m = loggermeta.New()
	}
	delete(m.KeyVals, key)

	return loggermeta.NewContext(ctx, m)
}
