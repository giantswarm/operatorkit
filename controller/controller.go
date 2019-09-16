package controller

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/micrologger/loggermeta"
	"github.com/prometheus/client_golang/prometheus"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/operatorkit/client/k8scrdclient"
	"github.com/giantswarm/operatorkit/controller/context/reconciliationcanceledcontext"
	"github.com/giantswarm/operatorkit/controller/context/resourcecanceledcontext"
	"github.com/giantswarm/operatorkit/informer"
	"github.com/giantswarm/operatorkit/resource"
)

const (
	loggerKeyController = "controller"
	loggerKeyEvent      = "event"
	loggerKeyLoop       = "loop"
	loggerKeyObject     = "object"
	loggerKeyResource   = "resource"
	loggerKeyVersion    = "version"
)

type Config struct {
	CRD       *apiextensionsv1beta1.CustomResourceDefinition
	CRDClient *k8scrdclient.CRDClient
	Informer  informer.Interface
	Logger    micrologger.Logger
	// ResourceSets is a list of resource sets. A resource set provides a specific
	// function to initialize the request context and a list of resources to be
	// executed for a reconciliation loop. That way each runtime object being
	// reconciled is executed against a desired list of resources. Since runtime
	// objects may differ in version and/or structure the resource router enables
	// custom inspection before each reconciliation loop. That way the complete
	// list of resources being executed for the received runtime object can be
	// versioned and different resources can be executed depending on the runtime
	// object being reconciled.
	ResourceSets []*ResourceSet
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

	BackOffFactory func() backoff.Interface
	// Name is the name which the controller uses on finalizers for resources.
	// The name used should be unique in the kubernetes cluster, to ensure that
	// two operators which handle the same resource add two distinct finalizers.
	Name string
}

type Controller struct {
	crd          *apiextensionsv1beta1.CustomResourceDefinition
	crdClient    *k8scrdclient.CRDClient
	informer     informer.Interface
	restClient   rest.Interface
	logger       micrologger.Logger
	resourceSets []*ResourceSet

	bootOnce               sync.Once
	booted                 chan struct{}
	errorCollector         chan error
	removedFinalizersCache *stringCache
	mutex                  sync.Mutex

	backOffFactory func() backoff.Interface
	name           string
}

// New creates a new configured operator controller.
func New(config Config) (*Controller, error) {
	if config.CRD != nil && config.CRDClient == nil || config.CRD == nil && config.CRDClient != nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.CRD and %T.CRDClient must not be empty when either given", config, config)
	}
	if config.Informer == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Informer must not be empty", config)
	}
	if config.RESTClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if len(config.ResourceSets) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.ResourceSets must not be empty", config)
	}

	if config.BackOffFactory == nil {
		config.BackOffFactory = func() backoff.Interface { return backoff.NewMaxRetries(7, 1*time.Second) }
	}
	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Name must not be empty", config)
	}

	c := &Controller{
		crd:          config.CRD,
		crdClient:    config.CRDClient,
		informer:     config.Informer,
		restClient:   config.RESTClient,
		logger:       config.Logger,
		resourceSets: config.ResourceSets,

		bootOnce:               sync.Once{},
		booted:                 make(chan struct{}),
		errorCollector:         make(chan error, 1),
		removedFinalizersCache: newStringCache(config.Informer.ResyncPeriod() * 3),
		mutex:                  sync.Mutex{},

		backOffFactory: config.BackOffFactory,
		name:           config.Name,
	}

	return c, nil
}

func (c *Controller) Boot(ctx context.Context) {
	ctx = setLoggerCtxValue(ctx, loggerKeyController, c.name)

	c.bootOnce.Do(func() {
		operation := func() error {
			err := c.bootWithError(ctx)
			if err != nil {
				return microerror.Mask(err)
			}

			return nil
		}

		notifier := backoff.NewNotifier(c.logger, ctx)

		err := backoff.RetryNotify(operation, c.backOffFactory(), notifier)
		if err != nil {
			c.logger.LogCtx(ctx, "level", "error", "message", "stop controller boot retries due to too many errors", "stack", fmt.Sprintf("%#v", err))
			os.Exit(1)
		}
	})
}

func (c *Controller) Booted() chan struct{} {
	return c.booted
}

// DeleteFunc executes the controller's ProcessDelete function.
func (c *Controller) DeleteFunc(obj interface{}) {
	ctx := context.Background()
	c.deleteFunc(ctx, obj)
}

