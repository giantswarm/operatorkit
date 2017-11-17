package tpr

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cenkalti/backoff"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

const (
	// ResyncPeriod is the interval at which the Informer cache is invalidated,
	// and the lister function is called.
	ResyncPeriod = 1 * time.Minute

	tprInitMaxElapsedTime = 2 * time.Minute
)

// Config is a TPR configuration.
type Config struct {
	// Dependencies.
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger

	// Settings.

	// Description is free text description.
	Description string
	// Name takes the form <kind>.<group> (Note: The group is also called
	// a domain). You are expected to provide a unique kind and group name
	// in order to avoid conflicts with other ThirdPartyResource objects.
	// Kind names will be converted to CamelCase when creating instances of
	// the ThirdPartyResource. Hyphens in the kind are assumed to be word
	// breaks. For instance the kind camel-case would be converted to
	// CamelCase but camelcase would be converted to Camelcase.
	Name         string
	ResyncPeriod time.Duration
	// Version is TPR version, e.g. v1.
	Version string
}

// DefaultConfig provides a default configuration to create a new TPR service
// by best effort.
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
		K8sClient: nil,
		Logger:    newLogger,

		// Settings.
		Description:  "",
		Name:         "",
		ResyncPeriod: ResyncPeriod,
		Version:      "",
	}
}

// New creates a new TPR.
func New(config Config) (*TPR, error) {
	// Dependencies.
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	// Settings.
	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "name must not be empty")
	}
	if config.Version == "" {
		return nil, microerror.Maskf(invalidConfigError, "version must not be empty")
	}
	if config.Description == "" {
		return nil, microerror.Maskf(invalidConfigError, "description must not be empty")
	}

	kind, group, err := extractKindAndGroup(config.Name)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	tpr := &TPR{
		// Dependencies.
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		// Internals.
		resourceName: unsafeGuessKindToResource(kind),

		// Settings.
		name:         config.Name,
		kind:         kind,
		group:        group,
		version:      config.Version,
		apiVersion:   group + "/" + config.Version,
		description:  config.Description,
		resyncPeriod: config.ResyncPeriod,
	}

	return tpr, nil
}

// TPR allows easy operations on ThirdPartyResources. See
// https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-third-party-resource/
// for details.
type TPR struct {
	// Dependencies.
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	// Internals.

	// apiVersion is group/version.
	apiVersion   string
	group        string // see Config.Name
	kind         string // see Config.Name
	resourceName string

	// Settings.
	description  string
	resyncPeriod time.Duration
	name         string // see Config.Name
	version      string
}

// Kind returns a TPR kind extracted from Name. Useful when creating
// ThirdPartyObjects. See Config.Name godoc for details.
func (t *TPR) Kind() string {
	return t.kind
}

// APIVersion returns a TPR APIVersion created from group and version. It takes
// format <group>/<version>. Useful for creating ThirdPartyObjects. See
// Config.Name and Config.Version for details.
func (t *TPR) APIVersion() string {
	return t.apiVersion
}

// Name returns a TPR name provided with Config.Name.
func (t *TPR) Name() string {
	return t.name
}

// Group returns a TPR group extracted from Name. See Config.Name godoc for details.
func (t *TPR) Group() string {
	return t.group
}

// Endpoint returns a TPR resource endpoint registered in the Kubernetes API
// under a given namespace. The default namespace will be used when the
// argument is an empty string.
func (t *TPR) Endpoint(namespace string) string {
	nsResource := t.resourceName
	if len(namespace) != 0 {
		nsResource = "namespaces/" + namespace + "/" + t.resourceName
	}
	return "/apis/" + t.group + "/" + t.version + "/" + nsResource
}

// WatchEndpoint returns a TPR watch resource endpoint registered in the
// Kubernetes API under a given namespace. The default namespace will be used
// when the argument is an empty string.
func (t *TPR) WatchEndpoint(namespace string) string {
	nsResource := "watch/" + t.resourceName
	if len(namespace) != 0 {
		nsResource = "namespaces/" + namespace + "/watch/" + t.resourceName
	}
	return "/apis/" + t.group + "/" + t.version + "/" + nsResource
}

// CreateAndWait creates a TPR and waits till it is initialized in the cluster.
// Returns alreadyExistsError when the resource already exists.
func (t *TPR) CreateAndWait() error {
	initBackOff := backoff.NewExponentialBackOff()
	initBackOff.MaxElapsedTime = tprInitMaxElapsedTime
	return t.CreateAndWaitBackOff(initBackOff)
}

