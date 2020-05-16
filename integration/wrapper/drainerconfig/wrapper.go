package drainerconfig

import (
	"context"
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgruntime "k8s.io/apimachinery/pkg/runtime"

	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/integration/env"
	"github.com/giantswarm/operatorkit/resource"
)

type Config struct {
	Logger    micrologger.Logger
	Resources []resource.Interface

	Name      string
	Namespace string
}

type Wrapper struct {
	controller *controller.Controller

	k8sClient k8sclient.Interface
}

func New(config Config) (*Wrapper, error) {
	var err error

	if config.Logger == nil {
		c := micrologger.Config{}

		config.Logger, err = micrologger.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var k8sClient k8sclient.Interface
	{
		c := k8sclient.ClientsConfig{
			SchemeBuilder: k8sclient.SchemeBuilder{
				v1alpha1.AddToScheme,
			},
			Logger: config.Logger,

			KubeConfigPath: env.KubeConfigPath(),
		}

		k8sClient, err = k8sclient.NewClients(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var newController *controller.Controller
	{
		c := controller.Config{
			K8sClient: k8sClient,
			Logger:    config.Logger,
			Resources: config.Resources,
			NewRuntimeObjectFunc: func() pkgruntime.Object {
				return new(v1alpha1.DrainerConfig)
			},

			Name:         config.Name,
			ResyncPeriod: 2 * time.Second,
		}

		newController, err = controller.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	w := &Wrapper{
		controller: newController,
		k8sClient:  k8sClient,
	}

	return w, nil
}

func (w Wrapper) Controller() *controller.Controller {
	return w.controller
}

func (w Wrapper) MustSetup(namespace string) {
	ctx := context.Background()

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

	_, err := w.k8sClient.K8sClient().CoreV1().Namespaces().Create(ns)
	if err != nil {
		panic(err)
	}

	var backOffFactory func() backoff.Interface
	{
		backOffFactory = func() backoff.Interface {
			return backoff.NewMaxRetries(3, 1*time.Second)
		}
	}

	err = w.k8sClient.CRDClient().EnsureCreated(ctx, v1alpha1.NewDrainerConfigCRD(), backOffFactory())
	if err != nil {
		panic(err)
	}
}

func (w Wrapper) MustTeardown(namespace string) {
	ctx := context.Background()

	err := w.k8sClient.K8sClient().CoreV1().Namespaces().Delete(namespace, nil)
	if errors.IsNotFound(err) {
		// fall though
	} else if err != nil {
		panic(err)
	}

	var backOffFactory func() backoff.Interface
	{
		backOffFactory = func() backoff.Interface {
			return backoff.NewMaxRetries(3, 1*time.Second)
		}
	}

	err = w.k8sClient.CRDClient().EnsureDeleted(ctx, v1alpha1.NewDrainerConfigCRD(), backOffFactory())
	if err != nil {
		panic(err)
	}
}
