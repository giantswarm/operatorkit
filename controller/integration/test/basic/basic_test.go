// +build k8srequired

package basic

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/operatorkit/controller/integration/wrapper"
	"github.com/giantswarm/operatorkit/controller/integration/wrapper/configmap"
)

const (
	configMapName = "test-cm"
	operatorName  = "test-operator"
	testFinalizer = "operatorkit.giantswarm.io/test-operator"
	testNamespace = "finalizer-integration-basic-test"
)

// Test_Finalizer_Integration_Basic is a integration test for basic finalizer
// operations. The test verifies that finalizers are added and removed as
// expected. It does not cover correct behavior with reconciliation.
//
// !!! This test does not work with CRs, the controller is not booted !!!
//
func Test_Finalizer_Integration_Basic(t *testing.T) {
	expectedFinalizers := []string{
		testFinalizer,
	}

	c := configmap.Config{
		Name:      operatorName,
		Namespace: testNamespace,
	}
	configmapWrapper, err := configmap.New(c)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	testWrapper := wrapper.Interface(configmapWrapper)

	testWrapper.MustSetup(testNamespace)
	defer testWrapper.MustTeardown(testNamespace)

	operatorkitController := testWrapper.Controller()

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
	createdObj, err := testWrapper.CreateObject(testNamespace, cm)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We update the object with a meaningless label to ensure a change in the
	// ResourceVersion of the ConfigMap.
	cm.SetLabels(
		map[string]string{
			"testlabel": "testlabel",
		},
	)
	_, err = testWrapper.UpdateObject(testNamespace, cm)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We directly pass the _old_ configmap to UpdateFunc.
	// This is expected to only add one finalizer, we want to make sure that we
	// only use the latest ResourceVersion of an object.
	operatorkitController.UpdateFunc(nil, createdObj)

	// We run UpdateFunc multiple times on the old object to check for duplicates.
	operatorkitController.UpdateFunc(nil, createdObj)

	// We get the current configmap.
	resultObj, err := testWrapper.GetObject(configMapName, testNamespace)
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
	err = testWrapper.DeleteObject(configMapName, testNamespace)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	resultObj, err = testWrapper.GetObject(configMapName, testNamespace)
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
	operatorkitController.DeleteFunc(resultObj)

	// We verify that our object is completely gone now.
	_, err = testWrapper.GetObject(configMapName, testNamespace)
	if !configmap.IsNotFound(err) {
		t.Fatalf("error == %#v, want NotFound error", err)
	}
}
