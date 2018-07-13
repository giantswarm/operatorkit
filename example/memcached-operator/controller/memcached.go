package controller

import (
	examplev1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/example/v1alpha1"
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/microerror"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/operatorkit/client/k8scrdclient"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/example/memcached-operator/controller/resource"
	"github.com/giantswarm/operatorkit/example/memcached-operator/logger"
	"github.com/giantswarm/operatorkit/informer"
)

const name = "memcached-operator"

// Config represents the configuration used to create a new memcached controller.
type Config struct {
	K8sClient    kubernetes.Interface
	K8sExtClient apiextensionsclient.Interface
	G8sClient    versioned.Interface
}

// Memcached is a type containing the OperatorKit controller.
type Memcached struct {
	*controller.Controller
}

// New creates a new memcached controller.
func NewMemcached(config Config) (*Memcached, error) {
	var err error

	var (
		crd        = examplev1alpha1.NewMemcachedConfigCRD()
		restClient = config.G8sClient.ExampleV1alpha1().RESTClient()
		watcher    = config.G8sClient.ExampleV1alpha1().MemcachedConfigs("")
	)

	var deploymentsResource controller.Resource
	{
		c := resource.DeploymentsConfig{
			K8sClient: config.K8sClient,
		}

		deploymentsResource, err = resource.NewDeployments(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var servicesResource controller.Resource
	{
		c := resource.ServicesConfig{}

		servicesResource, err = resource.NewServices(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []controller.Resource{
		deploymentsResource,
		servicesResource,
	}

	// Below is a common controller wiring. This code doesn't change unless
	// you want to reconcile non-CRD objects or to use more sophisticated
	// object routing.

	// crdClient ensures that the configured CRD exists on the cluster when the
	// operator boots.
	var crdClient *k8scrdclient.CRDClient
	{
		c := k8scrdclient.Config{
			Logger: logger.Default,

			K8sExtClient: config.K8sExtClient,
		}

		crdClient, err = k8scrdclient.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

	}

	// resourceSets is used in more complex operators that support multiple
	// resources and multiple versions of those resources.
	resourceSets, err := newSimpleResourceSets(resources)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// memcachedInformer implements a list watch of the memcachedconfig custom
	// resources.
	var memcachedInformer *informer.Informer
	{
		c := informer.Config{
			Logger: logger.Default,

			Watcher: watcher,
		}

		memcachedInformer, err = informer.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	// underlying is the OperatorKit controller. It implements a control loop
	// that reconciles the current state of the resources towards their
	// desired state.
	var underlying *controller.Controller
	{
		c := controller.Config{
			CRD:          crd,
			CRDClient:    crdClient,
			Informer:     memcachedInformer,
			Logger:       logger.Default,
			ResourceSets: resourceSets,
			RESTClient:   restClient,

			Name: name,
		}

		underlying, err = controller.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	m := &Memcached{
		Controller: underlying,
	}

	return m, nil
}
