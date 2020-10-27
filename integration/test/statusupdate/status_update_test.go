// +build k8srequired

package statusupdate

import (
	"context"
	"testing"
	"time"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/operatorkit/v4/integration/wrapper/drainerconfig"
	"github.com/giantswarm/operatorkit/v4/pkg/resource"
)

const (
	conditionStatus = "testStatus"
	conditionType   = "testType"
)

const (
	objName      = "test-obj"
	operatorName = "test-operator"
)

const (
	testNamespace = "finalizer-integration-statusupdate-test"
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

	var w *drainerconfig.Wrapper
	{
		c := drainerconfig.Config{
			Resources: []resource.Interface{
				r,
			},

			Name:      operatorName,
			Namespace: testNamespace,
		}

		w, err = drainerconfig.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		w.MustSetup(testNamespace)
		defer w.MustTeardown(testNamespace)
	}

	{
		c := w.Controller()

		go c.Boot(ctx)
		<-c.Booted()
	}

	{
		o := func() error {
			drainerConfig := &v1alpha1.DrainerConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      objName,
					Namespace: testNamespace,
				},
				TypeMeta: v1alpha1.NewDrainerTypeMeta(),
			}
			_, err := w.CreateObject(ctx, testNamespace, drainerConfig)
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

		customResource := newObj.(*v1alpha1.DrainerConfig)

		if len(customResource.Status.Conditions) != 1 {
			t.Fatal("expected one status condition")
		}
		if customResource.Status.Conditions[0].Status != conditionStatus {
			t.Fatalf("expected status condition status %#q", conditionStatus)
		}
		if customResource.Status.Conditions[0].Type != conditionType {
			t.Fatalf("expected status condition type %#q", conditionType)
		}
	}
}
