// +build k8srequired

package integration

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Test_Finalizer_Integration_Basic is a integration test for basic finalizer
// operations. The test verifies that finalizers are added and removed as
// expected. It does not cover correct behavior with reconciliation.
func Test_Finalizer_Integration_Basic(t *testing.T) {
	mustSetup()
	defer mustTeardown()

	operatorName := "test-operator"
	configMapName := "test-cm"
	expectedFinalizers := []string{
		"operatorkit.giantswarm.io/test-operator",
	}
	operatorkitFramework, err := newFramework(operatorName)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

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
	createdConfigMap, err := createConfigMap(cm)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We directly pass the configmap to UpdateFunc.
	operatorkitFramework.UpdateFunc(createdConfigMap, createdConfigMap)

	resultConfigMap, err := getConfigMap(configMapName)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We verify that the correct finalizer has been set during UpdateFunc.
	if !reflect.DeepEqual(resultConfigMap.GetFinalizers(), expectedFinalizers) {
		t.Fatalf("finalizers == %v, want %v", resultConfigMap.GetFinalizers(), expectedFinalizers)
	}

	// We delete out configmap.
	err = deleteConfigMap(configMapName)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	resultConfigMap, err = getConfigMap(configMapName)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We verify, that our configmap still exists, but has a DeletionTimestamp set.
	if resultConfigMap.GetDeletionTimestamp() == nil {
		t.Fatalf("DeletionTimestamp == nil, want non-nil")
	}

	// We directly pass the configmap to DeleteFunc to remove the finalizer.
	operatorkitFramework.DeleteFunc(resultConfigMap)

	// We verify that our configmap is completely gone now.
	_, err = getConfigMap(configMapName)
	if !errors.IsNotFound(err) {
		t.Fatalf("error == %#v, want NotFound error", err)
	}

}