// CreateAndWaitBackOff creates a TPR and waits till it is initialized in the
// cluster. It allows passing a custom initialization back off policy used to
// poll for TPR readiness. Returns alreadyExistsError when the resource already
// exists.
func (t *TPR) CreateAndWaitBackOff(initBackOff backoff.BackOff) error {
	err := t.create(initBackOff)
	if err != nil {
		return microerror.Maskf(err, "creating TPR %s", t.name)
	}
	err = t.waitInit(initBackOff)
	if err != nil {
		return microerror.Maskf(err, "waiting for TPR %s initialization", t.name)
	}
	return nil
}

func (t *TPR) NewInformer(resourceEventHandler cache.ResourceEventHandler, zeroObjectFactory ZeroObjectFactory) cache.Controller {
	listWatch := &cache.ListWatch{
		ListFunc: func(options apismetav1.ListOptions) (runtime.Object, error) {
			t.logger.Log("debug", "executing the reconciler's list function", "event", "list")

			req := t.k8sClient.Core().RESTClient().Get().AbsPath(t.Endpoint(""))
			b, err := req.DoRaw()
			if err != nil {
				return nil, microerror.Mask(err)
			}

			v := zeroObjectFactory.NewObjectList()
			if err := json.Unmarshal(b, v); err != nil {
				return nil, microerror.Mask(err)
			}

			return v, nil
		},
		WatchFunc: func(options apismetav1.ListOptions) (watch.Interface, error) {
			t.logger.Log("debug", "executing the reconciler's watch function", "event", "watch")

			req := t.k8sClient.CoreV1().RESTClient().Get().AbsPath(t.WatchEndpoint(""))
			stream, err := req.Stream()
			if err != nil {
				return nil, microerror.Mask(err)
			}

			watcher := watch.NewStreamWatcher(newDecoder(stream, zeroObjectFactory))
			return watcher, nil
		},
	}

	_, informer := cache.NewInformer(listWatch, zeroObjectFactory.NewObject(), t.resyncPeriod, resourceEventHandler)

	return informer
}

// create is extracted for testing because fake REST client does not work.
// Therefore waitInit can not be tested.
func (t *TPR) create(retry backoff.BackOff) error {
	tpr := &v1beta1.ThirdPartyResource{
		ObjectMeta: apismetav1.ObjectMeta{
			Name: t.name,
		},
		Versions: []v1beta1.APIVersion{
			{Name: t.version},
		},
		Description: t.description,
	}

	createTpr := func() error {
		_, err := t.k8sClient.ExtensionsV1beta1().ThirdPartyResources().Create(tpr)
		if err != nil && apierrors.IsAlreadyExists(err) {
			return backoff.Permanent(microerror.Mask(alreadyExistsError))
		}
		if err != nil {
			return microerror.Maskf(err, "creating TPR %s", t.name)
		}

		return nil
	}

	err := backoff.Retry(createTpr, retry)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (t *TPR) waitInit(retry backoff.BackOff) error {
	endpoint := t.Endpoint("")
	op := func() error {
		_, err := t.k8sClient.CoreV1().RESTClient().Get().RequestURI(endpoint).DoRaw()
		return err
	}

	err := backoff.Retry(op, retry)

	if apierrors.IsNotFound(err) {
		err = tprInitTimeoutError
	}
	return microerror.Maskf(err, "requesting TPR %s", t.name)
}

func (t *TPR) CollectMetrics(ctx context.Context) {
	go func() {
		t.logger.Log("info", "starting metrics collection")

		ticker := time.NewTicker(t.resyncPeriod)

		for {
			select {
			case <-ctx.Done():
				t.logger.Log("info", "context done, stopping metrics collection")
				return

			case <-ticker.C:
				t.logger.Log("info", "listing TPOs for metrics")

				operation := func() error {
					req := t.k8sClient.Core().RESTClient().Get().AbsPath(t.Endpoint(""))
					b, err := req.DoRaw()
					if err != nil {
						return microerror.Mask(err)
					}

					list := TPOList{}
					if err := json.Unmarshal(b, &list); err != nil {
						return microerror.Mask(err)
					}

					tpoCount.WithLabelValues(t.Kind(), t.APIVersion(), t.Name(), t.Group()).Set(float64(len(list.Items)))

					return nil
				}

				if err := backoff.Retry(operation, backoff.NewExponentialBackOff()); err != nil {
					t.logger.Log("error", "could not get tpo metrics", "message", err.Error())
				}
			}
		}
	}()
}
