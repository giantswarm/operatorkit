// +build k8srequired

package reconciliation

import (
	"reflect"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/integration/client"
	"github.com/giantswarm/operatorkit/framework/integration/testresource"
)

// Test_Finalizer_Integration_Reconciliation is a integration test for
// the proper replay and reconciliation of delete events with finalizers.
func Test_Finalizer_Integration_Reconciliation(t *testing.T) {
	configMapName := "test-cm"
	testFinalizer := "operatorkit.giantswarm.io/test-operator"
	testNamespace := "finalizer-integration-reconciliation-test"
	testOtherFinalizer := "operatorkit.giantswarm.io/other-operator"
	operatorName := "test-operator"

	client.MustSetup(testNamespace)
	defer client.MustTeardown(testNamespace)
	var err error
	var tr *testresource.Resource
	{
		c := testresource.Config{}

		tr, err = testresource.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	resources := []framework.Resource{
		framework.Resource(tr),
	}

	c := client.Config{
		Resources: resources,

		Name:      operatorName,
		Namespace: testNamespace,
	}

	operatorkitFramework, err := client.NewFramework(c)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We create a configmap, but add a finalizer of another operator. This will
	// cause the configmap to continue existing after the framework removes it own
	// finalizer.
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: testNamespace,
			Finalizers: []string{
				testOtherFinalizer,
			},
		},
		Data: map[string]string{},
	}
	_, err = client.CreateConfigMap(cm)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We start the framework.
	go operatorkitFramework.Boot()

	// We wait the absolute maximum amount of time here:
	// 20 second ResyncPeriod + 2 second RateWait + 3 second for safety.
	// The framework should now add the finalizer and EnsureCreated should be hit
	// once immediatly.
	//
	// 		EnsureCreated: 1, EnsureDeleted: 0
	//
	// The framework should reconcile twice in this period.
	//
	// 		EnsureCreated: 3, EnsureDeleted: 0
	//
	time.Sleep(25 * time.Second)

	// We get the ConfigMap after the framework has been started.
	resultConfigMap, err := client.GetConfigMap(configMapName, testNamespace)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We verify, that the DeletionTimestamp has not been set.
	if resultConfigMap.GetDeletionTimestamp() != nil {
		t.Fatalf("DeletionTimestamp != nil, want nil")
	}

	// We define which finalizers we currently expect.
	expectedFinalizers := []string{
		testOtherFinalizer,
		testFinalizer,
	}

	// We verify, that our finalizer is still set.
	if !reflect.DeepEqual(resultConfigMap.GetFinalizers(), expectedFinalizers) {
		t.Fatalf("finalizers == %v, want %v", resultConfigMap.GetFinalizers(), expectedFinalizers)
	}

	// Verify that we hit the resource functions for the expected amounts.
	if tr.GetCreateCount() != 3 {
		t.Fatalf("EnsureCreated was hit %v times, want %v", tr.GetCreateCount(), 3)
	}

	if tr.GetDeleteCount() != 0 {
		t.Fatalf("EnsureDeleted was hit %v times, want %v", tr.GetDeleteCount(), 0)
	}

	// We delete the configmap now.
	err = client.DeleteConfigMap(configMapName, testNamespace)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We wait the absolute maximum amount of time here:
	// 20 second ResyncPeriod + 2 second RateWait + 3 second for safety.
	// The framework should now remove the finalizer and EnsureDeleted should be
	// hit twice immediatly. See https://github.com/giantswarm/giantswarm/issues/2897
	//
	// 		EnsureCreated: 3, EnsureDeleted: 2
	//
	// The framework should also reconcile twice in this period. (The other
	// finalizer is still set, so we reconcile.)
	//
	// 		EnsureCreated: 3, EnsureDeleted: 4
	//
	time.Sleep(25 * time.Second)

	// We get the ConfigMap after the framework has handled the deletion event.
	resultConfigMap, err = client.GetConfigMap(configMapName, testNamespace)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We verify, that our configmap still exists, but has a DeletionTimestamp set.
	if resultConfigMap.GetDeletionTimestamp() == nil {
		t.Fatalf("DeletionTimestamp == nil, want non-nil")
	}

	// We define which finalizers we currently expect.
	expectedFinalizers = []string{
		testOtherFinalizer,
	}

	// We verify, that our finalizer is still set.
	if !reflect.DeepEqual(resultConfigMap.GetFinalizers(), expectedFinalizers) {
		t.Fatalf("finalizers == %v, want %v", resultConfigMap.GetFinalizers(), expectedFinalizers)
	}

	// Verify that we hit the resource functions for the expected amounts.
	if tr.GetCreateCount() != 3 {
		t.Fatalf("EnsureCreated was hit %v times, want %v", tr.GetCreateCount(), 3)
	}

	if tr.GetDeleteCount() != 4 {
		t.Fatalf("EnsureDeleted was hit %v times, want %v", tr.GetDeleteCount(), 4)
	}

}
