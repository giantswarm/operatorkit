// +build integration

package informer

/*
	Usage:

		go test -tags=integration ./informer [FLAGS]

	Flags:

		-integration.address string
			Kubernetes API server address (default "https://$(minikube ip):8443")
		-integration.ca string
			CA file path (default "$HOME/.minikube/ca.crt")
		-integration.crt string
			certificate file path (default "$HOME/.minikube/apiserver.crt")
		-integration.key string
			key file path (default "$HOME/.minikube/apiserver.key")
*/

import (
	"context"
	"flag"
	"fmt"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/cenk/backoff"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/operatorkit/client/crdclient"
	"github.com/giantswarm/operatorkit/crd"
	fakecrd "github.com/giantswarm/operatorkit/crd/fake"
)

var (
	address string

	caFile  string
	crtFile string
	keyFile string

	newCRD            *crd.CRD
	newCRDClient      apiextensionsclient.Interface
	newInformer       *Informer
	newWatcherFactory WatcherFactory
)

func init() {
	var err error

	{
		u, err := user.Current()
		homePath := func(relativePath string) string {
			if err != nil {
				return ""
			}
			return path.Join(u.HomeDir, relativePath)
		}

		var serverDefault string
		{
			out, err := exec.Command("minikube", "ip").Output()
			if err == nil {
				minikubeIP := strings.TrimSpace(string(out))
				serverDefault = "https://" + string(minikubeIP) + ":8443"
			}
		}

		flag.StringVar(&address, "integration.address", serverDefault, "Kubernetes API server address.")

		flag.StringVar(&caFile, "integration.ca", homePath(".minikube/ca.crt"), "CA file path.")
		flag.StringVar(&crtFile, "integration.crt", homePath(".minikube/apiserver.crt"), "Certificate file path.")
		flag.StringVar(&keyFile, "integration.key", homePath(".minikube/apiserver.key"), "Key file path.")
	}

	{
		c := crdclient.DefaultConfig()

		c.Logger = microloggertest.New()

		c.Address = address
		c.InCluster = false
		c.TLS.CAFile = caFile
		c.TLS.CrtFile = crtFile
		c.TLS.KeyFile = keyFile

		newCRDClient, err = crdclient.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v", err))
		}
	}

	{
		c := crd.DefaultConfig()

		c.Group = fakecrd.Group
		c.Kind = fakecrd.Kind
		c.Name = fakecrd.Name
		c.Plural = fakecrd.Plural
		c.Singular = fakecrd.Singular
		c.Scope = fakecrd.Scope
		c.Version = fakecrd.VersionV1

		newCRD, err = crd.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v", err))
		}
	}

	{
		zeroObjectFactory := &ZeroObjectFactoryFuncs{
			NewObjectFunc:     func() runtime.Object { return &fakecrd.CustomObject{} },
			NewObjectListFunc: func() runtime.Object { return &fakecrd.List{} },
		}
		newWatcherFactory = NewWatcherFactory(newCRDClient.Discovery().RESTClient(), newCRD.WatchEndpoint(), zeroObjectFactory)
	}
}

func testNewInformer(t *testing.T, rateWait, resyncPeriod time.Duration) *Informer {
	c := DefaultConfig()

	c.BackOff = backoff.NewExponentialBackOff()
	c.WatcherFactory = newWatcherFactory

	c.RateWait = rateWait
	c.ResyncPeriod = resyncPeriod

	newInformer, err := New(c)
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}

	return newInformer
}

func testAssertCROWithID(t *testing.T, e watch.Event, IDs ...string) {
	m, err := meta.Accessor(e.Object)
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}

	name := m.GetName()
	for _, ID := range IDs {
		if name == ID {
			return
		}
	}

	t.Fatalf("expected one of %#v got %#v", IDs, name)
}

func testCreateCRO(t *testing.T, ID string) {
	p := newCRD.CreateEndpoint()
	b := fakecrd.NewCRO(ID)

	err := newCRDClient.Discovery().RESTClient().Post().AbsPath(p).Body(b).Do().Error()
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}
}

func testDeleteCRO(t *testing.T, ID string) {
	p := newCRD.ResourceEndpoint(ID)

	err := newCRDClient.Discovery().RESTClient().Delete().AbsPath(p).Do().Error()
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}
}

func testSetup(t *testing.T) {
	testTeardown(t)

	err := crd.Ensure(context.TODO(), newCRD, newCRDClient, backoff.NewExponentialBackOff())
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}
}

func testTeardown(t *testing.T) {
	err := newCRDClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(newCRD.Name(), nil)
	if errors.IsNotFound(err) {
		// fall though
	} else if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}
}
