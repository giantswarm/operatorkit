// +build k8srequired

package reconciliation

import (
	"reflect"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/integration/client"
	"github.com/giantswarm/operatorkit/framework/integration/client/nodeconfig"
	"github.com/giantswarm/operatorkit/framework/integration/testresource"
)

// Test_Finalizer_Integration_Reconciliation is a integration test for
// the proper replay and reconciliation of delete events with finalizers.
func Test_Finalizer_Integration_Reconciliation(t *testing.T) {
	objName := "test-obj"
	testFinalizer := "operatorkit.giantswarm.io/test-operator"
	testNamespace := "finalizer-integration-reconciliation-test"
	testOtherFinalizer := "operatorkit.giantswarm.io/other-operator"
	operatorName := "test-operator"

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

	nodeconfigClient, err := nodeconfig.New(c)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	testClient := client.Interface(nodeconfigClient)

	testClient.MustSetup(testNamespace)
	defer testClient.MustTeardown(testNamespace)

	operatorkitFramework := testClient.GetFramework()

	// We start the framework.
	go operatorkitFramework.Boot()

	// We create an object, but add a finalizer of another operator. This will
	// cause the object to continue existing after the framework removes its own
	// finalizer.
	//
	//Creation is retried because the existance of a CRD might have to be ensured.
	obj := &v1alpha1.NodeConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      objName,
			Namespace: testNamespace,
			Finalizers: []string{
				testOtherFinalizer,
			},
		},
	}
	operation := func() error {
		_, err = testClient.CreateObject(testNamespace, obj)
		if err != nil {
			return microerror.Mask(err)
		}
		return nil
	}
	err = backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

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

	// We get the object after the framework has been started.
	resultObj, err := testClient.GetObject(objName, testNamespace)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	resultObjAccessor, err := meta.Accessor(resultObj)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We verify, that the DeletionTimestamp has not been set.
	if resultObjAccessor.GetDeletionTimestamp() != nil {
		t.Fatalf("DeletionTimestamp != nil, want nil")
	}

	// We define which finalizers we currently expect.
	expectedFinalizers := []string{
		testOtherFinalizer,
		testFinalizer,
	}

	// We verify, that our finalizer is still set.
	if !reflect.DeepEqual(resultObjAccessor.GetFinalizers(), expectedFinalizers) {
		t.Fatalf("finalizers == %v, want %v", resultObjAccessor.GetFinalizers(), expectedFinalizers)
	}

	// Verify that we hit the resource functions for the expected amounts.
	if tr.GetCreateCount() != 3 {
		t.Fatalf("EnsureCreated was hit %v times, want %v", tr.GetCreateCount(), 3)
	}

	if tr.GetDeleteCount() != 0 {
		t.Fatalf("EnsureDeleted was hit %v times, want %v", tr.GetDeleteCount(), 0)
	}

	// We delete the object now.
	err = testClient.DeleteObject(objName, testNamespace)
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

	// We get the object after the framework has handled the deletion event.
	resultObj, err = testClient.GetObject(objName, testNamespace)
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

	// We define which finalizers we currently expect.
	expectedFinalizers = []string{
		testOtherFinalizer,
	}

	// We verify, that our finalizer is still set.
	if !reflect.DeepEqual(resultObjAccessor.GetFinalizers(), expectedFinalizers) {
		t.Fatalf("finalizers == %v, want %v", resultObjAccessor.GetFinalizers(), expectedFinalizers)
	}

	// Verify that we hit the resource functions for the expected amounts.
	if tr.GetCreateCount() != 3 {
		t.Fatalf("EnsureCreated was hit %v times, want %v", tr.GetCreateCount(), 3)
	}

	if tr.GetDeleteCount() != 4 {
		t.Fatalf("EnsureDeleted was hit %v times, want %v", tr.GetDeleteCount(), 4)
	}

}