func (c *Controller) deleteFunc(ctx context.Context, obj interface{}) {
	// DeleteFunc/UpdateFunc is synchronized to make sure only one of them is
	// executed at a time. DeleteFunc/UpdateFunc is not thread safe. This is
	// important because the source of truth for an operator are the reconciled
	// resources. In case we would run the operator logic in parallel, we would
	// run into race conditions.
	c.mutex.Lock()
	defer c.mutex.Unlock()

	rs, err := c.resourceSet(obj)
	if IsNoResourceSet(err) {
		// In case the resource router is not able to find any resource set to
		// handle the reconciled runtime object, we stop here. Note that we just
		// remove the finalizer regardless because at this point there will never be
		// a chance to remove it otherwhise because nobody wanted to handle this
		// runtime object anyway. Otherwise we can end up in deadlock
		// trying to reconcile this object over and over.
		err = c.removeFinalizer(ctx, obj)
		if err != nil {
			c.logger.LogCtx(ctx, "level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
			return
		}

		c.logger.LogCtx(ctx, "level", "debug", "message", "did not find any resource set")
		c.logger.LogCtx(ctx, "level", "debug", "message", "canceling reconciliation")
		return

	} else if err != nil {
		c.logger.LogCtx(ctx, "level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}

	{
		// Memorize the old context. We need to use it for logging in
		// case of InitCtx failure because the context returned by
		// failing InitCtx can be nil.
		oldCtx := ctx
		ctx, err = rs.InitCtx(ctx, obj)
		if err != nil {
			c.logger.LogCtx(oldCtx, "level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
			return
		}
	}

	hasFinalizer, err := c.hasFinalizer(ctx, obj)
	if err != nil {
		c.logger.LogCtx(ctx, "level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}
	if hasFinalizer {
		err = ProcessDelete(ctx, obj, rs.Resources())
		if err != nil {
			c.errorCollector <- err
			c.logger.LogCtx(ctx, "level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
			return
		}
	} else {
		c.logger.LogCtx(ctx, "level", "debug", "message", "did not find any finalizer")
		c.logger.LogCtx(ctx, "level", "debug", "message", "canceling reconciliation")
		return
	}

	err = c.removeFinalizer(ctx, obj)
	if err != nil {
		c.logger.LogCtx(ctx, "level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}
}

// ProcessEvents takes the event channels created by the operatorkit informer
// and executes the controller's event functions accordingly.
func (c *Controller) ProcessEvents(ctx context.Context, deleteChan chan watch.Event, updateChan chan watch.Event, errChan chan error) error {
	loop := -1

	for {
		loop++

		// Set loop specific logger context.
		{
			ctx = setLoggerCtxValue(ctx, loggerKeyLoop, strconv.Itoa(loop))
		}

		select {
		case e := <-deleteChan:
			event := "delete"

			t := prometheus.NewTimer(controllerHistogram.WithLabelValues(event))

			// Set event specific logger context.
			{
				ctx = setLoggerCtxValue(ctx, loggerKeyEvent, event)

				accessor, err := meta.Accessor(e.Object)
				if err != nil {
					c.logger.LogCtx(ctx, "level", "error", "message", fmt.Sprintf("failed to create accessor %T", e.Object), "stack", fmt.Sprintf("%#v", err))
				} else {
					ctx = setLoggerCtxValue(ctx, loggerKeyObject, accessor.GetSelfLink())
					ctx = setLoggerCtxValue(ctx, loggerKeyVersion, accessor.GetResourceVersion())
				}
			}

			c.logger.LogCtx(ctx, "level", "debug", "message", "reconciling object")
			c.deleteFunc(ctx, e.Object)
			c.logger.LogCtx(ctx, "level", "debug", "message", "reconciled object")

			t.ObserveDuration()
		case e := <-updateChan:
			event := "update"

			t := prometheus.NewTimer(controllerHistogram.WithLabelValues(event))

			// Set event specific logger context.
			{
				ctx = setLoggerCtxValue(ctx, loggerKeyEvent, event)

				accessor, err := meta.Accessor(e.Object)
				if err != nil {
					c.logger.LogCtx(ctx, "level", "error", "message", fmt.Sprintf("failed to create accessor %T", e.Object), "stack", fmt.Sprintf("%#v", err))
				} else {
					ctx = setLoggerCtxValue(ctx, loggerKeyObject, accessor.GetSelfLink())
					ctx = setLoggerCtxValue(ctx, loggerKeyVersion, accessor.GetResourceVersion())
				}
			}

			c.logger.LogCtx(ctx, "level", "debug", "message", "reconciling object")
			c.updateFunc(ctx, e.Object)
			c.logger.LogCtx(ctx, "level", "debug", "message", "reconciled object")

			t.ObserveDuration()
		case err := <-errChan:
			if IsStatusForbidden(err) {
				c.logger.LogCtx(ctx, "level", "error", "message", fmt.Sprintf("controller might be missing RBAC rule for %s CRD", c.crd.Name), "stack", fmt.Sprintf("%#v", err))
			} else if err != nil {
				c.logger.LogCtx(ctx, "level", "error", "message", "failed to watch object", "stack", fmt.Sprintf("%#v", err))
			}

			time.Sleep(time.Second)
		case <-ctx.Done():
			return nil
		}
	}
}

// UpdateFunc executes the controller's ProcessUpdate function.
func (c *Controller) UpdateFunc(oldObj, newObj interface{}) {
	ctx := context.Background()
	c.updateFunc(ctx, newObj)
}

func (c *Controller) updateFunc(ctx context.Context, obj interface{}) {
	// DeleteFunc/UpdateFunc is synchronized to make sure only one of them is
	// executed at a time. DeleteFunc/UpdateFunc is not thread safe. This is
	// important because the source of truth for an operator are the reconciled
	// resources. In case we would run the operator logic in parallel, we would
	// run into race conditions.
	c.mutex.Lock()
	defer c.mutex.Unlock()

	rs, err := c.resourceSet(obj)
	if IsNoResourceSet(err) {
		// In case the resource router is not able to find any resource set to
		// handle the reconciled runtime object, we stop here.
		c.logger.LogCtx(ctx, "level", "debug", "message", "did not find any resource set")
		c.logger.LogCtx(ctx, "level", "debug", "message", "canceling reconciliation")
		return

	} else if err != nil {
		c.logger.LogCtx(ctx, "level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
		return
	}

	{
		// Memorize the old context. We need to use it for logging in
		// case of InitCtx failure because the context returned by
		// failing InitCtx can be nil.
		oldCtx := ctx
		ctx, err = rs.InitCtx(ctx, obj)
		if err != nil {
			c.logger.LogCtx(oldCtx, "level", "error", "message", "stop reconciliation due to error", "stack", fmt.Sprintf("%#v", err))
			return
		}
	}

	ok, err := c.addFinalizer(ctx, obj)
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

	err = ProcessUpdate(ctx, obj, rs.Resources())
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
		resetWait := c.informer.ResyncPeriod() * 4

		for {
			select {
			case <-c.errorCollector:
				controllerErrorGauge.Inc()
			case <-time.After(resetWait):
				controllerErrorGauge.Set(0)
			}
		}
	}()

	// We overwrite the k8s error handlers so they do not intercept our log
	// streams. The format is way easier to parse for us that way. Here we also
	// emit metrics for the occured errors to ensure we create more awareness of
	// anything going wrong in our operators.
	{
		runtime.ErrorHandlers = []func(err error){
			func(err error) {
				// When we see a port forwarding error we ignore it because we cannot do
				// anything about it. Errors like we check here would have to be dealt
				// with in the third party tools we use. The port forwarding in general
				// is broken by design which will go away with Helm 3, soon TM.
				if IsPortforward(err) {
					return
				}

				c.errorCollector <- err
				c.logger.LogCtx(ctx, "level", "error", "message", "caught third party runtime error", "stack", fmt.Sprintf("%#v", err))
			},
		}
	}

	// Initializing the watch gives us all necessary event channels we need to
	// further process them within the controller. Once we got the event channels
	// everything is set up for the operator's reconciliation. We put the
	// controller into a booted state by closing its booted channel so users know
	// when to go ahead. Note that ProcessEvents below blocks the boot process
	// until it fails.
	{
		c.logger.LogCtx(ctx, "level", "debug", "message", "starting list-watch")

		deleteChan, updateChan, errChan := c.informer.Watch(ctx)
		close(c.booted)

		c.logger.LogCtx(ctx, "level", "debug", "message", "processing object events")

		err := c.ProcessEvents(ctx, deleteChan, updateChan, errChan)
		if err != nil {
			return microerror.Mask(err)
		}

		c.logger.LogCtx(ctx, "level", "debug", "message", "processed object events")
	}

	return nil
}

// resourceSet tries to lookup the appropriate resource set based on the
// received runtime object. There might be not any resource set for an observed
// runtime object if an operator uses multiple controllers for reconciliations.
// There must not be multiple resource sets per observed runtime object though.
// If this is the case, ResourceSet returns an error.
func (c *Controller) resourceSet(obj interface{}) (*ResourceSet, error) {
	var found []*ResourceSet

	for _, r := range c.resourceSets {
		if r.Handles(obj) {
			found = append(found, r)
		}
	}

	if len(found) == 0 {
		return nil, microerror.Mask(noResourceSetError)
	}
	if len(found) > 1 {
		return nil, microerror.Mask(tooManyResourceSetsError)
	}

	return found[0], nil
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
func ProcessDelete(ctx context.Context, obj interface{}, resources []resource.Interface) error {
	if len(resources) == 0 {
		return microerror.Maskf(executionFailedError, "resources must not be empty")
	}

	ctx = reconciliationcanceledcontext.NewContext(ctx, make(chan struct{}))

	for _, r := range resources {
		ctx = setLoggerCtxValue(ctx, loggerKeyResource, r.Name())
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
func ProcessUpdate(ctx context.Context, obj interface{}, resources []resource.Interface) error {
	if len(resources) == 0 {
		return microerror.Maskf(executionFailedError, "resources must not be empty")
	}

	ctx = reconciliationcanceledcontext.NewContext(ctx, make(chan struct{}))

	for _, r := range resources {
		ctx = setLoggerCtxValue(ctx, loggerKeyResource, r.Name())
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
		ctx = loggermeta.NewContext(ctx, m)
	}

	m.KeyVals[key] = value

	return ctx
}
