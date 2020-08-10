// +build k8srequired

package pause

import (
	"context"
	"testing"
	"time"

	"github.com/giantswarm/apiextensions/v2/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/operatorkit/v2/integration/testresource"
	"github.com/giantswarm/operatorkit/v2/integration/wrapper"
	"github.com/giantswarm/operatorkit/v2/integration/wrapper/drainerconfig"
	"github.com/giantswarm/operatorkit/v2/pkg/resource"
)

const (
	objName      = "test-obj"
	objNamespace = "integration-controlflow-test"

	controllerName = "test-controller"
	resourceName   = "test-resource"
)

// Test_Integration_Pause is an integration test to check that the pausing
// annotation feature works as expected.
func Test_Integration_Pause(t *testing.T) {
	var err error

	ctx := context.Background()

	var r *testresource.Resource
	{
		c := testresource.Config{
			Name: resourceName,
		}

		r, err = testresource.New(c)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	var w wrapper.Interface
	{
		c := drainerconfig.Config{
			Resources: []resource.Interface{
				r,
			},

			Name:      controllerName,
			Namespace: objNamespace,
		}

		w, err = drainerconfig.New(c)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Start controller.
	{
		controller := w.Controller()

		go controller.Boot(ctx)
		select {
		case <-controller.Booted():
		case <-time.After(30 * time.Second):
			t.Fatalf("failed to wait for controller to boot")
		}
	}

	// Setup the test namespace.
	{
		w.MustSetup(objNamespace)
		defer w.MustTeardown(objNamespace)
	}

	// Create a runtime object we can reconcile for our integration test.
	{
		o := func() error {
			drainerConfig := &v1alpha1.DrainerConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      objName,
					Namespace: objNamespace,
				},
			}

			_, err := w.CreateObject(objNamespace, drainerConfig)
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
			if r.CreateCount() <= 2 {
				return microerror.Maskf(waitError, "resource.CreateCount() == %v, want more than %v", r.CreateCount(), 2)
			}
			if r.DeleteCount() != 0 {
				return microerror.Maskf(waitError, "resource.DeleteCount() == %v, want %v", r.DeleteCount(), 0)
			}

			return nil
		}
		b := backoff.NewMaxRetries(30, 1*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// Add the pausing annotation to the reconciled runtime object. This should
	// cause the reconciliation to be paused.
	{
		obj, err := w.GetObject(objName, objNamespace)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}

		accessor, err := meta.Accessor(obj)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}

		accessor.SetAnnotations(
			map[string]string{
				"operatorkit.giantswarm.io/paused": "true",
			},
		)
		_, err = w.UpdateObject(objNamespace, accessor)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// We wait a bit to make sure reconciliation does really not happen at this
	// point anymore meanwhile.
	{
		time.Sleep(30 * time.Second)
	}

	// Verify the controller did not reconcile the object anymore due to the
	// pausing annotation. There should now still be 2 the already registered
	// create counts and 0 delete counts.
	//
	// 		EnsureCreated: >2, EnsureDeleted: =0
	//
	{
		o := func() error {
			if r.CreateCount() <= 2 {
				return microerror.Maskf(waitError, "resource.CreateCount() == %v, want more than %v", r.CreateCount(), 2)
			}
			if r.DeleteCount() != 0 {
				return microerror.Maskf(waitError, "resource.DeleteCount() == %v, want %v", r.DeleteCount(), 0)
			}

			return nil
		}
		b := backoff.NewMaxRetries(30, 1*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	// We remove the pausing annotation in order to see the runtime object to be
	// reconciled again.
	{
		obj, err := w.GetObject(objName, objNamespace)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}

		accessor, err := meta.Accessor(obj)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}

		accessor.SetAnnotations(
			map[string]string{},
		)
		_, err = w.UpdateObject(objNamespace, accessor)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// We wait a bit to make sure reconciliation does now happen again
	{
		time.Sleep(30 * time.Second)
	}

	// Verify the controller did reconcile the object again due to the removed
	// pausing annotation. There should now be 4 registered create counts and 0
	// delete counts.
	//
	//      EnsureCreated: >4, EnsureDeleted: =0
	//
	{
		o := func() error {
			if r.CreateCount() <= 4 {
				return microerror.Maskf(waitError, "resource.CreateCount() == %v, want more than %v", r.CreateCount(), 2)
			}
			if r.DeleteCount() != 0 {
				return microerror.Maskf(waitError, "resource.DeleteCount() == %v, want %v", r.DeleteCount(), 0)
			}

			return nil
		}
		b := backoff.NewMaxRetries(30, 1*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}
}
