// +build k8srequired

package integration

import (
	"time"

	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/operatorkit/informer"
)

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

func newOperatorkitInformer(rateWait, resyncPeriod time.Duration) (*informer.Informer, error) {
	c := DefaultConfig()

	c.Watcher = k8sClient.CoreV1().ConfigMaps(namespace)

	c.RateWait = rateWait
	c.ResyncPeriod = resyncPeriod

	operatorkitInformer, err := New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return operatorkitInformer, nil
}

func mustSetup() {
	mustTeardown()

	k8sClient, err := newK8sClient()
	if err != nil {
		panic(err)
	}

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

	_, err = k8sClient.CoreV1().Namespaces().Create(ns)
	if err != nil {
		panic(err)
	}
}

func mustTeardown() {
	k8sClient, err := newK8sClient()
	if err != nil {
		panic(err)
	}

	err = k8sClient.CoreV1().Namespaces().Delete(namespace, nil)
	if errors.IsNotFound(err) {
		// fall though
	} else if err != nil {
		panic(err)
	}
}
