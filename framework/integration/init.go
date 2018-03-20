// +build k8srequired

package integration

import (
	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/integration/test_resource_set"
	"github.com/giantswarm/operatorkit/informer"
)

const (
	namespace = "finalizer-integration-test"
)

var (
	err error

	k8sClient     kubernetes.Interface
	k8sRESTClient rest.Interface
)

func init() {
	k8sClient, err = newK8sClient()
	if err != nil {
		panic(err)
	}
	k8sRESTClient, err = newK8sRESTClient()
	if err != nil {
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

func newK8sRESTClient() (rest.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", harness.DefaultKubeConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	k8sRESTClient, err := rest.RESTClientFor(config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return k8sRESTClient, nil
}

func newFramework(name string) (*framework.Framework, error) {
	logger, err := micrologger.New(micrologger.DefaultConfig())
	if err != nil {
		return nil, err
	}
	var newInformer *informer.Informer
	{
		c := informer.Config{
			Watcher: k8sClient.CoreV1().Pods(namespace),
		}
		newInformer, err = informer.New(c)
		if err != nil {
			return nil, err
		}
	}
	var resourceSet *framework.ResourceSet
	{
		c := testresourceset.ResourceSetConfig{
			K8sClient: k8sClient,
			Logger:    logger,

			ProjectName: name,
		}

		resourceSet, err = testresourceset.NewResourceSet(c)
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
		Logger:         logger,
		RestClient:     k8sRESTClient,
		Name:           name,
		ResourceRouter: resourceRouter,
	}
	f, err := framework.New(cf)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func mustSetup() {
	mustTeardown()

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

func mustTeardown() {
	err := k8sClient.CoreV1().Namespaces().Delete(namespace, nil)
	if errors.IsNotFound(err) {
		// fall though
	} else if err != nil {
		panic(err)
	}
}
