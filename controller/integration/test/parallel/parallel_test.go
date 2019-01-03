// +build k8srequired

package parallel

import (
	"reflect"
	"sort"
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
)

// Test_Finalizer_Integration_Parallel is a integration test to
// check that multiple controllers can function in parallel.
func Test_Finalizer_Integration_Parallel(t *testing.T) {
	var err error

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

	var resourceB *testresource.Resource
	{
		c := testresource.Config{
			Name: resourceNameB,
		}

		resourceB, err = testresource.New(c)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	var harnessA, harnessB *nodeconfig.Wrapper
	{
		harnessA, err = newHarness(objNamespace, controllerNameA, resourceA)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
		harnessB, err = newHarness(objNamespace, controllerNameB, resourceB)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Start controllers.
	{
		controllerA := harnessA.Controller()
		controllerB := harnessB.Controller()

		go controllerA.Boot()
		go controllerB.Boot()
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
	}

	// Setup the test namespace. We use the harness A. It makes no
	// difference if we use the harness A or B.
	{
		harnessA.MustSetup(objNamespace)
		defer harnessA.MustTeardown(objNamespace)
	}

	// We create an object without any finalizers.
	// Creation is retried because the existance of a CRD might have to be ensured.
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
		b := backoff.NewExponential(2*time.Minute, 10*time.Second)
		err = backoff.Retry(o, b)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Verify we reconcile creation of that object. We should have also
	// 2 ResyncPeriods in 30 seconds so we check if there were more than
	// 2 create events.
	{
		o := func() error {
			if resourceA.CreateCount() < 2 {
				return microerror.Maskf(waitError, "resourceA.CreateCount() == %v, want more than %v", resourceA.CreateCount(), 2)
			}
			if resourceA.DeleteCount() != 0 {
				return microerror.Maskf(waitError, "resourceA.DeleteCount() == %v, want %v", resourceA.DeleteCount(), 0)
			}
			if resourceB.CreateCount() < 2 {
				return microerror.Maskf(waitError, "resourceB.CreateCount() == %v, want more than %v", resourceB.CreateCount(), 2)
			}
			if resourceB.DeleteCount() != 0 {
				return microerror.Maskf(waitError, "resourceB.DeleteCount() == %v, want %v", resourceB.DeleteCount(), 0)
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
		sort.Strings(finalizers)
		expectedFinalizers := []string{
			finalizerA,
			finalizerB,
		}
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

	// Verify resources received exactly one deletion event after
	// successful deletion.
	{
		if resourceA.DeleteCount() > 0 {
			t.Fatalf("resourceA.DeleteCount() == %v, want more than %v", resourceA.DeleteCount(), 0)
		}
		if resourceB.DeleteCount() > 0 {
			t.Fatalf("resourceB.DeleteCount() == %v, want more than %v", resourceB.DeleteCount(), 0)
		}
	}
}

func newHarness(namespace string, operatorName string, resource *testresource.Resource) (*nodeconfig.Wrapper, error) {
	resources := []controller.Resource{
		controller.Resource(resource),
	}

	c := nodeconfig.Config{
		Resources: resources,

		Name:      operatorName,
		Namespace: namespace,
	}

	harness, err := nodeconfig.New(c)
	if err != nil {
		return nil, err
	}

	return harness, nil
}
