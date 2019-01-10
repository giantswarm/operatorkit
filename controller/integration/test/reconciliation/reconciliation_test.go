// +build k8srequired

package reconciliation

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/controller/integration/testresource"
	"github.com/giantswarm/operatorkit/controller/integration/wrapper/nodeconfig"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	objName            = "test-obj"
	operatorName       = "test-operator"
	testFinalizer      = "operatorkit.giantswarm.io/test-operator"
	testNamespace      = "finalizer-integration-reconciliation-test"
	testOtherFinalizer = "operatorkit.giantswarm.io/other-operator"
)

// Test_Finalizer_Integration_Reconciliation is a integration test for
// the proper replay and reconciliation of delete events with finalizers.
func Test_Finalizer_Integration_Reconciliation(t *testing.T) {
	var err error

	ctx := context.Background()

	var tr *testresource.Resource
	{
		c := testresource.Config{}

		tr, err = testresource.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	resources := []controller.Resource{
		controller.Resource(tr),
	}

	c := nodeconfig.Config{
		Resources: resources,

		Name:      operatorName,
		Namespace: testNamespace,
	}

	nodeConfigWrapper, err := nodeconfig.New(c)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	nodeConfigWrapper.MustSetup(testNamespace)
	defer nodeConfigWrapper.MustTeardown(testNamespace)

	controller := nodeConfigWrapper.Controller()

	// We start the controller.
	go controller.Boot(ctx)
	<-controller.Booted()

	// We create an object, but add a finalizer of another operator. This will
	// cause the object to continue existing after the controller removes its own
	// finalizer.
	//
	//Creation is retried because the existance of a CRD might have to be ensured.
	var createdNodeConfig *v1alpha1.NodeConfig
	{
		o := func() error {
			nodeConfig := &v1alpha1.NodeConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      objName,
					Namespace: testNamespace,
					Finalizers: []string{
						testOtherFinalizer,
					},
				},
			}
			v, err := nodeConfigWrapper.CreateObject(testNamespace, nodeConfig)
			if err != nil {
				return microerror.Mask(err)
			}
			createdNodeConfig = v.(*v1alpha1.NodeConfig)

			return nil
		}
		b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)

		err = backoff.Retry(o, b)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// We update the object with a meaningless label to ensure a change in the
	// ResourceVersion of the object.
	{
		o := func() error {
			createdNodeConfig.SetLabels(map[string]string{"testlabel": "testlabel"})

			_, err = nodeConfigWrapper.UpdateObject(testNamespace, createdNodeConfig)
			if err != nil {
				return microerror.Mask(err)
			}

			return nil
		}
		b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)

		err = backoff.Retry(o, b)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// We use backoff with the absolute maximum amount:
	// 20 second ResyncPeriod + 2 second RateWait + 8 second for safety.
	// The controller should now add the finalizer and EnsureCreated should be hit
	// once immediatly.
	//
	// 		EnsureCreated: 1, EnsureDeleted: 0
	//
	// Then we hit EnsureCreated once more because we updated the object with a new
	// label.
	//
	// 		EnsureCreated: 2, EnsureDeleted: 0
	//
	// The controller should reconcile twice in this period.
	//
	// 		EnsureCreated: 4, EnsureDeleted: 0
	//
	{
		o := func() error {
			if tr.CreateCount() != 4 {
				return microerror.Maskf(countMismatchError, "EnsureCreated was hit %v times, want %v", tr.CreateCount(), 4)
			}
			if tr.DeleteCount() != 0 {
				return microerror.Maskf(countMismatchError, "EnsureDeleted was hit %v times, want %v", tr.DeleteCount(), 0)
			}
			return nil
		}
		b := backoff.NewMaxRetries(30, 1*time.Second)

		err = backoff.Retry(o, b)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// We get the object after the controller has been started.
	resultObj, err := nodeConfigWrapper.GetObject(objName, testNamespace)
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

	// We delete the object now.
	err = nodeConfigWrapper.DeleteObject(objName, testNamespace)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// Verify the deletion event was processed.
	//
	// 		EnsureCreated: 4, EnsureDeleted: 1
	{
		o := func() error {
			if tr.CreateCount() != 4 {
				return microerror.Maskf(countMismatchError, "EnsureCreated was hit %v times, want %v", tr.CreateCount(), 4)
			}
			if tr.DeleteCount() != 1 {
				return microerror.Maskf(countMismatchError, "EnsureDeleted was hit %v times, want %v", tr.DeleteCount(), 1)
			}
			return nil
		}
		b := backoff.NewMaxRetries(30, 1*time.Second)

		err = backoff.Retry(o, b)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// Verify deletion timestamp and finalizer.
	{
		o := func() error {
			obj, err := nodeConfigWrapper.GetObject(objName, testNamespace)
			if err != nil {
				return microerror.Mask(err)
			}

			accessor, err := meta.Accessor(obj)
			if err != nil {
				return microerror.Mask(err)
			}

			if accessor.GetDeletionTimestamp() == nil {
				microerror.Maskf(waitError, "DeletionTimestamp == %v, want non nil", accessor.GetDeletionTimestamp())
			}

			finalizers := accessor.GetFinalizers()
			expectedFinalizers := []string{
				testOtherFinalizer,
			}
			if !reflect.DeepEqual(finalizers, expectedFinalizers) {
				microerror.Maskf(waitError, "finalizers == %v, want %v", finalizers, expectedFinalizers)
			}

			return nil
		}
		b := backoff.NewMaxRetries(10, 1*time.Second)

		err := backoff.Retry(o, b)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}
}
