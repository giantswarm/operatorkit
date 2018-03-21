// +build k8srequired

package integration

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_Finalizer_Integration_Basic(t *testing.T) {
	mustSetup()
	defer mustTeardown()
	operatorName := "test-operator"
	configMapName := "test-cm"
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
	err = createConfigMap(cm)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	operatorkitFramework.UpdateFunc(cm, cm)

	resultConfigMap, err := getConfigMap(configMapName)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	expectedFinalizers := []string{
		"operatorkit.giantswarm.io/test-operator",
	}

	if !reflect.DeepEqual(resultConfigMap.GetFinalizers(), expectedFinalizers) {
		t.Fatalf("finalizers == %v, want %v", resultConfigMap.GetFinalizers(), expectedFinalizers)
	}

}
