// +build k8srequired

package client

import (
	"time"

	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/integration/testresourceset"
	"github.com/giantswarm/operatorkit/informer"
)

var (
	err error

	k8sClient kubernetes.Interface
)

type Config struct {
	Informer  *informer.Informer
	Resources []framework.Resource

	Name      string
	Namespace string
}

func init() {
	k8sClient, err = newK8sClient()
	if err != nil {
		panic(err)
	}
}

func NewFramework(config Config) (*framework.Framework, error) {
	logger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		return nil, microerror.Mask(err)
	}
	var newInformer *informer.Informer
	{

		if config.Informer == nil {
			c := informer.Config{
				Watcher: k8sClient.CoreV1().ConfigMaps(config.Namespace),

				RateWait:     time.Second * 2,
				ResyncPeriod: time.Second * 10,
			}
			newInformer, err = informer.New(c)
			if err != nil {
				return nil, microerror.Mask(err)
			}
		} else {
			newInformer = config.Informer
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
	return f, nil
}

func MustSetup(namespace string) {
	MustTeardown(namespace)

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

	_, err := k8sClient.CoreV1().Namespaces().Create(ns)
	if err != nil {
		panic(err)
	}
}

func MustTeardown(namespace string) {
	err := k8sClient.CoreV1().Namespaces().Delete(namespace, nil)
	if errors.IsNotFound(err) {
		// fall though
	} else if err != nil {
		panic(err)
	}
}

func newK8sClient() (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", harness.DefaultKubeConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return k8sClient, nil
}
