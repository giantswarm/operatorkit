// +build k8srequired

package basic

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/operatorkit/framework/integration/client"
	"github.com/giantswarm/operatorkit/framework/integration/client/configmap"
)

// Test_Finalizer_Integration_Basic is a integration test for basic finalizer
// operations. The test verifies that finalizers are added and removed as
// expected. It does not cover correct behavior with reconciliation.
//
// !!! This test does not work with CRs, the framework is not booted !!!
//
func Test_Finalizer_Integration_Basic(t *testing.T) {
	configMapName := "test-cm"
	expectedFinalizers := []string{
		"operatorkit.giantswarm.io/test-operator",
	}
	operatorName := "test-operator"
	testNamespace := "finalizer-integration-basic-test"

	c := client.Config{
		Name:      operatorName,
		Namespace: testNamespace,
	}
	configmapClient, err := configmap.New(c)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	testClient := client.Interface(configmapClient)

	testClient.MustSetup(testNamespace)
	defer testClient.MustTeardown(testNamespace)

	operatorkitFramework := testClient.GetFramework()

	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: testNamespace,
		},
		Data: map[string]string{},
	}
	// We create an object which does not have any finalizers.
	createdObj, err := testClient.CreateObject(testNamespace, cm)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We directly pass the object to UpdateFunc.
	operatorkitFramework.UpdateFunc(createdObj, createdObj)

	resultObj, err := testClient.GetObject(configMapName, testNamespace)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	resultObjAccessor, err := meta.Accessor(resultObj)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We verify that the correct finalizer has been set during UpdateFunc.
	if !reflect.DeepEqual(resultObjAccessor.GetFinalizers(), expectedFinalizers) {
		t.Fatalf("finalizers == %v, want %v", resultObjAccessor.GetFinalizers(), expectedFinalizers)
	}

	// We delete our object.
	err = testClient.DeleteObject(configMapName, testNamespace)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	resultObj, err = testClient.GetObject(configMapName, testNamespace)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	resultObjAccessor, err = meta.Accessor(resultObj)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We verify, that our object still exists, but has a DeletionTimestamp set.
	if resultObjAccessor.GetDeletionTimestamp() == nil {
		t.Fatalf("DeletionTimestamp == nil, want non-nil")
	}

	// We verify, that our finalizer is still set.
	if !reflect.DeepEqual(resultObjAccessor.GetFinalizers(), expectedFinalizers) {
		t.Fatalf("finalizers == %v, want %v", resultObjAccessor.GetFinalizers(), expectedFinalizers)
	}

	// We directly pass the object to DeleteFunc to remove the finalizer.
	operatorkitFramework.DeleteFunc(resultObj)

	// We verify that our object is completely gone now.
	_, err = testClient.GetObject(configMapName, testNamespace)
	if !errors.IsNotFound(err) {
		t.Fatalf("error == %#v, want NotFound error", err)
	}

}
