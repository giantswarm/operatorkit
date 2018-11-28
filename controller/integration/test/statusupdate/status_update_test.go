// +build k8srequired

package statusupdate

import (
	"testing"
	"time"

	"github.com/giantswarm/apiextensions/pkg/apis/core/v1alpha1"
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/controller/integration/wrapper/nodeconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	var r controller.Resource
	{
		c := ResourceConfig{
			T: t,
		}

		r, err = NewResource(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	var nodeConfigWrapper *nodeconfig.Wrapper
	{
		c := nodeconfig.Config{
			Resources: []controller.Resource{
				r,
			},

			Name:      operatorName,
			Namespace: testNamespace,
		}

		nodeConfigWrapper, err = nodeconfig.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		nodeConfigWrapper.MustSetup(testNamespace)
		defer nodeConfigWrapper.MustTeardown(testNamespace)
	}

	{
		c := nodeConfigWrapper.Controller()

		go c.Boot()
		<-c.Booted()
	}

	{
		o := func() error {
			nodeConfig := &v1alpha1.NodeConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      objName,
					Namespace: testNamespace,
				},
				TypeMeta: v1alpha1.NewNodeTypeMeta(),
			}
			_, err := nodeConfigWrapper.CreateObject(testNamespace, nodeConfig)
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
		newObj, err := nodeConfigWrapper.GetObject(objName, testNamespace)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		customResource := newObj.(*v1alpha1.NodeConfig)

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
