package controller

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	apiextensionsscheme "github.com/giantswarm/apiextensions/v3/pkg/clientset/versioned/scheme"
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/k8sclient/v5/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/micrologger/loggermeta"
	"github.com/giantswarm/to"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgruntime "k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/giantswarm/operatorkit/v4/pkg/controller/collector"
	"github.com/giantswarm/operatorkit/v4/pkg/controller/context/cachekeycontext"
	"github.com/giantswarm/operatorkit/v4/pkg/controller/context/finalizerskeptcontext"
	"github.com/giantswarm/operatorkit/v4/pkg/controller/context/reconciliationcanceledcontext"
	"github.com/giantswarm/operatorkit/v4/pkg/controller/context/resourcecanceledcontext"
	"github.com/giantswarm/operatorkit/v4/pkg/controller/context/updateallowedcontext"
	"github.com/giantswarm/operatorkit/v4/pkg/controller/internal/recorder"
	"github.com/giantswarm/operatorkit/v4/pkg/controller/internal/selector"
	"github.com/giantswarm/operatorkit/v4/pkg/controller/internal/sentry"
	"github.com/giantswarm/operatorkit/v4/pkg/resource"
)

const (
	DefaultResyncPeriod = 5 * time.Minute
)

const (
	DisableMetricsServing = "0"
)

const (
	loggerKeyController = "controller"
	loggerKeyEvent      = "event"
	loggerKeyLoop       = "loop"
	loggerKeyObject     = "object"
	loggerKeyResource   = "resource"
	loggerKeyVersion    = "version"
)

var (
	defaultPauseAnnotations = map[string]string{
		"cluster.x-k8s.io/paused":          "true",
		"operatorkit.giantswarm.io/paused": "true",
	}
)

type Config struct {
	// InitCtx is deprecated and should not be used anymore.
	InitCtx func(ctx context.Context, obj interface{}) (context.Context, error)
	// K8sClient is the client collection used to setup and manage certain
	// operatorkit primitives. The Controller Client is used to fetch runtime
	// objects. It therefore must be properly configured using the AddToScheme
	// option. The REST Client is used to patch finalizers on runtime objects.
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
	// NewRuntimeObjectFunc returns a new initialized pointer of a type
	// implementing the runtime object interface. The object returned is used with
	// the controller-runtime client to fetch the latest version of the object
	// itself. That way we can manage all runtime objects in a somewhat generic
	// way. See the example below.
	//
	//     func() pkgruntime.Object {
	//        return new(corev1.ConfigMap)
	//     }
	//
	NewRuntimeObjectFunc func() pkgruntime.Object
	// Pause is a map of additional pausing annotations, defining their
	// key-value pairs. This can be used to stop reconciliation in case the
	// configured pausing annotations are present in the reconciled runtime
	// object watched by operatorkit.
	Pause map[string]string
	// Resources is the list of controller resources being executed on runtime
	// object reconciliation. Resources are executed in given order.
	Resources []resource.Interface
	// Selector is used to filter objects before passing them to the controller.
	Selector selector.Selector

	// Name is the name which the controller uses on finalizers for resources.
	// The name used should be unique in the kubernetes cluster, to ensure that
	// two operators which handle the same resource add two distinct finalizers.
	Name string
	// Namespace is where the controller would reconcile the runtime objects.
	// Empty string means all namespaces.
	Namespace string
	// ResyncPeriod is the duration after which a complete sync with all known
	// runtime objects the controller watches is performed. Defaults to
	// DefaultResyncPeriod.
	ResyncPeriod time.Duration
	// SentryDSN is the optional URL used to forward runtime errors to the sentry.io service.
	// If this field is empty, logs will not be forwarded.
	SentryDSN string
}

