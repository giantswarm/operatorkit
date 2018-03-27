// +build k8srequired

package integration

import (
	"testing"
	"time"

	"github.com/giantswarm/operatorkit/informer"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Test_Finalizer_Integration_Reconciliation is a integration test for basic finalizer
// operations. The test verifies that finalizers are added and removed as
// expected. It does not cover correct behavior with reconciliation.
func Test_Finalizer_Integration_Reconciliation(t *testing.T) {
	namespace := "finalizer-integration-reconciliation-test"
	mustSetup(namespace)
	defer mustTeardown(namespace)

	operatorName := "test-operator"
	configMapName := "test-cm"

	c := informer.Config{
		Watcher: k8sClient.CoreV1().ConfigMaps(namespace),

		RateWait:     time.Second * 2,
		ResyncPeriod: time.Second * 10,
	}
	operatorkitInformer, err := informer.New(c)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	operatorkitFramework, err := newFramework(operatorName, namespace, operatorkitInformer)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	operatorkitFramework.Boot()

	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
		},
		Data: map[string]string{},
	}

	// We create a configmap which does not have any finalizers.
	_, err = createConfigMap(namespace, cm)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	deleteConfigMap(namespace, configMapName)

}
