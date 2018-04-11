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
	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/integration/testresourceset"
	"github.com/giantswarm/operatorkit/informer"
)

type Client struct {
	framework *framework.Framework

	k8sClient kubernetes.Interface
}

type Config struct {
	Resources []framework.Resource

	Name      string
	Namespace string
}

func New(config Config) (*Client, error) {
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
			Watcher: k8sClient.CoreV1().ConfigMaps(config.Namespace),

			RateWait:     time.Second * 2,
			ResyncPeriod: time.Second * 10,
		}
		newInformer, err = informer.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}
	var resourceSet *framework.ResourceSet
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
	var resourceRouter *framework.ResourceRouter
	{
		c := framework.ResourceRouterConfig{
			Logger: logger,

			ResourceSets: []*framework.ResourceSet{
				resourceSet,
			},
		}

		resourceRouter, err = framework.NewResourceRouter(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}
	cf := framework.Config{
		Informer:       newInformer,
		K8sClient:      k8sClient,
		Logger:         logger,
		ResourceRouter: resourceRouter,

		Name: config.Name,
	}
	f, err := framework.New(cf)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	client := &Client{
		framework: f,
		k8sClient: k8sClient,
	}
	return client, nil
}

func (c Client) Framework() *framework.Framework {
	return c.framework
}

func (c Client) MustSetup(namespace string) {
	c.MustTeardown(namespace)

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

	_, err := c.k8sClient.CoreV1().Namespaces().Create(ns)
	if err != nil {
		panic(err)
	}
}

func (c Client) MustTeardown(namespace string) {
	err := c.k8sClient.CoreV1().Namespaces().Delete(namespace, nil)
	if errors.IsNotFound(err) {
		// fall though
	} else if err != nil {
		panic(err)
	}
}
