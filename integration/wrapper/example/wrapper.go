package example

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/giantswarm/k8sclient/v8/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	v1 "github.com/giantswarm/operatorkit/v7/api/v1"
	"github.com/giantswarm/operatorkit/v7/integration/env"
	"github.com/giantswarm/operatorkit/v7/pkg/controller"
	"github.com/giantswarm/operatorkit/v7/pkg/resource"
)

type Config struct {
	Resources []resource.Interface

	Name      string
	Namespace string
}

type Wrapper struct {
	controller *controller.Controller

	ctrlClient client.Client
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
			Namespace: config.Namespace,
			NewRuntimeObjectFunc: func() client.Object {
				return new(v1.Example)
			},
			Selector:     labels.Everything(),
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
		ctrlClient: k8sClient.CtrlClient(),
	}

	return wrapper, nil
}

func (w Wrapper) Controller() *controller.Controller {
	return w.controller
}

func (w Wrapper) MustSetup(ctx context.Context, namespace string) {
	w.MustTeardown(ctx, namespace)

	ns := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	err := w.ctrlClient.Create(ctx, &ns)
	if errors.IsAlreadyExists(err) {
		// fall though
	} else if err != nil {
		panic(err)
	}

	crdYAML, err := os.ReadFile("../../../config/crd/testing.giantswarm.io_examples.yaml")
	if err != nil {
		panic(microerror.JSON(err))
	}

	var crd apiextensionsv1.CustomResourceDefinition
	err = yaml.Unmarshal(crdYAML, &crd)
	if err != nil {
		panic(microerror.JSON(err))
	}

	err = w.ctrlClient.Create(ctx, &crd)
	if errors.IsAlreadyExists(err) {
		// fall though
	} else if err != nil {
		panic(err)
	}
}

func (w Wrapper) MustTeardown(ctx context.Context, namespace string) {
	var crd apiextensionsv1.CustomResourceDefinition
	err := w.ctrlClient.Get(ctx, client.ObjectKey{
		Name: "examples.testing.giantswarm.io",
	}, &crd)
	if errors.IsNotFound(err) {
		// fall though
	} else if err != nil {
		panic(err)
	} else {
		var existing v1.ExampleList
		err = w.ctrlClient.List(ctx, &existing)
		if err != nil {
			panic(err)
		}

		for _, item := range existing.Items {
			var filteredFinalizers []string
			for _, finalizer := range item.Finalizers {
				if !strings.HasPrefix(finalizer, "operatorkit.giantswarm.io") {
					filteredFinalizers = append(filteredFinalizers, finalizer)
				}
			}
			if len(filteredFinalizers) != len(item.Finalizers) {
				item.Finalizers = filteredFinalizers
				err = w.ctrlClient.Update(ctx, &item) //nolint:gosec
				if errors.IsNotFound(err) {
					// fall though
				} else if err != nil {
					panic(err)
				}
			}
		}

		err = w.ctrlClient.Delete(ctx, &crd)
		if errors.IsNotFound(err) {
			// fall though
		} else if err != nil {
			panic(err)
		}
	}

	ns := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	err = w.ctrlClient.Delete(ctx, &ns)
	if errors.IsNotFound(err) {
		// fall though
	} else if err != nil {
		panic(err)
	}
}
