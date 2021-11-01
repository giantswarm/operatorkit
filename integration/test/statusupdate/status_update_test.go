//go:build k8srequired
// +build k8srequired

package statusupdate

import (
	"context"
	"testing"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/giantswarm/operatorkit/v5/api/v1"
	"github.com/giantswarm/operatorkit/v5/integration/wrapper/example"
	"github.com/giantswarm/operatorkit/v5/pkg/resource"
)

const (
	objName       = "test"
	operatorName  = "test-operator"
	testNamespace = "integration-status-update-test"
)

func Test_Finalizer_Integration_StatusUpdate(t *testing.T) {
	var err error

	ctx := context.Background()

	var r resource.Interface
	{
		c := ResourceConfig{
			T: t,
		}

		r, err = NewResource(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	var w *example.Wrapper
	{
		c := example.Config{
			Resources: []resource.Interface{
				r,
			},

			Name:      operatorName,
			Namespace: testNamespace,
		}

		w, err = example.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		w.MustSetup(ctx, testNamespace)
		defer func(w *example.Wrapper) {
			w.Controller().Stop(ctx)
			w.MustTeardown(ctx, testNamespace)
		}(w)
	}

	{
		c := w.Controller()

		go c.Boot(ctx)
		<-c.Booted()
	}

	{
		o := func() error {
			obj := &v1.Example{
				ObjectMeta: metav1.ObjectMeta{
					Name:      objName,
					Namespace: testNamespace,
				},
				Spec: v1.ExampleSpec{
					Field1: "a",
				},
			}
			_, err := w.CreateObject(ctx, testNamespace, obj)
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

	time.Sleep(5 * time.Second)

	{
		newObj, err := w.GetObject(ctx, objName, testNamespace)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		newObjTyped := newObj.(*v1.Example)

		if len(newObjTyped.Status.Conditions) != 1 {
			t.Error("expected one status condition")
		}
		if newObjTyped.Status.Conditions[0].Status != conditionStatus {
			t.Errorf("expected status condition status %#q", conditionStatus)
		}
		if newObjTyped.Status.Conditions[0].Type != conditionType {
			t.Errorf("expected status condition type %#q", conditionType)
		}
	}
}
