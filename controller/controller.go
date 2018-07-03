package controller

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
	"k8s.io/client-go/rest"

	"github.com/giantswarm/operatorkit/client/k8scrdclient"
	"github.com/giantswarm/operatorkit/controller/context/finalizerskeptcontext"
	"github.com/giantswarm/operatorkit/controller/context/reconciliationcanceledcontext"
	"github.com/giantswarm/operatorkit/controller/context/resourcecanceledcontext"
	"github.com/giantswarm/operatorkit/informer"
)

const (
	loggerResourceKey = "resource"
)

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
	// RESTClient needs to be configured with a serializer capable of serializing
	// and deserializing the object which is watched by the informer. Otherwise
	// deserialization will fail when trying to add a finalizer.
	//
	// For standard k8s object this is going to be e.g.
	//
	// 		k8sClient.CoreV1().RESTClient()
	//
	// For CRs of giantswarm this is going to be e.g.
	//
	// 		g8sClient.CoreV1alpha1().RESTClient()
	//
	RESTClient rest.Interface

	BackOffFactory func() backoff.BackOff
	// Name is the name which the controller uses on finalizers for resources.
	// The name used should be unique in the kubernetes cluster, to ensure that
	// two operators which handle the same resource add two distinct finalizers.
	Name string
}

type Controller struct {
	crd            *apiextensionsv1beta1.CustomResourceDefinition
	crdClient      *k8scrdclient.CRDClient
	informer       informer.Interface
	restClient     rest.Interface
	logger         micrologger.Logger
	resourceRouter *ResourceRouter

	bootOnce       sync.Once
	errorCollector chan error
	mutex          sync.Mutex

	backOffFactory func() backoff.BackOff
	name           string
}

// New creates a new configured operator controller.
func New(config Config) (*Controller, error) {
	if config.CRD != nil && config.CRDClient == nil || config.CRD == nil && config.CRDClient != nil {
		return nil, microerror.Maskf(invalidConfigError, "config.CRD and config.CRDClient must not be empty when either given")
	}
	if config.Informer == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Informer must not be empty")
	}
	if config.RESTClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
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
	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Name must not be empty")
	}

	c := &Controller{
		crd:            config.CRD,
		crdClient:      config.CRDClient,
		informer:       config.Informer,
		restClient:     config.RESTClient,
		logger:         config.Logger,
		resourceRouter: config.ResourceRouter,

		bootOnce:       sync.Once{},
		errorCollector: make(chan error, 1),
		mutex:          sync.Mutex{},

		backOffFactory: config.BackOffFactory,
		name:           config.Name,
	}

	return c, nil
}

func (c *Controller) Boot() {
	ctx := context.TODO()

	c.bootOnce.Do(func() {
		operation := func() error {
			err := c.bootWithError(ctx)
			if err != nil {
				return microerror.Mask(err)
			}

			return nil
		}

		notifier := func(err error, d time.Duration) {
			c.logger.LogCtx(ctx, "level", "warning", "message", "retrying controller boot due to error", "stack", fmt.Sprintf("%#v", err))
		}

		err := backoff.RetryNotify(operation, c.backOffFactory(), notifier)
		if err != nil {
			c.logger.LogCtx(ctx, "level", "error", "message", "stop controller boot retries due to too many errors", "stack", fmt.Sprintf("%#v", err))
			os.Exit(1)
		}
	})
}

