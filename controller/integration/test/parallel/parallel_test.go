// +build k8srequired

package parallel

import (
	"testing"
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/controller/integration/testresource"
	"github.com/giantswarm/operatorkit/controller/integration/wrapper"
	"github.com/giantswarm/operatorkit/controller/integration/wrapper/nodeconfig"
)

const (
	objName       = "test-obj"
	testNamespace = "integration-parallel-test"

	controllerNameA = "test-controller-a"
	testFinalizerA  = "operatorkit.giantswarm.io/test-controller-a"

	controllerNameB = "test-controller-b"
	testFinalizerB  = "operatorkit.giantswarm.io/test-controller-b"
)

// Test_Finalizer_Integration_Parallel is a integration test to
// check that multiple controllers can function in parallel.
func Test_Finalizer_Integration_Parallel(t *testing.T) {
	var err error

	// We create the first resource "A" here with its own resource.
	var trA *testresource.Resource
	{
		c := testresource.Config{}

		trA, err = testresource.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	testWrapperA, err := setupController(testNamespace, controllerNameA, trA)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	controllerA := testWrapperA.Controller()

	// We create the second resource "B" and give it a different resource.
	var trB *testresource.Resource
	{
		c := testresource.Config{}

		trB, err = testresource.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	testWrapperB, err := setupController(testNamespace, controllerNameB, trB)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	controllerB := testWrapperB.Controller()

	// We setup the namespace in which we test. We use the wrapper of controller A
	// here, it makes no difference if we use the wrapper of A or B.
	testWrapperA.MustSetup(testNamespace)
	defer testWrapperA.MustTeardown(testNamespace)

	// We start the controllers.
	go controllerA.Boot()
	go controllerB.Boot()
	<-controllerA.Booted()
	<-controllerB.Booted()

	// We create an object without any finalizers.
	// Creation is retried because the existance of a CRD might have to be ensured.
	{
		o := func() error {
			nodeConfig := &v1alpha1.NodeConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      objName,
					Namespace: testNamespace,
				},
			}

			_, err := testWrapperA.CreateObject(testNamespace, nodeConfig)
			if err != nil {
				return microerror.Mask(err)
			}
			return nil
		}
		b := backoff.NewExponential(2*time.Minute, 10*time.Second)
		err = backoff.Retry(o, b)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// We use backoff with the absolute maximum amount:
	// 10 second ResyncPeriod + 2 second RateWait + 8 second for safety.
	// The controllers should now add their finalizers and EnsureCreated should be
	// hit once immediatly.
	//
	// 		EnsureCreated: 1, EnsureDeleted: 0
	//
	// The controllers should reconcile once in this period.
	//
	// 		EnsureCreated: 2, EnsureDeleted: 0
	//
	operation := func() error {
		// We are more forgiving here compared to other tests, the controllers might
		// receive events at different times. Checking the count exactly might fail
		// if a controller is slower and the other one reconciles one more time.
		if trA.CreateCount() < 2 {
			return microerror.Maskf(countMismatchError, "EnsureCreated of controller A was hit %v times, want atleast %v", trA.CreateCount(), 2)
		}
		if trA.DeleteCount() != 0 {
			return microerror.Maskf(countMismatchError, "EnsureDeleted of controller A was hit %v times, want %v", trA.DeleteCount(), 0)
		}
		if trB.CreateCount() < 2 {
			return microerror.Maskf(countMismatchError, "EnsureCreated of controller B was hit %v times, want atleast %v", trB.CreateCount(), 2)
		}
		if trB.DeleteCount() != 0 {
			return microerror.Maskf(countMismatchError, "EnsureDeleted of controller B was hit %v times, want %v", trB.DeleteCount(), 0)
		}
		return nil
	}
	err = backoff.Retry(operation, backoff.NewMaxRetries(20, 1*time.Second))
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
	// We check this individually because we are not sure in which order the
	// finalizers were added.
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

	// Verify that our resources received single delete events.
	{
		o := func() error {
			if trA.DeleteCount() != 1 {
				return microerror.Maskf(waitError, "EnsureDeleted of controller A was hit %v times, want %v", trA.DeleteCount(), 1)
			}
			if trB.DeleteCount() != 1 {
				return microerror.Maskf(waitError, "EnsureDeleted of controller B was hit %v times, want %v", trB.DeleteCount(), 1)
			}

			return nil
		}
		b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)
		backoff.Retry(o, b)
		err := backoff.RetryNotify(o, b, n)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	// Verify that our object is completely gone now.
	{
		o := func() error {
			_, err := testWrapperA.GetObject(objName, testNamespace)
			if !nodeconfig.IsNotFound(err) {
				return microerror.Maskf(waitError, "object %#q in namespace %#q still exists", objName, testNamespace)
			}

			return nil
		}
		b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)
		err := backoff.Retry(o, b)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	// Veryfy resources didn't receive more deletion events after
	// successful deleteion.
	{
		o := func() error {
			if trA.DeleteCount() != 1 {
				return microerror.Maskf(waitError, "EnsureDeleted of controller A was hit %v times, want %v", trA.DeleteCount(), 1)
			}
			if trB.DeleteCount() != 1 {
				return microerror.Maskf(waitError, "EnsureDeleted of controller B was hit %v times, want %v", trB.DeleteCount(), 1)
			}

			return nil
		}
		b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)
		backoff.Retry(o, b)
		err := backoff.RetryNotify(o, b, n)
		if err != nil {
			return microerror.Mask(err)
		}
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
