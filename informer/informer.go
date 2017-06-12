package informer

import (
	"encoding/json"
	"time"

	"github.com/giantswarm/operatorkit/tpr"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

// Config represents the configuration used to create a new informer.
type Config struct {
	// Dependencies.
	K8sClient            kubernetes.Interface
	Logger               micrologger.Logger
	Observer             Observer
	ResourceEventHandler cache.ResourceEventHandler
	TPR                  *tpr.TPR
	ZeroObjectFactory    ZeroObjectFactory

	// Settings.
	ResyncPeriod time.Duration
}

// DefaultConfig provides a default configuration to create a new create service
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		K8sClient:            nil,
		Logger:               nil,
		Observer:             nil,
		ResourceEventHandler: nil,
		TPR:                  nil,
		ZeroObjectFactory:    nil,

		// Settings.
		ResyncPeriod: tpr.ResyncPeriod,
	}
}

// New returns a new configured informer.
func New(config Config) (*cache.Controller, error) {
	// Dependencies.
	if config.K8sClient == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.ResourceEventHandler == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.ResourceEventHandler must not be empty")
	}
	if config.TPR == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.TPR must not be empty")
	}
	if config.ZeroObjectFactory == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.ZeroObjectFactory must not be empty")
	}

	if config.Observer == nil {
		config.Observer = ObserverFuncs{
			OnListFunc: func() {
				config.Logger.Log("debug", "executing the reconciler's list function", "event", "list")
			},
			OnWatchFunc: func() {
				config.Logger.Log("debug", "executing the reconciler's watch function", "event", "watch")
			},
		}
	}

	err := config.TPR.CreateAndWait()
	if tpr.IsAlreadyExists(err) {
		config.Logger.Log("debug", "third party resource already exists")
	} else if err != nil {
		return nil, microerror.MaskAny(err)
	}
	config.Logger.Log("debug", "successfully created third party resource")

	listWatch := &cache.ListWatch{
		ListFunc: func(options api.ListOptions) (runtime.Object, error) {
			config.Observer.OnList()

			req := config.K8sClient.Core().RESTClient().Get().AbsPath(config.TPR.Endpoint(""))
			b, err := req.DoRaw()
			if err != nil {
				return nil, microerror.MaskAny(err)
			}

			v := config.ZeroObjectFactory.NewObjectList()
			if err := json.Unmarshal(b, v); err != nil {
				return nil, microerror.MaskAny(err)
			}

			return v, nil
		},
		WatchFunc: func(options api.ListOptions) (watch.Interface, error) {
			config.Observer.OnWatch()

			req := config.K8sClient.CoreV1().RESTClient().Get().AbsPath(config.TPR.WatchEndpoint(""))
			stream, err := req.Stream()
			if err != nil {
				return nil, microerror.MaskAny(err)
			}

			watcher := watch.NewStreamWatcher(&decoder{
				stream: stream,
				obj:    config.ZeroObjectFactory,
			})
			return watcher, nil
		},
	}

	_, informer := cache.NewInformer(listWatch, config.ZeroObjectFactory.NewObject(), config.ResyncPeriod, config.ResourceEventHandler)

	return informer, nil
}
