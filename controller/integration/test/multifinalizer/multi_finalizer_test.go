// +build k8srequired

package multifinalizer

import (
	"context"
	"testing"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger/microloggertest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/controller/integration/testresource"
	"github.com/giantswarm/operatorkit/controller/integration/wrapper/configmap"
)

const (
	objName       = "test-obj"
	testNamespace = "integration-parallel-test"

	controllerNameA = "test-controller-a"
	resourceNameA   = "test-resource-a"
	testFinalizerA  = "operatorkit.giantswarm.io/test-controller-a"

	controllerNameB = "test-controller-b"
	resourceNameB   = "test-resource-b"
	testFinalizerB  = "operatorkit.giantswarm.io/test-controller-b"
)

func Test_MultiFinalizer(t *testing.T) {
	var err error

	var (
		ctx    = context.Background()
		logger = microloggertest.New()
	)

	var resourceA *testresource.Resource
	{
		c := testresource.Config{
			Name: resourceNameA,
		}

		resourceA, err = testresource.New(c)
		if err != nil {
			t.Fatalf("failed to create resource: %#v", err)
		}
	}

	var resourceB *testresource.Resource
	{
		c := testresource.Config{
			Name: resourceNameB,
			ReturnErrorFunc: func(obj interface{}) error {
				return microerror.Maskf(executionError, "I fail to keep the finalizer forever")
			},
		}

		resourceB, err = testresource.New(c)
		if err != nil {
			t.Fatalf("failed to create resource: %#v", err)
		}
	}

	var harnessA, harnessB *configmap.Wrapper
	{
		harnessA, err = setupHarness(controllerNameA, resourceA)
		if err != nil {
			t.Fatalf("failed to setup controller %#q: %#v", controllerNameA, err)
		}
		harnessB, err = setupHarness(controllerNameB, resourceB)
		if err != nil {
			t.Fatalf("failed to setup controller %#q: %#v", controllerNameB, err)
		}
	}

	// Start controllers.
	{
		controllerA := harnessA.Controller()
		controllerB := harnessB.Controller()

		go controllerA.Boot()
		go controllerB.Boot()
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
	}

	// We setup the namespace in which we test. We use the harness A. It
	// makes no difference if we use the harness A or B.
	{
		harnessA.MustSetup(testNamespace)
		defer harnessA.MustTeardown(testNamespace)
	}

	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm",
			Namespace: testNamespace,
		},
		Data: map[string]string{},
	}

	// Create the object and wait till it has both finalizers.
	{
		_, err := harnessA.CreateObject(cm.Namespace, cm)
		if err != nil {
			t.Fatalf("failed to create ConfigMap %#q in namespace %#q: %#v", cm.Name, cm.Namespace, err)
		}

		o := func() error {
			obj, err := harnessA.GetObject(cm.Name, cm.Namespace)
			if err != nil {
				t.Fatalf("failed to get ConfigMap %#q in namespace %#q: %#v", cm.Name, cm.Namespace, err)
			}

			cm, err := configmap.ToCustomObject(obj)
			if err != nil {
				t.Fatalf("failed to convert the object to ConfigMap: %#v", err)
			}

			var hasFinalizerA, hasFinalizerB bool
			{
				for _, f := range cm.Finalizers {
					switch f {
					case testFinalizerA:
						hasFinalizerA = true
					case testFinalizerB:
						hasFinalizerB = true
					}
				}
			}

			if !hasFinalizerA {
				return microerror.Maskf(waitError, "finalizer %#q is not present in %#v", testFinalizerA, cm.Finalizers)
			}
			if !hasFinalizerB {
				return microerror.Maskf(waitError, "finalizer %#q is not present in %#v", testFinalizerB, cm.Finalizers)
			}

			return nil
		}
		b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)
		n := backoff.NewNotifier(logger, ctx)
		err = backoff.RetryNotify(o, b, n)
		if err != nil {
			t.Fatalf("failed to wait for ConfigMap to have both finalizers: %#v", err)
		}
	}

	// Delete objects and check if it has only finalizer B from
	// a constantly failing resource.
	{
		err := harnessA.DeleteObject(cm.Name, cm.Namespace)
		if err != nil {
			t.Fatalf("failed to delete ConfigMap %#q in namespace %#q: %#v", cm.Name, cm.Namespace, err)
		}

		o := func() error {
			obj, err := harnessA.GetObject(cm.Name, cm.Namespace)
			if err != nil {
				t.Fatalf("failed to get ConfigMap %#q in namespace %#q: %#v", cm.Name, cm.Namespace, err)
			}

			cm, err := configmap.ToCustomObject(obj)
			if err != nil {
				t.Fatalf("failed to convert the object to ConfigMap: %#v", err)
			}

			var hasFinalizerA, hasFinalizerB bool
			{
				for _, f := range cm.Finalizers {
					switch f {
					case testFinalizerA:
						hasFinalizerA = true
					case testFinalizerB:
						hasFinalizerB = true
					}
				}
			}

			if hasFinalizerA {
				return microerror.Maskf(waitError, "finalizer %#q is still present in %#v", testFinalizerA, cm.Finalizers)
			}
			if !hasFinalizerB {
				t.Fatalf("expected finalizer %#q in %#v", testFinalizerB, cm.Finalizers)
			}

			return nil
		}
		b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)
		n := backoff.NewNotifier(logger, ctx)
		err = backoff.RetryNotify(o, b, n)
		if err != nil {
			t.Fatalf("failed to wait for ConfigMap to have %#q finalizer removed: %#v", testFinalizerB, err)
		}
	}

	// Reset the resource and check if deletion counter stays 1 in resource
	// A and increases in resource B. That means only the controller with
	// coresponding finalizer receives deletion events.
	{
		o := func() error {
			if resourceA.DeleteCount() != 1 {
				t.Fatalf("resourceA.DeleteCount == %d, want 1", resourceA.DeleteCount())
			}
			if resourceB.DeleteCount() < 4 {
				return microerror.Maskf(waitError, "resourceB.DeleteCount() is still less than 4")
			}

			return nil
		}
		b := backoff.NewExponential(backoff.ShortMaxWait, backoff.ShortMaxInterval)
		n := backoff.NewNotifier(logger, ctx)
		err := backoff.RetryNotify(o, b, n)
		if err != nil {
			t.Fatalf("failed to wait for resourceB.DeleteCount() being bigger than 3: %#v", err)
		}
	}
}

func setupHarness(controllerName string, resource controller.Resource) (*configmap.Wrapper, error) {
	resources := []controller.Resource{
		resource,
	}

	c := configmap.Config{
		Resources: resources,

		Name:      controllerName,
		Namespace: testNamespace,
	}

	wrapper, err := configmap.New(c)
	if err != nil {
		return nil, err
	}

	return wrapper, nil
}
