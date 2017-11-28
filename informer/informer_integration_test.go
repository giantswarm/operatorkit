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
	"encoding/json"
	"flag"
	"fmt"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/operatorkit/client/k8scrdclient"
	"github.com/giantswarm/operatorkit/client/k8sextclient"
)

var (
	address string

	caFile  string
	crtFile string
	keyFile string

	newCRD            *apiextensionsv1beta1.CustomResourceDefinition
	newCRDClient      *k8scrdclient.CRDClient
	newK8sExtClient   apiextensionsclient.Interface
	newInformer       *Informer
	newWatcherFactory WatcherFactory
)

type Test struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              TestSpec `json:"spec"`
}

type TestSpec struct {
	ID string `json:"id" yaml:"id"`
}

type TestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Test `json:"items"`
}

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
		c := k8sextclient.DefaultConfig()

		c.Logger = microloggertest.New()

		c.Address = address
		c.InCluster = false
		c.TLS.CAFile = caFile
		c.TLS.CrtFile = crtFile
		c.TLS.KeyFile = keyFile

		newK8sExtClient, err = k8sextclient.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v", err))
		}
	}

	{
		c := k8scrdclient.DefaultConfig()

		c.Logger = microloggertest.New()
		c.K8sExtClient = newK8sExtClient

		newCRDClient, err = k8scrdclient.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v", err))
		}
	}

	{
		newCRD = &apiextensionsv1beta1.CustomResourceDefinition{
			TypeMeta: metav1.TypeMeta{
				APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
				Kind:       "CustomResourceDefinition",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "tests.foo.giantswarm.io",
			},
			Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
				Group:   "foo.giantswarm.io",
				Scope:   "Cluster",
				Version: "v1alpha1",
				Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
					Kind:     "Test",
					Plural:   "tests",
					Singular: "test",
				},
			},
		}
	}

	{
		zeroObjectFactory := &ZeroObjectFactoryFuncs{
			NewObjectFunc:     func() runtime.Object { return &Test{} },
			NewObjectListFunc: func() runtime.Object { return &TestList{} },
		}
		newWatcherFactory = NewWatcherFactory(newK8sExtClient.Discovery().RESTClient(), filepath.Join("apis", "foo.giantswarm.io", "v1alpha1", "watch", "tests"), zeroObjectFactory)
	}
}

func testNewInformer(t *testing.T, rateWait, resyncPeriod time.Duration) *Informer {
	c := DefaultConfig()

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
	p := filepath.Join("apis", "foo.giantswarm.io", "v1alpha1", "tests")

	o := &Test{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "foo.giantswarm.io/v1alpha1",
			Kind:       "Test",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: ID,
		},
		Spec: TestSpec{
			ID: ID,
		},
	}

	b, err := json.Marshal(o)
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}

	err = newK8sExtClient.Discovery().RESTClient().Post().AbsPath(p).Body(b).Do().Error()
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}
}

func testDeleteCRO(t *testing.T, ID string) {
	p := filepath.Join("apis", "foo.giantswarm.io", "v1alpha1", "tests", ID)

	err := newK8sExtClient.Discovery().RESTClient().Delete().AbsPath(p).Do().Error()
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}
}

func testSetup(t *testing.T) {
	testTeardown(t)

	err := newCRDClient.Ensure(context.TODO(), newCRD, backoff.NewExponentialBackOff())
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}
}

func testTeardown(t *testing.T) {
	err := newK8sExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(newCRD.Name, nil)
	if errors.IsNotFound(err) {
		// fall though
	} else if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}
}
