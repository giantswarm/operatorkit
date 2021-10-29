//go:build k8srequired
// +build k8srequired

package error

import (
	"context"
	"testing"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/operatorkit/v5/integration/testresource"
	"github.com/giantswarm/operatorkit/v5/integration/wrapper/configmap"
	"github.com/giantswarm/operatorkit/v5/pkg/resource"
)

const (
	testNamespace = "error-test"
	testObjectA   = "test-object-a"
	testObjectB   = "test-object-b"
)

func Test_Controller_Integration_Error(t *testing.T) {
	var err error

	ctx := context.Background()

	var rA *testresource.Resource
	{
		c := testresource.Config{
			Name: "test-resource-a",
			ReturnErrorFunc: func(obj interface{}) error {
				a, err := meta.Accessor(obj)
				if err != nil {
					return microerror.Mask(err)
				}

				if a.GetName() == testObjectA {
					return microerror.Mask(testError)
				}

				return nil
			},
		}

		rA, err = testresource.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	var rB *testresource.Resource
	{
		c := testresource.Config{
			Name: "test-resource-b",
		}

		rB, err = testresource.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	resources := []resource.Interface{
		rA,
		rB,
	}

	var wrapper *configmap.Wrapper
	{
		c := configmap.Config{
			Resources: resources,

			Name:      "operator-name",
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
	go controller.Boot(ctx)
	<-controller.Booted()

	// We create two test objects. One is used by one resource to error out.
	{
		o := func() error {
			a := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      testObjectA,
					Namespace: testNamespace,
				},
			}

			_, err = wrapper.CreateObject(ctx, testNamespace, a)
			if err != nil {
				return microerror.Mask(err)
			}

			b := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      testObjectB,
					Namespace: testNamespace,
				},
			}

			_, err = wrapper.CreateObject(ctx, testNamespace, b)
			if err != nil {
				return microerror.Mask(err)
			}

			return nil
		}
		b := backoff.NewExponential(2*time.Minute, 10*time.Second)
		err = backoff.Retry(o, b)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// We use a backoff with a max wait of 30 seconds. 20 second ResyncPeriod + 2
	// second RateWait + 8 second for safety. Test resource A should be executed
	// as twice as much as test resource B. This is because test resource A errors
	// when it perceives the test object A.
	{
		o := func() error {
			if rA.CreateCount() != 6 {
				return microerror.Maskf(countMismatchError, "EnsureCreated was hit %v times, want %v", rA.CreateCount(), 6)
			}
			if rA.DeleteCount() != 0 {
				return microerror.Maskf(countMismatchError, "EnsureDeleted was hit %v times, want %v", rA.DeleteCount(), 0)
			}

			if rB.CreateCount() != 4 {
				return microerror.Maskf(countMismatchError, "EnsureCreated was hit %v times, want %v", rB.CreateCount(), 4)
			}
			if rB.DeleteCount() != 0 {
				return microerror.Maskf(countMismatchError, "EnsureDeleted was hit %v times, want %v", rB.DeleteCount(), 0)
			}

			return nil
		}
		b := backoff.NewMaxRetries(30, 1*time.Second)
		err = backoff.Retry(o, b)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}
}
