// +build k8srequired

package controlflow

import (
	"context"
	"reflect"
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
	"github.com/giantswarm/operatorkit/controller/integration/wrapper/drainerconfig"
)

const (
	objName      = "test-obj"
	objNamespace = "integration-controlflow-test"

	controllerName = "test-controller"
	finalizer      = "operatorkit.giantswarm.io/test-controller"
	resourceName   = "test-resource"
)

// Test_Finalizer_Integration_Controlflow is an integration test to check that
// errors during deletion prevent the finalizer from removal.
func Test_Finalizer_Integration_Controlflow(t *testing.T) {
	var err error

	ctx := context.Background()

	var resource *testresource.Resource
	{
		c := testresource.Config{
			Name: resourceName,
		}

		resource, err = testresource.New(c)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	var harness wrapper.Interface
	{
		harness, err = newHarness(objNamespace, controllerName, resource)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Start controller.
	{
		controller := harness.Controller()

		go controller.Boot(ctx)
		select {
		case <-controller.Booted():
		case <-time.After(30 * time.Second):
			t.Fatalf("failed to wait for controller to boot")
		}
	}

	// Setup the test namespace.
	{
		harness.MustSetup(objNamespace)
		defer harness.MustTeardown(objNamespace)
	}

	// Create an object and wait for the controller to add a finalizer.
	// Creation is retried because the CRD might still not be ensured.
	{
		o := func() error {
			drainerConfig := &v1alpha1.DrainerConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      objName,
					Namespace: objNamespace,
				},
			}

			_, err := harness.CreateObject(objNamespace, drainerConfig)
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

	// Verify the controller reconciles creation of that object. There
	// should be 2 ResyncPeriods in 30 seconds so verify there were more
	// than 2 create events.
	//
	// 		EnsureCreated: >2, EnsureDeleted: =0
	//
	{
		o := func() error {
			if resource.CreateCount() <= 2 {
				return microerror.Maskf(waitError, "resource.CreateCount() == %v, want more than %v", resource.CreateCount(), 2)
			}
			if resource.DeleteCount() != 0 {
				return microerror.Maskf(waitError, "resource.DeleteCount() == %v, want %v", resource.DeleteCount(), 0)
			}

			return nil
		}
		b := backoff.NewMaxRetries(30, 1*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Verify deletion timestamp and finalizer.
	{
		obj, err := harness.GetObject(objName, objNamespace)
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
			finalizer,
		}
		if !reflect.DeepEqual(finalizers, expectedFinalizers) {
			t.Fatalf("finalizers == %v, want %v", finalizers, expectedFinalizers)
		}
	}

	// Set an error function to return an error. This makes the resource
	// always return an error and should therefore prevent the removal of
	// the finalizer.
	{
		resource.SetReturnErrorFunc(func(obj interface{}) error {
			return microerror.Mask(testError)
		})
	}

	// Delete the object.
	{
		err := harness.DeleteObject(objName, objNamespace)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Verify the controller reconciles deletion the object. There should
	// be 2 ResyncPeriods in 30 seconds so verify there were more than
	// 2 delete events because of there error the resource returns.
	//
	// 		EnsureCreated: >2, EnsureDeleted: >2
	//
	{
		o := func() error {
			if resource.CreateCount() <= 2 {
				return microerror.Maskf(waitError, "resource.CreateCount() == %v, want more than %v", resource.CreateCount(), 2)
			}
			if resource.DeleteCount() <= 2 {
				return microerror.Maskf(waitError, "resource.DeleteCount() == %v, want more than %v", resource.DeleteCount(), 2)
			}

			return nil
		}
		b := backoff.NewMaxRetries(30, 1*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Verify deletion timestamp and finalizer.
	{
		obj, err := harness.GetObject(objName, objNamespace)
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
			finalizer,
		}
		if !reflect.DeepEqual(finalizers, expectedFinalizers) {
			t.Fatalf("finalizers == %v, want %v", finalizers, expectedFinalizers)
		}
	}

	// Set the error function to nil to not return any error anymore. The
	// finalizer should be removed with the next reconciliation now.
	{
		resource.SetReturnErrorFunc(nil)
	}

	// We verify that the object is completely gone now.
	{
		o := func() error {
			_, err = harness.GetObject(objName, objNamespace)
			if drainerconfig.IsNotFound(err) {
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
}

func newHarness(namespace string, controllerName string, resource *testresource.Resource) (*drainerconfig.Wrapper, error) {
	resources := []controller.Resource{
		controller.Resource(resource),
	}

	c := drainerconfig.Config{
		Resources: resources,

		Name:      controllerName,
		Namespace: namespace,
	}

	harness, err := drainerconfig.New(c)
	if err != nil {
		return nil, err
	}

	return harness, nil
}
