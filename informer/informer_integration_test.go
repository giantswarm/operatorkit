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
	"flag"
	"fmt"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/operatorkit/client/k8srestconfig"
)

const (
	namespace = "informer-integration-test"
)

var (
	address string

	caFile  string
	crtFile string
	keyFile string

	k8sClient kubernetes.Interface
)

type Test struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              TestSpec `json:"spec"`
}

func (t *Test) DeepCopyObject() runtime.Object {
	return &Test{
		TypeMeta:   t.TypeMeta,
		ObjectMeta: *t.ObjectMeta.DeepCopy(),
		Spec:       t.Spec,
	}
}

type TestSpec struct {
	ID string `json:"id" yaml:"id"`
}

type TestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Test `json:"items"`
}

func (t *TestList) DeepCopyObject() runtime.Object {
	itemsCopy := make([]Test, len(t.Items))
	for i, item := range t.Items {
		itemsCopy[i] = item
	}

	return &TestList{
		TypeMeta: t.TypeMeta,
		ListMeta: *t.ListMeta.DeepCopy(),
		Items:    itemsCopy,
	}
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

	var restConfig *rest.Config
	{
		c := k8srestconfig.DefaultConfig()

		c.Logger = microloggertest.New()

		c.Address = address
		c.InCluster = false
		c.TLS.CAFile = caFile
		c.TLS.CrtFile = crtFile
		c.TLS.KeyFile = keyFile

		restConfig, err = k8srestconfig.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v", err))
		}
	}

	k8sClient, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		panic(fmt.Sprintf("%#v", err))
	}
}

func testNewInformer(t *testing.T, rateWait, resyncPeriod time.Duration) *Informer {
	c := DefaultConfig()

	c.Watcher = k8sClient.CoreV1().ConfigMaps(namespace)

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

func testCreateObj(t *testing.T, ID string) {
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ID,
			Namespace: namespace,
		},
		Data: map[string]string{},
	}

	_, err := k8sClient.CoreV1().ConfigMaps(namespace).Create(cm)
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}
}

func testDeleteObj(t *testing.T, ID string) {
	err := k8sClient.CoreV1().ConfigMaps(namespace).Delete(ID, nil)
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}
}

func testSetup(t *testing.T) {
	testTeardown(t)

	ns := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
		Spec: corev1.NamespaceSpec{},
	}

	_, err := k8sClient.CoreV1().Namespaces().Create(ns)
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}
}

func testTeardown(t *testing.T) {
	err := k8sClient.CoreV1().Namespaces().Delete(namespace, nil)
	if errors.IsNotFound(err) {
		// fall though
	} else if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}
}
