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

	"github.com/giantswarm/operatorkit/controller/integration/testresource"
	"github.com/giantswarm/operatorkit/controller/integration/wrapper/drainerconfig"
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
// multiple controllers do not interfere with finalizer management where they
// should not act. This is to prevent situations in which multiple operators
// reconcile the same runtime object in different versions. Then finalizer
// management should only be done by controllers that express to handle the
// reconciliation for an object.
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

	// We create three wrappers for controllers and resource sets respectively.
	// They are run to reconcile the same object in order to verify the finalizer
	// management. Only wrapperA expresses to handle the object via its handle
	// func. So neither wrapperB nor wrapperC are allowed to remove the finalizer
	// of the reconciled runtime object.
	var wrapperA *drainerconfig.Wrapper
	var wrapperB *drainerconfig.Wrapper
	var wrapperC *drainerconfig.Wrapper
	{
		wrapperA, err = newWrapper(r, newWrapperLogger("a"), newTrueHandlesFunc())
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
		wrapperB, err = newWrapper(r, newWrapperLogger("b"), newFalseHandlesFunc())
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
		wrapperC, err = newWrapper(r, newWrapperLogger("c"), newFalseHandlesFunc())
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Start controllers.
	{
		controllerA := wrapperA.Controller()
		controllerB := wrapperB.Controller()
		controllerC := wrapperC.Controller()

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

	// Setup the test namespace. We use the wrapperA. It makes no difference which
	// wrapper we use to do this. It is only important to do the setup once for
	// all.
	{
		wrapperA.MustSetup(objNamespace)
		defer wrapperA.MustTeardown(objNamespace)
	}

	// Create an object. Creation is retried because the CRD might still
	// not be ensured. Again, we just use wrapperA to do this.
	{
		o := func() error {
			drainerConfig := &v1alpha1.DrainerConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      objName,
					Namespace: objNamespace,
				},
			}

			_, err := wrapperA.CreateObject(objNamespace, drainerConfig)
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

	// Verify controllers reconcile creation of that object. There should be 2
	// ResyncPeriods in 30 seconds so verify there were at least 2 create events.
	// They should be made by the controller of wrapperA.
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
	// controllers, verify the deletion timestamp and finalizers.
	{
		obj, err := wrapperA.GetObject(objName, objNamespace)
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
		err := wrapperA.DeleteObject(objName, objNamespace)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Verify that the test resource received at least 1 deletion event. When the
	// fix for the test is not applied, the test fails with DeleteCount() == 2.
	// This is because controllers of wrapperB and wrapperC remove the finalizer
	// early, which they should not.
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
		obj, err := wrapperA.GetObject(objName, objNamespace)
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

	// Set the resourceC error function to nil to not return any error
	// anymore. The finalizer should be removed with the next
	// reconciliation now.
	{
		r.SetReturnErrorFunc(nil)
	}

	// Verify that the object is completely gone now.
	{
		o := func() error {
			_, err := wrapperA.GetObject(objName, objNamespace)
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

func newWrapper(r *testresource.Resource, l micrologger.Logger, h func(obj interface{}) bool) (*drainerconfig.Wrapper, error) {
	c := drainerconfig.Config{
		HandlesFunc: h,
		Logger:      l,
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

func newTrueHandlesFunc() func(obj interface{}) bool {
	return func(obj interface{}) bool {
		return true
	}
}

func newFalseHandlesFunc() func(obj interface{}) bool {
	return func(obj interface{}) bool {
		return false
	}
}
