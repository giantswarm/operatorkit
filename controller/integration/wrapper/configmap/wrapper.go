// +build k8srequired

package configmap

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/controller/integration/testresourceset"
	"github.com/giantswarm/operatorkit/informer"
)

type Config struct {
	Resources []controller.Resource

	Name      string
	Namespace string
}

type Wrapper struct {
	controller *controller.Controller

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

	logger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		return nil, microerror.Mask(err)
	}
	var newInformer *informer.Informer
	{
		c := informer.Config{
			Logger:  logger,
			Watcher: k8sClient.CoreV1().ConfigMaps(config.Namespace),

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
			K8sClient: k8sClient,
			Logger:    logger,
			Resources: config.Resources,

			ProjectName: config.Name,
		}

		resourceSet, err = testresourceset.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}
	var resourceRouter *controller.ResourceRouter
	{
		c := controller.ResourceRouterConfig{
			Logger: logger,

			ResourceSets: []*controller.ResourceSet{
				resourceSet,
			},
		}

		resourceRouter, err = controller.NewResourceRouter(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}
	cf := controller.Config{
		Informer:       newInformer,
		RESTClient:     k8sClient.CoreV1().RESTClient(),
		Logger:         logger,
		ResourceRouter: resourceRouter,

		Name: config.Name,
	}
	f, err := controller.New(cf)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	wrapper := &Wrapper{
		controller: f,
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