type Controller struct {
	event                recorder.Interface
	initCtx              func(ctx context.Context, obj interface{}) (context.Context, error)
	k8sClient            k8sclient.Interface
	logger               micrologger.Logger
	newRuntimeObjectFunc func() pkgruntime.Object
	pause                map[string]string
	resources            []resource.Interface
	selector             selector.Selector

	backOffFactory         func() backoff.Interface
	bootOnce               sync.Once
	booted                 chan struct{}
	collector              *collector.Set
	loop                   int64
	removedFinalizersCache *stringCache
	sentry                 sentry.Interface

	name         string
	namespace    string
	resyncPeriod time.Duration
}

// New creates a new configured operator controller.
func New(config Config) (*Controller, error) {
	if config.InitCtx == nil {
		config.InitCtx = func(ctx context.Context, obj interface{}) (context.Context, error) {
			return ctx, nil
		}
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.NewRuntimeObjectFunc == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.NewRuntimeObjectFunc must not be empty", config)
	}
	{
		if config.Pause == nil {
			config.Pause = map[string]string{}
		}
		for k, v := range defaultPauseAnnotations {
			config.Pause[k] = v
		}
	}
	if len(config.Resources) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resources must not be empty", config)
	}
	if config.Selector == nil {
		config.Selector = selector.NewSelectorEverything()
	}

	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Name must not be empty", config)
	}
	if config.ResyncPeriod == 0 {
		config.ResyncPeriod = DefaultResyncPeriod
	}

	var err error

	var collectorSet *collector.Set
	{
		c := collector.SetConfig{
			Logger:               config.Logger,
			K8sClient:            config.K8sClient,
			NewRuntimeObjectFunc: config.NewRuntimeObjectFunc,
			Selector:             config.Selector,

			Controller: config.Name,
		}

		collectorSet, err = collector.NewSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var eventRecorder recorder.Interface
	{
		c := recorder.Config{
			K8sClient: config.K8sClient,

			Component: config.Name,
		}

		eventRecorder = recorder.New(c)
	}

	var sentryClient sentry.Interface
	{
		c := sentry.Config{
			DSN: config.SentryDSN,
		}

		sentryClient, err = sentry.New(c)
		if err != nil {
			// Error during sentry initialization.
			return nil, microerror.Mask(err)
		}
	}

	c := &Controller{
		event:                eventRecorder,
		initCtx:              config.InitCtx,
		k8sClient:            config.K8sClient,
		logger:               config.Logger,
		newRuntimeObjectFunc: config.NewRuntimeObjectFunc,
		pause:                config.Pause,
		resources:            config.Resources,
		selector:             config.Selector,

		backOffFactory:         func() backoff.Interface { return backoff.NewMaxRetries(7, 1*time.Second) },
		bootOnce:               sync.Once{},
		booted:                 make(chan struct{}),
		collector:              collectorSet,
		loop:                   -1,
		removedFinalizersCache: newStringCache(config.ResyncPeriod * 3),
		sentry:                 sentryClient,

		name:         config.Name,
		namespace:    config.Namespace,
		resyncPeriod: config.ResyncPeriod,
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
			c.sentry.Capture(ctx, err)
			c.logger.Errorf(ctx, err, "stop controller boot retries due to too many errors")
			os.Exit(1)
		}
	})
}

func (c *Controller) Booted() chan struct{} {
	return c.booted
}

