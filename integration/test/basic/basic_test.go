// +build k8srequired

package basic

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/giantswarm/operatorkit/integration/testresource"
	"github.com/giantswarm/operatorkit/integration/wrapper/configmap"
	"github.com/giantswarm/operatorkit/resource"
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
	var err error

	expectedFinalizers := []string{
		testFinalizer,
	}

	var r *testresource.Resource
	{
		c := testresource.Config{
			Name: "test-resource",
		}

		r, err = testresource.New(c)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	var wrapper *configmap.Wrapper
	{
		c := configmap.Config{
			Resources: []resource.Interface{r},
			Name:      operatorName,
			Namespace: testNamespace,
		}

		wrapper, err = configmap.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	wrapper.MustSetup(testNamespace)
	defer wrapper.MustTeardown(testNamespace)

	controller := wrapper.Controller()

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
	_, err = wrapper.CreateObject(testNamespace, cm)
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
	_, err = wrapper.UpdateObject(testNamespace, cm)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We reconcile the ConfigMap using its name and namespace.
	// This is expected to only add one finalizer, we want to make sure that we
	// only use the latest ResourceVersion of an object.
	_, err = controller.Reconcile(reconcile.Request{NamespacedName: types.NamespacedName{
		Name:      cm.GetName(),
		Namespace: cm.GetNamespace(),
	}})
	if err != nil {
		t.Fatal("failed reconciliation", nil, "got", err)
	}

	// We run Reconcile multiple times to check for duplicates.
	_, err = controller.Reconcile(reconcile.Request{NamespacedName: types.NamespacedName{
		Name:      cm.GetName(),
		Namespace: cm.GetNamespace(),
	}})
	if err != nil {
		t.Fatal("failed reconciliation", nil, "got", err)
	}

	// We get the current configmap.
	resultObj, err := wrapper.GetObject(configMapName, testNamespace)
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
	err = wrapper.DeleteObject(configMapName, testNamespace)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	resultObj, err = wrapper.GetObject(configMapName, testNamespace)
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
	// We run Reconcile multiple times to check for duplicates.
	_, _ = controller.Reconcile(reconcile.Request{NamespacedName: types.NamespacedName{
		Name:      cm.GetName(),
		Namespace: cm.GetNamespace(),
	}})

	// We verify that our object is completely gone now.
	_, err = wrapper.GetObject(configMapName, testNamespace)
	if !configmap.IsNotFound(err) {
		t.Fatalf("error == %#v, want NotFound error", err)
	}
}
