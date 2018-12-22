// +build k8srequired

package controlflow

import (
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
	"github.com/giantswarm/operatorkit/controller/integration/wrapper/nodeconfig"
)

const (
	objName       = "test-obj"
	operatorName  = "test-operator"
	testFinalizer = "operatorkit.giantswarm.io/test-operator"
	testNamespace = "finalizer-integration-reconciliation-test"
)

// Test_Finalizer_Integration_Controlflow is an integration test to check that
// errors during deletion prevent the finalizer from removal.
func Test_Finalizer_Integration_Controlflow(t *testing.T) {
	var err error

	var tr *testresource.Resource
	{
		c := testresource.Config{}

		tr, err = testresource.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// We setup and start the controller.
	var testWrapper wrapper.Interface
	{
		c := nodeconfig.Config{
			Resources: []controller.Resource{
				tr,
			},

			Name:      operatorName,
			Namespace: testNamespace,
		}

		testWrapper, err = nodeconfig.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		testWrapper.MustSetup(testNamespace)
		defer testWrapper.MustTeardown(testNamespace)
		newController := testWrapper.Controller()
		go newController.Boot()
		<-newController.Booted()
	}

	// We create an object which is valid and wait for the framework to add a
	// finalizer.
	{
		o := func() error {
			nodeConfig := &v1alpha1.NodeConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      objName,
					Namespace: testNamespace,
				},
			}

			_, err := testWrapper.CreateObject(testNamespace, nodeConfig)
			if err != nil {
				return microerror.Mask(err)
			}

			return nil
		}
		b := backoff.NewExponential(2*time.Second, 10*time.Second)

		err = backoff.Retry(o, b)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
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
	{
		o := func() error {
			if tr.CreateCount() != 2 {
				return microerror.Maskf(waitError, "EnsureCreated was hit %v times, want %v", tr.CreateCount(), 2)
			}
			if tr.DeleteCount() != 0 {
				return microerror.Maskf(waitError, "EnsureDeleted was hit %v times, want %v", tr.DeleteCount(), 0)
			}

			return nil
		}
		b := backoff.NewMaxRetries(20, 1*time.Second)

		err := backoff.Retry(operation, b)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// Verify deletion timestamp and finalizer.
	{
		obj, err := testWrapper.GetObject(objName, testNamespace)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		accessor, err := meta.Accessor(obj)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if accessor.GetDeletionTimestamp() != nil {
			t.Fatalf("DeletionTimestamp != nil, want nil")
		}

		finalizers := accessor.GetFinalizers()
		expectedFinalizers := []string{
			testFinalizer,
		}
		if !reflect.DeepEqual(finalizers, expectedFinalizers) {
			t.Fatalf("finalizers == %v, want %v", finalizers, expectedFinalizers)
		}
	}

	// We set an error function to return an error. This causes the
	// resource to always return an error and should therefore prevent the
	// removal of our finalizer.
	{
		tr.SetReturnErrorFunc(func(obj interface{}) error {
			return microerror.Mask(testError)
		})
	}

	// We delete the object now.
	{
		err := testWrapper.DeleteObject(objName, testNamespace)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
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
	{
		o := func() error {
			if tr.CreateCount() != 2 {
				return microerror.Maskf(waitError, "EnsureCreated was hit %v times, want %v", tr.CreateCount(), 2)
			}
			if tr.DeleteCount() != 2 {
				return microerror.Maskf(waitError, "EnsureDeleted was hit %v times, want %v", tr.DeleteCount(), 2)
			}

			return nil
		}
		b := backoff.NewMaxRetries(20, 1*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// Verify deletion timestamp and finalizer again.
	{
		obj, err := testWrapper.GetObject(objName, testNamespace)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		accessor, err := meta.Accessor(obj)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		if accessor.GetDeletionTimestamp() != nil {
			t.Fatalf("DeletionTimestamp != nil, want nil")
		}

		finalizers := accessor.GetFinalizers()
		expectedFinalizers := []string{
			testFinalizer,
		}
		if !reflect.DeepEqual(finalizers, expectedFinalizers) {
			t.Fatalf("finalizers == %v, want %v", finalizers, expectedFinalizers)
		}
	}

	// We use backoff with the absolute maximum amount:
	// 10 second ResyncPeriod + 2 second RateWait + 8 second for safety.
	//
	// 		EnsureCreated: 2, EnsureDeleted: >3
	//
	{
		o := func() error {
			if tr.CreateCount() != 2 {
				return microerror.Maskf(waitError, "EnsureCreated was hit %v times, want %v", tr.CreateCount(), 2)
			}
			if tr.DeleteCount() > 3 {
				return microerror.Maskf(waitError, "EnsureDeleted was hit %v times, want more than %v", tr.DeleteCount(), 3)
			}

			return nil
		}
		b := backoff.NewExponential(1*time.Second, 20*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// We set the error function to nil to not return any error anymore. Our
	// finalizer should be removed with the next reconciliation now.
	{
		tr.SetReturnErrorFunc(nil)
	}

	// We verify that our object is completely gone now.
	{
		o := func() error {
			_, err = testWrapper.GetObject(objName, testNamespace)
			if nodeconfig.IsNotFound(err) {
				return nil
			} else if err != nil {
				return microerror.Mask(err)
			}

			return microerror.Maskf(waitError, "object %#q in namespace %#q is still not deleted", objName, testNamespace)
		}
		b := backoff.NewExponential(1*time.Second, 30*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}
}
