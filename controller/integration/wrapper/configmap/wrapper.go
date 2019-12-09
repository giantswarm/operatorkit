package configmap

import (
	"time"

	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/controller/integration/testresourceset"
	"github.com/giantswarm/operatorkit/resource"
)

type Config struct {
	Resources []resource.Interface

	Name      string
	Namespace string
}

type Wrapper struct {
	controller *controller.Controller

	k8sClient kubernetes.Interface
}

func New(config Config) (*Wrapper, error) {
	var err error

	var newLogger micrologger.Logger
	{
		c := micrologger.Config{}

		newLogger, err = micrologger.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var k8sClient *k8sclient.Clients
	{
		c := k8sclient.ClientsConfig{
			Logger: newLogger,

			KubeConfigPath: harness.DefaultKubeConfig,
		}

		k8sClient, err = k8sclient.NewClients(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var resourceSet *controller.ResourceSet
	{
		c := testresourceset.Config{
			K8sClient: k8sClient.K8sClient(),
			Logger:    newLogger,
			Resources: config.Resources,

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
			K8sClient: k8sClient,
			Logger:    newLogger,
			ResourceSets: []*controller.ResourceSet{
				resourceSet,
			},
			NewRuntimeObjectFunc: func() pkgruntime.Object {
				return new(corev1.ConfigMap)
			},

			Name:         config.Name,
			ResyncPeriod: 2 * time.Second,
		}

		newController, err = controller.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	wrapper := &Wrapper{
		controller: newController,
		k8sClient:  k8sClient.K8sClient(),
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
