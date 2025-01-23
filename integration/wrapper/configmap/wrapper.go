package configmap

import (
	"context"
	"strings"
	"time"

	"github.com/giantswarm/k8sclient/v8/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

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
			Namespace: config.Namespace,
			NewRuntimeObjectFunc: func() client.Object {
				return new(corev1.ConfigMap)
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
		k8sClient:  k8sClient.K8sClient(),
	}

	return wrapper, nil
}

func (w Wrapper) Controller() *controller.Controller {
	return w.controller
}

func (w Wrapper) MustSetup(ctx context.Context, namespace string) {
	w.MustTeardown(ctx, namespace)

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

	_, err := w.k8sClient.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func (w Wrapper) MustTeardown(ctx context.Context, namespace string) {
	_, err := w.k8sClient.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return
	} else if err != nil {
		panic(err)
	}

	configMaps, err := w.k8sClient.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	if errors.IsNotFound(err) {
		// fall though
	} else if err != nil {
		panic(err)
	}

	for _, configMap := range configMaps.Items {
		var filteredFinalizers []string
		for _, finalizer := range configMap.Finalizers {
			if !strings.HasPrefix(finalizer, "operatorkit.giantswarm.io") {
				filteredFinalizers = append(filteredFinalizers, finalizer)
			}
		}
		if len(filteredFinalizers) != len(configMap.Finalizers) {
			configMap.Finalizers = filteredFinalizers
			_, err = w.k8sClient.CoreV1().ConfigMaps(namespace).Update(ctx, &configMap, metav1.UpdateOptions{}) //nolint:gosec
			if errors.IsNotFound(err) {
				// fall though
			} else if err != nil {
				panic(err)
			}
		}
	}

	err = w.k8sClient.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		// fall though
	} else if err != nil {
		panic(err)
	}
}

func (w Wrapper) Events(ctx context.Context, namespace string) ([]corev1.Event, error) {
	events, err := w.k8sClient.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return events.Items, nil
}
