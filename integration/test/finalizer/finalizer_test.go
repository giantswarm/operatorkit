// +build k8srequired

package parallel

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/operatorkit/integration/testresource"
	"github.com/giantswarm/operatorkit/integration/wrapper/drainerconfig"
	"github.com/giantswarm/operatorkit/resource"
)

const (
	// controllerName is the same for all three controllers in order to simulate
	// multiple VOO deployments of the same operator.
	controllerName = "test-controller"
	// finalizer is the only finalizer shared by the three test controllers.
	finalizer    = "operatorkit.giantswarm.io/test-controller"
	objName      = "test-obj"
	objNamespace = "integration-finalizer-test"
)

// Test_Controller_Integration_Finalizer is an integration test to check that
// finalizer management works as expected.
func Test_Controller_Integration_Finalizer(t *testing.T) {
	var err error

	ctx := context.Background()

	// r fails all the time so we can replay the delete event until the end of the
	// test.
	var r *testresource.Resource
	{
		c := testresource.Config{
			Name: "test-resource",
			ReturnErrorFunc: func(obj interface{}) error {
				return microerror.Maskf(testError, "I fail all the time")
			},
		}

		r, err = testresource.New(c)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	var wrapper *drainerconfig.Wrapper
	{
		wrapper, err = newWrapper(r, newWrapperLogger("a"))
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Start controllers.
	{
		controllerA := wrapper.Controller()

		go controllerA.Boot(ctx)
		select {
		case <-controllerA.Booted():
		case <-time.After(30 * time.Second):
			t.Fatalf("failed to wait for controllerA to boot")
		}
	}

	// Setup the test namespace.
	{
		wrapper.MustSetup(objNamespace)
		defer wrapper.MustTeardown(objNamespace)
	}

	// Create an object. Creation is retried because the CRD might still not be
	// ensured.
	{
		o := func() error {
			drainerConfig := &v1alpha1.DrainerConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      objName,
					Namespace: objNamespace,
				},
			}

			_, err := wrapper.CreateObject(objNamespace, drainerConfig)
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

	// Verify the test controller reconciles the creation of the runtime object.
	// There should be 2 ResyncPeriods in 30 seconds so verify there were at least
	// 2 create events.
	//
	// 		EnsureCreated: >=2, EnsureDeleted: =0
	//
	{
		o := func() error {
			if r.CreateCount() < 2 {
				return microerror.Maskf(waitError, "r.CreateCount() == %v, want more than %v", r.CreateCount(), 2)
			}
			if r.DeleteCount() != 0 {
				return microerror.Maskf(waitError, "r.DeleteCount() == %v, want %v", r.DeleteCount(), 0)
			}

			return nil
		}
		b := backoff.NewMaxRetries(30, 1*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Once we ensured the reconciled runtime object got processed by the
	// controller, verify the deletion timestamp and finalizers.
	{
		obj, err := wrapper.GetObject(objName, objNamespace)
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

		expectedFinalizers := []string{
			finalizer,
		}
		sort.Strings(expectedFinalizers)

		finalizers := accessor.GetFinalizers()
		sort.Strings(finalizers)

		if !reflect.DeepEqual(finalizers, expectedFinalizers) {
			t.Fatalf("finalizers == %v, want %v", finalizers, expectedFinalizers)
		}
	}

	// Delete the object.
	{
		err := wrapper.DeleteObject(objName, objNamespace)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Verify that the test resource received at least 1 deletion event.
	{
		o := func() error {
			if r.DeleteCount() < 3 {
				return microerror.Maskf(waitError, "r.DeleteCount() == %v, want at least %v", r.DeleteCount(), 3)
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
		obj, err := wrapper.GetObject(objName, objNamespace)
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

		expectedFinalizers := []string{
			finalizer,
		}
		sort.Strings(expectedFinalizers)

		finalizers := accessor.GetFinalizers()
		sort.Strings(finalizers)

		if !reflect.DeepEqual(finalizers, expectedFinalizers) {
			t.Fatalf("finalizers == %v, want %v", finalizers, expectedFinalizers)
		}
	}

	// Set the resource error function to nil to not return any error anymore. The
	// finalizer should be removed with the next reconciliation now.
	{
		r.SetReturnErrorFunc(nil)
	}

	// Verify that the object is completely gone now.
	{
		o := func() error {
			_, err := wrapper.GetObject(objName, objNamespace)
			if drainerconfig.IsNotFound(err) {
				return nil
			} else if err != nil {
				return microerror.Mask(err)
			}

			return microerror.Maskf(waitError, "object %#q still exists", fmt.Sprintf("%s/%s", objName, objNamespace))
		}
		b := backoff.NewMaxRetries(30, 1*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}
}

func newWrapperLogger(w string) micrologger.Logger {
	var err error

	var l micrologger.Logger
	{
		c := micrologger.Config{}

		l, err = micrologger.New(c)
		if err != nil {
			panic(err)
		}
	}

	return l.With("wrapper", w)
}

func newWrapper(r *testresource.Resource, l micrologger.Logger) (*drainerconfig.Wrapper, error) {
	c := drainerconfig.Config{
		Logger: l,
		Resources: []resource.Interface{
			r,
		},

		Name:      controllerName,
		Namespace: objNamespace,
	}

	w, err := drainerconfig.New(c)
	if err != nil {
		return nil, err
	}

	return w, nil
}