// DeleteFunc executes the controller's ProcessDelete function.
func (c *Controller) DeleteFunc(obj interface{}) {
	// DeleteFunc/UpdateFunc is synchronized to make sure only one of them is
	// executed at a time. DeleteFunc/UpdateFunc is not thread safe. This is
	// important because the source of truth for an operator are the reconciled
	// resources. In case we would run the operator logic in parallel, we would
	// run into race conditions.
	c.mutex.Lock()
	defer c.mutex.Unlock()

	resourceSet, err := c.resourceRouter.ResourceSet(obj)
	if IsNoResourceSet(err) {
		// In case the resource router is not able to find any resource set to
		// handle the reconciled runtime object, we stop here. Note that we just
		// remove the finalizer regardless because at this point there will never be
		// a chance to remove it otherwhise because nobody wanted to handle this
		// runtime object anyway.

		c.logger.Log("level", "debug", "message", "removing finalizer from runtime object")

		err = c.removeFinalizer(obj)
		if err != nil {
			c.logger.Log("level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
			return
		}

		c.logger.Log("level", "debug", "message", "removed finalizer from runtime object")

		return
	} else if err != nil {
		c.logger.Log("level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}

	ctx, err := resourceSet.InitCtx(context.Background(), obj)
	if err != nil {
		c.logger.Log("level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}

	{
		meta, ok := loggermeta.FromContext(ctx)
		if !ok {
			meta = loggermeta.New()
		}
		meta.KeyVals["event"] = "delete"

		ctx = loggermeta.NewContext(ctx, meta)
	}

	err = ProcessDelete(ctx, obj, resourceSet.Resources())
	if err != nil {
		c.errorCollector <- err
		c.logger.LogCtx(ctx, "level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}

	if !finalizerskeptcontext.IsKept(ctx) {
		c.logger.LogCtx(ctx, "level", "debug", "message", "removing finalizer from runtime object")

		err = c.removeFinalizer(obj)
		if err != nil {
			c.logger.LogCtx(ctx, "level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
			return
		}

		c.logger.LogCtx(ctx, "level", "debug", "message", "removed finalizer from runtime object")
	} else {
		c.logger.LogCtx(ctx, "level", "debug", "message", "not removing finalizer from runtime object due to request of keeping it")
	}
}

// ProcessEvents takes the event channels created by the operatorkit informer
// and executes the controller's event functions accordingly.
func (c *Controller) ProcessEvents(ctx context.Context, deleteChan chan watch.Event, updateChan chan watch.Event, errChan chan error) {
	operation := func() error {
		for {
			select {
			case e := <-deleteChan:
				t := prometheus.NewTimer(controllerHistogram.WithLabelValues("delete"))
				c.DeleteFunc(e.Object)
				t.ObserveDuration()
			case e := <-updateChan:
				t := prometheus.NewTimer(controllerHistogram.WithLabelValues("update"))
				c.UpdateFunc(nil, e.Object)
				t.ObserveDuration()
			case err := <-errChan:
				if IsStatusForbidden(err) {
					return microerror.Maskf(statusForbiddenError, "controller might be missing RBAC rule for %s CRD", c.crd.Name)
				} else if err != nil {
					return microerror.Mask(err)
				}
			case <-ctx.Done():
				return nil
			}
		}
	}

	notifier := func(err error, d time.Duration) {
		c.logger.LogCtx(ctx, "level", "warning", "message", "retrying event processing due to error", "stack", fmt.Sprintf("%#v", err))
	}

	err := backoff.RetryNotify(operation, c.backOffFactory(), notifier)
	if err != nil {
		c.logger.LogCtx(ctx, "level", "error", "message", "stop event processing retries due to too many errors", "stack", fmt.Sprintf("%#v", err))
		os.Exit(1)
	}
}

// UpdateFunc executes the controller's ProcessUpdate function.
func (c *Controller) UpdateFunc(oldObj, newObj interface{}) {
	obj := newObj

	// DeleteFunc/UpdateFunc is synchronized to make sure only one of them is
	// executed at a time. DeleteFunc/UpdateFunc is not thread safe. This is
	// important because the source of truth for an operator are the reconciled
	// resources. In case we would run the operator logic in parallel, we would
	// run into race conditions.
	c.mutex.Lock()
	defer c.mutex.Unlock()

	resourceSet, err := c.resourceRouter.ResourceSet(obj)
	if IsNoResourceSet(err) {
		// In case the resource router is not able to find any resource set to
		// handle the reconciled runtime object, we stop here.
		return
	} else if err != nil {
		c.logger.Log("level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}

	ctx, err := resourceSet.InitCtx(context.Background(), obj)
	if err != nil {
		c.logger.Log("level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}

	{
		meta, ok := loggermeta.FromContext(ctx)
		if !ok {
			meta = loggermeta.New()
		}
		meta.KeyVals["event"] = "update"

		ctx = loggermeta.NewContext(ctx, meta)
	}

	ok, err := c.addFinalizer(obj)
	if IsInvalidRESTClient(err) {
		panic("invalid REST client configured for controller")
	} else if err != nil {
		c.logger.LogCtx(ctx, "level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}
	if ok {
		// A finalizer was added, this causes a new update event, so we stop
		// reconciling here and will pick up the new event.
		c.logger.LogCtx(ctx, "level", "debug", "message", "stop reconciliation due to finalizer added")
		return
	}

	err = ProcessUpdate(ctx, obj, resourceSet.Resources())
	if err != nil {
		c.errorCollector <- err
		c.logger.LogCtx(ctx, "level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}
}

func (c *Controller) bootWithError(ctx context.Context) error {
	if c.crd != nil {
		c.logger.LogCtx(ctx, "level", "debug", "message", "ensuring custom resource definition exists")

		err := c.crdClient.EnsureCreated(ctx, c.crd, c.backOffFactory())
		if err != nil {
			return microerror.Mask(err)
		}

		c.logger.LogCtx(ctx, "level", "debug", "message", "ensured custom resource definition exists")
	}

	{
		c.logger.LogCtx(ctx, "level", "debug", "message", "booting informer")

		err := c.informer.Boot(ctx)
		if err != nil {
			return microerror.Mask(err)
		}

		c.logger.LogCtx(ctx, "level", "debug", "message", "booted informer")
	}

	go func() {
		resetWait := c.informer.ResyncPeriod() * 3

		for {
			select {
			case <-c.errorCollector:
				controllerErrorGauge.Inc()
			case <-time.After(resetWait):
				controllerErrorGauge.Set(0)
			}
		}
	}()

	c.logger.LogCtx(ctx, "level", "debug", "message", "starting list-watch")

	deleteChan, updateChan, errChan := c.informer.Watch(ctx)
	c.ProcessEvents(ctx, deleteChan, updateChan, errChan)

	return nil
}

// ProcessDelete is a drop-in for an informer's DeleteFunc. It receives the
// custom object observed during custom resource watches and anything that
// implements Resource. ProcessDelete takes care about all necessary
// reconciliation logic for delete events.
//
//     func deleteFunc(obj interface{}) {
//         err := c.ProcessDelete(obj, resources)
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

	ctx = reconciliationcanceledcontext.NewContext(ctx, make(chan struct{}))

	defer unsetLoggerCtxValue(ctx, loggerResourceKey)

	for _, r := range resources {
		ctx = setLoggerCtxValue(ctx, loggerResourceKey, r.Name())
		ctx = resourcecanceledcontext.NewContext(ctx, make(chan struct{}))

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
//         err := c.ProcessUpdate(newObj, resources)
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

	ctx = reconciliationcanceledcontext.NewContext(ctx, make(chan struct{}))

	defer unsetLoggerCtxValue(ctx, loggerResourceKey)

	for _, r := range resources {
		ctx = setLoggerCtxValue(ctx, loggerResourceKey, r.Name())
		ctx = resourcecanceledcontext.NewContext(ctx, make(chan struct{}))

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
