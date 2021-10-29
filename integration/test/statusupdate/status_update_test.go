//go:build k8srequired
// +build k8srequired

package statusupdate

import (
	"context"
	"testing"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/operatorkit/v5/integration/wrapper/customresourcedefinition"
	"github.com/giantswarm/operatorkit/v5/pkg/resource"
)

const (
	objName      = "tests.example.com"
	operatorName = "test-operator"
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

	var w *customresourcedefinition.Wrapper
	{
		c := customresourcedefinition.Config{
			Resources: []resource.Interface{
				r,
			},

			Name: operatorName,
		}

		w, err = customresourcedefinition.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		w.MustSetup(ctx, "")
		defer w.MustTeardown(ctx, "")
	}

	{
		c := w.Controller()

		go c.Boot(ctx)
		<-c.Booted()
	}

	{
		o := func() error {
			drainerConfig := &apiextensionsv1.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name: objName,
				},
				Spec: apiextensionsv1.CustomResourceDefinitionSpec{
					Group: "example.com",
					Names: apiextensionsv1.CustomResourceDefinitionNames{
						Plural:   "tests",
						Singular: "test",
						Kind:     "Test",
						ListKind: "TestList",
					},
					Scope: "Cluster",
					Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
						{
							Name:    "v1alpha1",
							Served:  true,
							Storage: true,
							Schema: &apiextensionsv1.CustomResourceValidation{
								OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
									Type: "object",
								},
							},
						},
					},
				},
			}
			_, err := w.CreateObject(ctx, "", drainerConfig)
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
		newObj, err := w.GetObject(ctx, objName, "")
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		customResource := newObj.(*apiextensionsv1.CustomResourceDefinition)

		if len(customResource.Status.Conditions) != 3 {
			t.Fatal("expected three status conditions")
		}
		if customResource.Status.Conditions[2].Status != conditionStatus {
			t.Fatalf("expected status condition status %#q", conditionStatus)
		}
		if customResource.Status.Conditions[2].Type != conditionType {
			t.Fatalf("expected status condition type %#q", conditionType)
		}
	}
}
