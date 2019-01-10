// +build k8srequired

package parallel

import (
	"context"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/controller/integration/testresource"
	"github.com/giantswarm/operatorkit/controller/integration/wrapper/nodeconfig"
)

const (
	objName      = "test-obj"
	objNamespace = "integration-parallel-test"

	controllerNameA = "test-controller-a"
	finalizerA      = "operatorkit.giantswarm.io/test-controller-a"
	resourceNameA   = "test-resource-a"

	controllerNameB = "test-controller-b"
	finalizerB      = "operatorkit.giantswarm.io/test-controller-b"
	resourceNameB   = "test-resource-b"

	controllerNameC = "test-controller-c"
	resourceNameC   = "test-resource-c"
	finalizerC      = "operatorkit.giantswarm.io/test-controller-c"
)

// Test_Finalizer_Integration_Parallel is a integration test to
// check that multiple controllers can function in parallel.
func Test_Finalizer_Integration_Parallel(t *testing.T) {
	var err error

	ctx := context.Background()

	var resourceA *testresource.Resource
	{
		c := testresource.Config{
			Name: resourceNameA,
		}

		resourceA, err = testresource.New(c)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// resourceB errors once during the deletion.
	var resourceB *testresource.Resource
	{
		var once sync.Once

		c := testresource.Config{
			Name: resourceNameB,
			ReturnErrorFunc: func(obj interface{}) error {
				// Do not return error for update events.
				{
					accessor, err := meta.Accessor(obj)
					if err != nil {
						return microerror.Mask(err)
					}
					if accessor.GetDeletionTimestamp() == nil {
						return nil
					}
				}

				// Return error for first deletion event.
				{
					var err error
					once.Do(func() {
						err = microerror.Maskf(testError, "I fail once during the deletion")
					})

					if err != nil {
						return microerror.Mask(err)
					}
				}

				return nil
			},
		}

		resourceB, err = testresource.New(c)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// resourceC fails all the time.
	var resourceC *testresource.Resource
	{
		c := testresource.Config{
			Name: resourceNameC,
			ReturnErrorFunc: func(obj interface{}) error {
				return microerror.Maskf(testError, "I fail all the time")
			},
		}

		resourceC, err = testresource.New(c)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	var harnessA, harnessB, harnessC *nodeconfig.Wrapper
	{
		harnessA, err = newHarness(objNamespace, controllerNameA, resourceA)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
		harnessB, err = newHarness(objNamespace, controllerNameB, resourceB)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
		harnessC, err = newHarness(objNamespace, controllerNameC, resourceC)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Start controllers.
	{
		controllerA := harnessA.Controller()
		controllerB := harnessB.Controller()
		controllerC := harnessC.Controller()

		go controllerA.Boot(ctx)
		go controllerB.Boot(ctx)
		go controllerC.Boot(ctx)
		select {
		case <-controllerA.Booted():
		case <-time.After(30 * time.Second):
			t.Fatalf("failed to wait for controllerA to boot")
		}
		select {
		case <-controllerB.Booted():
		case <-time.After(30 * time.Second):
			t.Fatalf("failed to wait for controllerB to boot")
		}
		select {
		case <-controllerC.Booted():
		case <-time.After(30 * time.Second):
			t.Fatalf("failed to wait for controllerC to boot")
		}
	}

	// Setup the test namespace. We use the harness A. It makes no
	// difference if we use the harness A or B.
	{
		harnessA.MustSetup(objNamespace)
		defer harnessA.MustTeardown(objNamespace)
	}

	// Create an object. Creation is retried because the CRD might still
	// not be ensured.
	{
		o := func() error {
			nodeConfig := &v1alpha1.NodeConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      objName,
					Namespace: objNamespace,
				},
			}

			_, err := harnessA.CreateObject(objNamespace, nodeConfig)
			if err != nil {
				return microerror.Mask(err)
			}
			return nil
		}
		b := backoff.NewMaxRetries(20, 1*time.Second)

		err = backoff.Retry(o, b)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Verify controllers reconcile creation of that object. There should
	// be 2 ResyncPeriods in 30 seconds so verify there were more than
	// 2 create events.
	//
	// 		EnsureCreated: >2, EnsureDeleted: =0
	//
	{
		o := func() error {
			if resourceA.CreateCount() <= 2 {
				return microerror.Maskf(waitError, "resourceA.CreateCount() == %v, want more than %v", resourceA.CreateCount(), 2)
			}
			if resourceA.DeleteCount() != 0 {
				return microerror.Maskf(waitError, "resourceA.DeleteCount() == %v, want %v", resourceA.DeleteCount(), 0)
			}
			if resourceB.CreateCount() <= 2 {
				return microerror.Maskf(waitError, "resourceB.CreateCount() == %v, want more than %v", resourceB.CreateCount(), 2)
			}
			if resourceB.DeleteCount() != 0 {
				return microerror.Maskf(waitError, "resourceB.DeleteCount() == %v, want %v", resourceB.DeleteCount(), 0)
			}
			if resourceC.CreateCount() <= 2 {
				return microerror.Maskf(waitError, "resourceC.CreateCount() == %v, want more than %v", resourceC.CreateCount(), 2)
			}
			if resourceC.DeleteCount() != 0 {
				return microerror.Maskf(waitError, "resourceC.DeleteCount() == %v, want %v", resourceC.DeleteCount(), 0)
			}

			return nil
		}
		b := backoff.NewMaxRetries(30, 1*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Verify deletion timestamp and finalizers.
	{
		obj, err := harnessA.GetObject(objName, objNamespace)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}

		accessor, err := meta.Accessor(obj)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}

		if accessor.GetDeletionTimestamp() != nil {
			t.Fatalf("DeletionTimestamp == %v, want %v", accessor.GetDeletionTimestamp(), nil)
		}

		finalizers := accessor.GetFinalizers()
		expectedFinalizers := []string{
			finalizerA,
			finalizerB,
			finalizerC,
		}
		sort.Strings(finalizers)
		sort.Strings(expectedFinalizers)
		if !reflect.DeepEqual(finalizers, expectedFinalizers) {
			t.Fatalf("finalizers == %v, want %v", finalizers, expectedFinalizers)
		}
	}

	// Delete the object.
	{
		err := harnessA.DeleteObject(objName, objNamespace)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Verify deletion timestamp and check there is only finzalierC left
	// from constantly failing resource.
	{
		o := func() error {
			obj, err := harnessA.GetObject(objName, objNamespace)
			if err != nil {
				t.Fatalf("err == %v, want %v", err, nil)
			}

			accessor, err := meta.Accessor(obj)
			if err != nil {
				t.Fatalf("err == %v, want %v", err, nil)
			}

			if accessor.GetDeletionTimestamp() == nil {
				t.Fatalf("DeletionTimestamp == %v, want non nil", accessor.GetDeletionTimestamp())
			}

			finalizers := accessor.GetFinalizers()
			expectedFinalizers := []string{
				finalizerC,
			}
			sort.Strings(finalizers)
			sort.Strings(expectedFinalizers)
			if !reflect.DeepEqual(finalizers, expectedFinalizers) {
				return microerror.Maskf(waitError, "finalizers == %v, want %v", finalizers, expectedFinalizers)
			}

			return nil
		}
		b := backoff.NewMaxRetries(30, 1*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Verify that:
	//
	//	- resourceA received exactly 1 deletion event as it never
	//	  fails.
	//	- resourceB received exectly 2 deletion events as it fails
	//	  during the deletion once.
	//	- resourceB received more than 3 deletion events as it fails
	//	  all the time.
	//
	{
		o := func() error {
			if resourceA.DeleteCount() != 1 {
				microerror.Maskf(waitError, "resourceA.DeleteCount() == %v, want %v", resourceA.DeleteCount(), 1)
			}
			if resourceB.DeleteCount() != 2 {
				microerror.Maskf(waitError, "resourceB.DeleteCount() == %v, want %v", resourceB.DeleteCount(), 2)
			}
			if resourceC.DeleteCount() <= 3 {
				microerror.Maskf(waitError, "resourceC.DeleteCount() == %v, want more than %v", resourceC.DeleteCount(), 3)
			}

			return nil
		}

		b := backoff.NewMaxRetries(30, 1*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Set the resourceC error function to nil to not return any error
	// anymore. The finalizer should be removed with the next
	// reconciliation now.
	{
		resourceC.SetReturnErrorFunc(nil)
	}

	// Verify that the object is completely gone now.
	{
		o := func() error {
			_, err := harnessA.GetObject(objName, objNamespace)
			if nodeconfig.IsNotFound(err) {
				return nil
			} else if err != nil {
				return microerror.Mask(err)
			}

			return microerror.Maskf(waitError, "object %#q in namespace %#q still exists", objName, objNamespace)
		}
		b := backoff.NewMaxRetries(30, 1*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatalf("failed to wait for object deletion: %#v", err)
		}
	}

	// Verify **again** that:
	//
	//	- resourceA received exactly 1 deletion event as it never
	//	  fails.
	//	- resourceB received exectly 2 deletion events as it fails
	//	  during the deletion once.
	//	- resourceB received more than 3 deletion events as it was
	//	  failing all the time till SetReturnErrorFunc(nil) was called.
	//
	{
		if resourceA.DeleteCount() != 1 {
			t.Fatalf("resourceA.DeleteCount() == %v, want %v", resourceA.DeleteCount(), 1)
		}
		if resourceB.DeleteCount() != 2 {
			t.Fatalf("resourceB.DeleteCount() == %v, want %v", resourceB.DeleteCount(), 2)
		}
		if resourceC.DeleteCount() <= 3 {
			t.Fatalf("resourceC.DeleteCount() == %v, want more than %v", resourceC.DeleteCount(), 3)
		}
	}
}

func newHarness(namespace string, controllerName string, resource *testresource.Resource) (*nodeconfig.Wrapper, error) {
	resources := []controller.Resource{
		controller.Resource(resource),
	}

	c := nodeconfig.Config{
		Resources: resources,

		Name:      controllerName,
		Namespace: namespace,
	}

	harness, err := nodeconfig.New(c)
	if err != nil {
		return nil, err
	}

	return harness, nil
}
