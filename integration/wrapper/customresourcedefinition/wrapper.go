package customresourcedefinition

import (
	"context"
	"time"

	"github.com/giantswarm/k8sclient/v6/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/operatorkit/v5/integration/env"
	"github.com/giantswarm/operatorkit/v5/pkg/controller"
	"github.com/giantswarm/operatorkit/v5/pkg/resource"
)

type Config struct {
	Resources []resource.Interface

	Name      string
}

type Wrapper struct {
	controller *controller.Controller

	extClient clientset.Interface
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
			Logger:    newLogger,
			Resources: config.Resources,
			NewRuntimeObjectFunc: func() client.Object {
				return new(apiextensionsv1.CustomResourceDefinition)
			},
			Selector: labels.Everything(),

			Name:         config.Name,
			ResyncPeriod: 10 * time.Second,
		}

		newController, err = controller.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	wrapper := &Wrapper{
		controller: newController,
		extClient:  k8sClient.ExtClient(),
		k8sClient:  k8sClient.K8sClient(),
	}

	return wrapper, nil
}

func (w Wrapper) Controller() *controller.Controller {
	return w.controller
}

func (w Wrapper) MustSetup(ctx context.Context, namespace string) {
	w.MustTeardown(ctx, namespace)
}

func (w Wrapper) MustTeardown(ctx context.Context, namespace string) {
	err := w.extClient.ApiextensionsV1().CustomResourceDefinitions().DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{})
	if errors.IsNotFound(err) {
		// fall though
	} else if err != nil {
		panic(err)
	}
}
