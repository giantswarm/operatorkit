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
	nameObj      = "test-obj"
	nameOperator = "test-operator"
)

const (
	namespaceTest = "finalizer-integration-statusupdate-test"
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

			Name:      nameOperator,
			Namespace: namespaceTest,
		}

		nodeConfigWrapper, err = nodeconfig.New(c)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		nodeConfigWrapper.MustSetup(namespaceTest)
		defer nodeConfigWrapper.MustTeardown(namespaceTest)
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
					Name:      nameObj,
					Namespace: namespaceTest,
				},
				TypeMeta: v1alpha1.NewNodeTypeMeta(),
			}
			_, err := nodeConfigWrapper.CreateObject(namespaceTest, nodeConfig)
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
		newObj, err := nodeConfigWrapper.GetObject(nameObj, namespaceTest)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		customResource := newObj.(*v1alpha1.NodeConfig)

		if len(customResource.Status.Conditions) != 1 {
			t.Fatal("expected one status condition")
		}
	}
}
