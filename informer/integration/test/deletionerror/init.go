// +build k8srequired

package deletionerror

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

	"github.com/giantswarm/operatorkit/informer"
)

const (
	namespace = "test-informer-integration-deletionerror"
)

var (
	err error

	k8sClient kubernetes.Interface
)

func init() {
	k8sClient, err = newK8sClient()
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

func newOperatorkitInformerAndWatcher(rateWait, resyncPeriod time.Duration) (*informer.Informer, *FilterWatcher, error) {
	var err error

	var filterWatcher *FilterWatcher
	{
		c := FilterWatcherConfig{
			Watcher: k8sClient.CoreV1().ConfigMaps(namespace),
		}

		filterWatcher, err = NewFilterWatcher(c)
		if err != nil {
			return nil, nil, microerror.Mask(err)
		}
	}

	var newLogger micrologger.Logger
	{
		c := micrologger.Config{}

		newLogger, err = micrologger.New(c)
		if err != nil {
			return nil, nil, microerror.Mask(err)
		}
	}

	var operatorkitInformer *informer.Informer
	{
		c := informer.Config{
			Logger:  newLogger,
			Watcher: filterWatcher,

			RateWait:     rateWait,
			ResyncPeriod: resyncPeriod,
		}

		operatorkitInformer, err = informer.New(c)
		if err != nil {
			return nil, nil, microerror.Mask(err)
		}
	}

	return operatorkitInformer, filterWatcher, nil
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
