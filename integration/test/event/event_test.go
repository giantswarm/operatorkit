// +build k8srequired

package event

import (
	"context"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/operatorkit/integration/testresource"
	"github.com/giantswarm/operatorkit/integration/wrapper/configmap"
	"github.com/giantswarm/operatorkit/resource"
)

const (
	configMapName = "test-cm"
	operatorName  = "test-operator"
	testNamespace = "event-test"
)

// Test_Kubernetes_Event is a integration test for Kubernetes events.
// The test verifies that error events are created when reconciliation returns error.
//
// !!! This test does not work with CRs, the controller is not booted !!!
//
func Test_Kubernetes_Event(t *testing.T) {
	ctx := context.Background()

	var err error
	var r *testresource.Resource
	{
		c := testresource.Config{
			Name: "test-resource",
			ReturnErrorFunc: func(obj interface{}) error {
				return microerror.Mask(eventError)
			},
		}

		r, err = testresource.New(c)
		if err != nil {
			t.Fatalf("err == %v, want %v", err, nil)
		}
	}

	var wrapper *configmap.Wrapper
	{
		c := configmap.Config{
			Resources: []resource.Interface{r},
			Name:      operatorName,
			Namespace: testNamespace,
		}

		wrapper, err = configmap.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	wrapper.MustSetup(ctx, testNamespace)
	defer wrapper.MustTeardown(ctx, testNamespace)

	controller := wrapper.Controller()

	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: testNamespace,
		},
		Data: map[string]string{},
	}

	_, err = wrapper.CreateObject(ctx, testNamespace, cm)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// We update the object with a meaningless label to ensure a change in the
	// ResourceVersion of the ConfigMap.
	cm.SetLabels(
		map[string]string{
			"testlabel": "testlabel",
		},
	)
	_, err = wrapper.UpdateObject(testNamespace, cm)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	_, err = controller.Reconcile(reconcile.Request{NamespacedName: types.NamespacedName{
		Name:      cm.GetName(),
		Namespace: cm.GetNamespace(),
	}})
	if err != nil {
		t.Fatal("failed reconciliation", nil, "got", err)
	}

	// run Reconcile multiple times to trigger error events.
	_, err = controller.Reconcile(reconcile.Request{NamespacedName: types.NamespacedName{
		Name:      cm.GetName(),
		Namespace: cm.GetNamespace(),
	}})
	if err != nil {
		t.Fatal("failed reconciliation", nil, "got", err)
	}

	// wait a bit to let events appear on the cm
	time.Sleep(5 * time.Second)

	events, err := wrapper.Events(ctx, cm.GetNamespace())
	if err != nil {
		t.Fatal("failed to get events from configmap", nil, "got", err)
	}
	if len(events) != 1 {
		t.Fatal("failed to create event for configmap")
	}
}
