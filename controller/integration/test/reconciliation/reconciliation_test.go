// +build k8srequired

package reconciliation

import (
	"reflect"
	"testing"

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

// Test_Finalizer_Integration_Reconciliation is a integration test for
// the proper replay and reconciliation of delete events with finalizers.
func Test_Finalizer_Integration_Reconciliation(t *testing.T) {
	objName := "test-obj"
	testFinalizer := "operatorkit.giantswarm.io/test-operator"
	testNamespace := "finalizer-integration-reconciliation-test"
	testOtherFinalizer := "operatorkit.giantswarm.io/other-operator"
	operatorName := "test-operator"

	var err error
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

	nodeconfigWrapper, err := nodeconfig.New(c)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	testWrapper := wrapper.Interface(nodeconfigWrapper)

	testWrapper.MustSetup(testNamespace)
	defer testWrapper.MustTeardown(testNamespace)

	controller := testWrapper.Controller()

	// We start the controller.
	go controller.Boot()

	// We create an object, but add a finalizer of another operator. This will
	// cause the object to continue existing after the controller removes its own
	// finalizer.
	//
	//Creation is retried because the existance of a CRD might have to be ensured.
	obj := &v1alpha1.NodeConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      objName,
			Namespace: testNamespace,
			Finalizers: []string{
				testOtherFinalizer,
			},
		},
	}
	var createdObj interface{}
	operation := func() error {
		createdObj, err = testWrapper.CreateObject(testNamespace, obj)
		if err != nil {
			return microerror.Mask(err)
		}
		return nil
	}
	err = backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We update the object with a meaningless label to ensure a change in the
	// ResourceVersion of the object.
	createdObjAccessor, err := meta.Accessor(createdObj)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	createdObjAccessor.SetLabels(
		map[string]string{
			"testlabel": "testlabel",
		},
	)
	// Setting the labels on createdObj works through the magic or accessors and
	// pointers here.
	_, err = testWrapper.UpdateObject(testNamespace, createdObj)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
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
	operation = func() error {
		if tr.CreateCount() != 4 {
			return microerror.Maskf(countMismatchError, "EnsureCreated was hit %v times, want %v", tr.CreateCount(), 4)
		}
		if tr.DeleteCount() != 0 {
			return microerror.Maskf(countMismatchError, "EnsureDeleted was hit %v times, want %v", tr.DeleteCount(), 0)
		}
		return nil
	}
	err = backoff.Retry(operation, newConstantBackoff(uint64(30)))
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We get the object after the controller has been started.
	resultObj, err := testWrapper.GetObject(objName, testNamespace)
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
	err = testWrapper.DeleteObject(objName, testNamespace)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We use backoff with the absolute maximum amount:
	// 20 second ResyncPeriod + 2 second RateWait + 8 second for safety.
	// The controller should now remove the finalizer and EnsureDeleted should be
	// hit twice immediatly. See https://github.com/giantswarm/giantswarm/issues/2897
	//
	// 		EnsureCreated: 4, EnsureDeleted: 2
	//
	// The controller should also reconcile twice in this period. (The other
	// finalizer is still set, so we reconcile.)
	//
	// 		EnsureCreated: 4, EnsureDeleted: 4
	//
	operation = func() error {
		if tr.CreateCount() != 4 {
			return microerror.Maskf(countMismatchError, "EnsureCreated was hit %v times, want %v", tr.CreateCount(), 4)
		}
		if tr.DeleteCount() != 4 {
			return microerror.Maskf(countMismatchError, "EnsureDeleted was hit %v times, want %v", tr.DeleteCount(), 4)
		}
		return nil
	}
	err = backoff.Retry(operation, newConstantBackoff(uint64(30)))
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We get the object after the controller has handled the deletion event.
	resultObj, err = testWrapper.GetObject(objName, testNamespace)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	resultObjAccessor, err = meta.Accessor(resultObj)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We verify, that our object still exists, but has a DeletionTimestamp set.
	if resultObjAccessor.GetDeletionTimestamp() == nil {
		t.Fatalf("DeletionTimestamp == nil, want non-nil")
	}

	// We define which finalizers we currently expect.
	expectedFinalizers = []string{
		testOtherFinalizer,
	}

	// We verify, that our finalizer is still set.
	if !reflect.DeepEqual(resultObjAccessor.GetFinalizers(), expectedFinalizers) {
		t.Fatalf("finalizers == %v, want %v", resultObjAccessor.GetFinalizers(), expectedFinalizers)
	}

}
