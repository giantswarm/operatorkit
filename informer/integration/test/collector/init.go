// +build k8srequired

package collector

import (
	"time"

	"github.com/giantswarm/e2e-harness/pkg/harness"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/informer"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	namespace = "test-informer-collector"
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

func newOperatorkitInformer(k8sClient kubernetes.Interface) (*informer.Informer, error) {
	logger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	operatorkitInformer, err := informer.New(informer.Config{
		Logger:  logger,
		Watcher: k8sClient.CoreV1().ConfigMaps(namespace),

		RateWait:     time.Second * 2,
		ResyncPeriod: time.Second * 10,
	})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return operatorkitInformer, nil
}

func mustSetup(k8sClient kubernetes.Interface) error {
	if err := mustTeardown(k8sClient); err != nil {
		return microerror.Mask(err)
	}

	namespace := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
		Spec: corev1.NamespaceSpec{},
	}

	if _, err := k8sClient.CoreV1().Namespaces().Create(namespace); err != nil {
		return err
	}

	return nil
}

func mustTeardown(k8sClient kubernetes.Interface) error {
	if err := k8sClient.CoreV1().Namespaces().Delete(namespace, nil); err != nil && !errors.IsNotFound(err) {
		return microerror.Mask(err)
	}

	// TODO: Wait on the namespace being actually deleted, this is a hack.
	time.Sleep(10 * time.Second)

	return nil
}