// Reconcile implements the reconciler given to the controller-runtime
// controller. Reconcile never returns any error as we deal with them in
// operatorkit internally.
func (c *Controller) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()

	// Add common keys to the logger context.
	{
		loop := strconv.FormatInt(atomic.AddInt64(&c.loop, 1), 10)

		ctx = cachekeycontext.NewContext(ctx, fmt.Sprintf("%s-%s", c.name, loop))
		ctx = finalizerskeptcontext.NewContext(ctx, make(chan struct{}))
		ctx = updateallowedcontext.NewContext(ctx, make(chan struct{}))

		ctx = setLoggerCtxValue(ctx, loggerKeyLoop, loop)
		ctx = setLoggerCtxValue(ctx, loggerKeyController, c.name)
	}

	obj := c.newRuntimeObjectFunc()
	err := c.k8sClient.CtrlClient().Get(ctx, req.NamespacedName, obj)
	if errors.IsNotFound(err) {
		// At this point the controller-runtime cache dispatches a runtime object
		// which is already being deleted, which is why it cannot be found here
		// anymore. We then likely perceive the last delete event of that runtime
		// object and it got purged from the controller-runtime cache. We do not
		// need to log these errors and just stop processing here in a more graceful
		// way.
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, microerror.Mask(err)
	}

	res, err := c.reconcile(ctx, req, obj)
	if err != nil {
		// Microerror creates an error event on the object when kind and description is set.
		c.event.Emit(ctx, obj, err)
		errorGauge.Inc()
		c.sentry.Capture(ctx, err)
		c.logger.Errorf(ctx, err, "failed to reconcile")
		return reconcile.Result{}, nil
	}

	reportLastReconciled(obj, c.name)

	return res, nil
}

