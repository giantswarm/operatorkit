// +build k8srequired

package parallel

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/controller/integration/testresource"
	"github.com/giantswarm/operatorkit/controller/integration/wrapper"
	"github.com/giantswarm/operatorkit/controller/integration/wrapper/nodeconfig"
)

// Test_Finalizer_Integration_Parallel is a integration test to
// check that multiple controllers can function in parallel.
func Test_Finalizer_Integration_Parallel(t *testing.T) {
	var err error
	objName := "test-obj"
	testNamespace := "finalizer-integration-parallel-test"

	testFinalizerA := "operatorkit.giantswarm.io/test-operator-a"
	operatorNameA := "test-operator-a"

	testFinalizerB := "operatorkit.giantswarm.io/test-operator-b"
	operatorNameB := "test-operator-b"

	//
	var trA *testresource.Resource
	{
		c := testresource.Config{}

		trA, err = testresource.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	testWrapperA, err := setupController(testNamespace, operatorNameA, trA)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	controllerA := testWrapperA.Controller()

	//
	var trB *testresource.Resource
	{
		c := testresource.Config{}

		trB, err = testresource.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	testWrapperB, err := setupController(testNamespace, operatorNameB, trB)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	controllerB := testWrapperB.Controller()

	//
	testWrapperA.MustSetup(testNamespace)
	defer testWrapperA.MustTeardown(testNamespace)

	// We start the controller.
	go controllerA.Boot()
	go controllerB.Boot()

	// We create an object, but add a finalizer of another operator. This will
	// cause the object to continue existing after the controller removes its own
	// finalizer.
	//
	//Creation is retried because the existance of a CRD might have to be ensured.
	obj := &v1alpha1.NodeConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      objName,
			Namespace: testNamespace,
		},
	}
	var createdObj interface{}
	operation := func() error {
		createdObj, err = testWrapperA.CreateObject(testNamespace, obj)
		if err != nil {
			return microerror.Mask(err)
		}
		return nil
	}
	err = backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We use backoff with the absolute maximum amount:
	// 10 second ResyncPeriod + 2 second RateWait + 8 second for safety.
	// The controller should now add the finalizer and EnsureCreated should be hit
	// once immediatly.
	//
	// 		EnsureCreated: 1, EnsureDeleted: 0
	//
	// The controller should reconcile once in this period.
	//
	// 		EnsureCreated: 2, EnsureDeleted: 0
	//
	operation = func() error {
		if trA.CreateCount() < 2 {
			return microerror.Maskf(countMismatchError, "EnsureCreated of controller A was hit %v times, want atleast %v", trA.CreateCount(), 2)
		}
		if trA.DeleteCount() != 0 {
			return microerror.Maskf(countMismatchError, "EnsureDeleted of controller A was hit %v times, want %v", trB.DeleteCount(), 0)
		}
		if trB.CreateCount() < 2 {
			return microerror.Maskf(countMismatchError, "EnsureCreated of controller B was hit %v times, want atleast %v", trB.CreateCount(), 2)
		}
		if trB.DeleteCount() != 0 {
			return microerror.Maskf(countMismatchError, "EnsureDeleted of controller B was hit %v times, want %v", trB.DeleteCount(), 0)
		}
		return nil
	}
	err = backoff.Retry(operation, newConstantBackoff(uint64(20)))
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We get the object after the controllers have been started.
	resultObj, err := testWrapperA.GetObject(objName, testNamespace)
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

	// We verify, that our finalizer is still set.
	if !containsFinalizer(resultObjAccessor.GetFinalizers(), testFinalizerA) {
		t.Fatalf("finalizers == %v, want to contain %v", resultObjAccessor.GetFinalizers(), testFinalizerA)
	}
	if !containsFinalizer(resultObjAccessor.GetFinalizers(), testFinalizerB) {
		t.Fatalf("finalizers == %v, want to contain %v", resultObjAccessor.GetFinalizers(), testFinalizerB)
	}

	// We delete the object now.
	err = testWrapperA.DeleteObject(objName, testNamespace)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We use backoff with the absolute maximum amount:
	// 10 second ResyncPeriod + 2 second RateWait + 8 second for safety.
	// The controller should get the deletion event immediatly but not remove the
	// finalizer because of the error we return in our resource.
	//
	// 		EnsureCreated: 2, EnsureDeleted: 1
	//
	// The controller should also reconcile once in this period. (The other
	// finalizer is still set, so we reconcile.)
	//
	// 		EnsureCreated: 2, EnsureDeleted: 2
	//
	operation = func() error {
		if trA.DeleteCount() < 2 {
			return microerror.Maskf(countMismatchError, "EnsureDeleted of controller A was hit %v times, want atleast %v", trA.DeleteCount(), 2)
		}
		if trB.DeleteCount() < 2 {
			return microerror.Maskf(countMismatchError, "EnsureDeleted of controller B was hit %v times, want atleast %v", trB.DeleteCount(), 2)
		}
		return nil
	}
	err = backoff.Retry(operation, newConstantBackoff(uint64(20)))
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We verify that our object is completely gone now.
	_, err = testWrapperA.GetObject(objName, testNamespace)
	if !errors.IsNotFound(err) {
		t.Fatalf("error == %#v, want NotFound error", err)
	}
}

func containsFinalizer(finalizers []string, finalizer string) bool {
	for _, f := range finalizers {
		if f == finalizer {
			return true
		}
	}
	return false
}

func setupController(namespace string, operatorName string, resource *testresource.Resource) (wrapper.Interface, error) {
	resources := []controller.Resource{
		controller.Resource(resource),
	}

	c := nodeconfig.Config{
		Resources: resources,

		Name:      operatorName,
		Namespace: namespace,
	}

	nodeconfigWrapper, err := nodeconfig.New(c)
	if err != nil {
		return nil, err
	}

	testWrapper := wrapper.Interface(nodeconfigWrapper)

	return testWrapper, nil
}
