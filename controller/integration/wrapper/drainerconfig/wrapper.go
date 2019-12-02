package drainerconfig

import (
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/apiextensions/pkg/clientset/versioned"
	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/k8sclient/k8scrdclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/controller/integration/testresourceset"
	"github.com/giantswarm/operatorkit/informer"
	"github.com/giantswarm/operatorkit/resource"
)

type Config struct {
	HandlesFunc func(obj interface{}) bool
	Logger      micrologger.Logger
	Resources   []resource.Interface

	Name      string
	Namespace string
}

type Wrapper struct {
	controller *controller.Controller

	g8sClient versioned.Interface
	k8sClient kubernetes.Interface
}

func New(config Config) (*Wrapper, error) {
	restConfig, err := clientcmd.BuildConfigFromFlags("", harness.DefaultKubeConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	k8sClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	g8sClient, err := versioned.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	k8sExtClient, err := clientset.NewForConfig(restConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	if config.Logger == nil {
		c := micrologger.Config{}

		config.Logger, err = micrologger.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var crdClient *k8scrdclient.CRDClient
	{
		c := k8scrdclient.Config{
			K8sExtClient: k8sExtClient,
			Logger:       config.Logger,
		}

		crdClient, err = k8scrdclient.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var newInformer *informer.Informer
	{
		c := informer.Config{
			Logger:  config.Logger,
			Watcher: g8sClient.CoreV1alpha1().DrainerConfigs(config.Namespace),

			RateWait:     time.Second * 2,
			ResyncPeriod: time.Second * 10,
		}
		newInformer, err = informer.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var resourceSet *controller.ResourceSet
	{
		c := testresourceset.Config{
			HandlesFunc: config.HandlesFunc,
			K8sClient:   k8sClient,
			Logger:      config.Logger,
			Resources:   config.Resources,

			ProjectName: config.Name,
		}

		resourceSet, err = testresourceset.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var newController *controller.Controller
	{
		c := controller.Config{
			CRD:       v1alpha1.NewDrainerConfigCRD(),
			CRDClient: crdClient,
			Informer:  newInformer,
			Logger:    config.Logger,
			ResourceSets: []*controller.ResourceSet{
				resourceSet,
			},
			RESTClient: g8sClient.CoreV1alpha1().RESTClient(),

			Name: config.Name,
		}

		newController, err = controller.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	wrapper := &Wrapper{
		controller: newController,
		g8sClient:  g8sClient,
		k8sClient:  k8sClient,
	}

	return wrapper, nil
}

func (w Wrapper) Controller() *controller.Controller {
	return w.controller
}

func (w Wrapper) MustSetup(namespace string) {
	w.MustTeardown(namespace)

	ns := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
		Spec: corev1.NamespaceSpec{},
	}

	_, err := w.k8sClient.CoreV1().Namespaces().Create(ns)
	if err != nil {
		panic(err)
	}
}

func (w Wrapper) MustTeardown(namespace string) {
	err := w.k8sClient.CoreV1().Namespaces().Delete(namespace, nil)
	if errors.IsNotFound(err) {
		// fall though
	} else if err != nil {
		panic(err)
	}
}