func (c *Controller) bootWithError(ctx context.Context) error {
	var err error

	// Boot the collector.
	err = c.collector.Boot(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	go func() {
		for {
			resetWait := c.resyncPeriod * 4
			time.Sleep(resetWait)
			errorGauge.Set(0)
		}
	}()

	// We overwrite the k8s error handlers so they do not intercept our log
	// streams. The format is way easier to parse for us that way. Here we also
	// emit metrics for the occured errors to ensure we create more awareness of
	// anything going wrong in our operators.
	{
		utilruntime.ErrorHandlers = []func(err error){
			func(err error) {
				// When we see a port forwarding error we ignore it because we cannot do
				// anything about it. Errors like we check here would have to be dealt
				// with in the third party tools we use. The port forwarding in general
				// is broken by design which will go away with Helm 3, soon TM.
				if IsPortforward(err) {
					return
				}

				errorGauge.Inc()
				c.logger.Errorf(ctx, err, "caught third party runtime error")
			},
		}
	}

	var mgr manager.Manager
	{
		o := manager.Options{
			// MetricsBindAddress is set to 0 in order to disable it. We do this
			// ourselves.
			MetricsBindAddress: DisableMetricsServing,
			Namespace:          c.namespace,
			SyncPeriod:         to.DurationP(c.resyncPeriod),
		}

		mgr, err = manager.New(c.k8sClient.RESTConfig(), o)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	{
		// We build our controller and set up its reconciliation.
		// We use the Complete() method instead of Build() because we don't
		// need the controller instance.
		err = builder.
			ControllerManagedBy(mgr).
			For(c.newRuntimeObjectFunc()).
			WithOptions(controller.Options{
				MaxConcurrentReconciles: 1,
				Reconciler:              c,
			}).
			WithEventFilter(predicate.Funcs{
				CreateFunc:  func(e event.CreateEvent) bool { return c.selector.Matches(selector.NewLabels(e.Meta.GetLabels())) },
				DeleteFunc:  func(e event.DeleteEvent) bool { return c.selector.Matches(selector.NewLabels(e.Meta.GetLabels())) },
				UpdateFunc:  func(e event.UpdateEvent) bool { return c.selector.Matches(selector.NewLabels(e.MetaNew.GetLabels())) },
				GenericFunc: func(e event.GenericEvent) bool { return c.selector.Matches(selector.NewLabels(e.Meta.GetLabels())) },
			}).
			Complete(c)
		if err != nil {
			return microerror.Mask(err)
		}

		// We put the controller into a booted state by closing its booted
		// channel once so users know when to go ahead.
		select {
		case <-c.booted:
		default:
			close(c.booted)
		}

		// mgr.Start() blocks the boot process until it ends gracefully or fails.
		err = mgr.Start(setupSignalHandler())
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (c *Controller) deleteFunc(ctx context.Context, obj interface{}) error {
	var err error

	hasFinalizer, err := c.hasFinalizer(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	}
	if !hasFinalizer {
		return nil
	}

	{
		ctx = reconciliationcanceledcontext.NewContext(ctx, make(chan struct{}))

		defer func() {
			ctx = unsetLoggerCtxValue(ctx, loggerKeyResource)
		}()

		for _, r := range c.resources {
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
	}

	err = c.removeFinalizer(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (c *Controller) hasPauseAnnotation(m map[string]string) (bool, string, string) {
	for k, v := range m {
		if hasAnnotation(c.pause, k, v) {
			return true, k, v
		}
	}

	return false, "", ""
}

func (c *Controller) reconcile(ctx context.Context, req reconcile.Request, obj interface{}) (reconcile.Result, error) {
	ctx, err := c.initCtx(ctx, obj)
	if err != nil {
		return reconcile.Result{}, microerror.Mask(err)
	}

	var m metav1.Object
	{
		m, err = meta.Accessor(obj)
		if err != nil {
			return reconcile.Result{}, microerror.Mask(err)
		}
	}

	ctx = setLoggerCtxValue(ctx, loggerKeyObject, m.GetSelfLink())
	ctx = setLoggerCtxValue(ctx, loggerKeyVersion, m.GetResourceVersion())

	if ok, k, v := c.hasPauseAnnotation(m.GetAnnotations()); ok {
		c.logger.Debugf(ctx, "cancelling reconciliation due to pause annotation %#q set to %#q", k, v)
		return reconcile.Result{}, nil
	}

	if m.GetDeletionTimestamp() != nil {
		event := "delete"

		t := prometheus.NewTimer(eventHistogram.WithLabelValues(event))
		ctx = setLoggerCtxValue(ctx, loggerKeyEvent, event)

		err = c.deleteFunc(ctx, obj)
		if err != nil {
			return reconcile.Result{}, microerror.Mask(err)
		}

		t.ObserveDuration()
	} else {
		event := "update"

		t := prometheus.NewTimer(eventHistogram.WithLabelValues(event))
		ctx = setLoggerCtxValue(ctx, loggerKeyEvent, event)

		err = c.updateFunc(ctx, obj)
		if err != nil {
			return reconcile.Result{}, microerror.Mask(err)
		}

		t.ObserveDuration()
	}

	return reconcile.Result{}, nil
}

func (c *Controller) updateFunc(ctx context.Context, obj interface{}) error {
	var err error

	ok, err := c.addFinalizer(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	}
	if ok {
		// A finalizer was added, this causes a new update event, so we stop
		// reconciling here and will pick up the new event.
		return nil
	}

	{
		ctx = reconciliationcanceledcontext.NewContext(ctx, make(chan struct{}))

		defer func() {
			ctx = unsetLoggerCtxValue(ctx, loggerKeyResource)
		}()

		for _, r := range c.resources {
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
	}

	return nil
}

func hasAnnotation(annotations map[string]string, targetKey string, targetValue string) bool {
	if annotations == nil {
		return false
	}

	actualValue, ok := annotations[targetKey]
	if !ok {
		return false
	}

	return actualValue == targetValue
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

func setupSignalHandler() (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

func unsetLoggerCtxValue(ctx context.Context, key string) context.Context {
	m, ok := loggermeta.FromContext(ctx)
	if !ok {
		m = loggermeta.New()
		ctx = loggermeta.NewContext(ctx, m)
		return ctx
	}

	delete(m.KeyVals, key)

	return ctx
}

func reportLastReconciled(o interface{}, controllerName string) {
	var kind string
	{
		obj, ok := o.(pkgruntime.Object)
		if !ok {
			return
		}

		gvks, _, err := scheme.Scheme.ObjectKinds(obj)
		if pkgruntime.IsNotRegisteredError(err) {
			gvks, _, err = apiextensionsscheme.Scheme.ObjectKinds(obj)
			if err != nil {
				return
			}
		} else if err != nil {
			return
		}

		if len(gvks) == 0 {
			return
		}

		kind = gvks[0].Kind
	}

	lastReconciledGauge.WithLabelValues(
		kind,
		controllerName,
	).SetToCurrentTime()
}
