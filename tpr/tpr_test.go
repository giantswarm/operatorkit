package tpr

import (
	"fmt"
	"testing"

	"github.com/cenkalti/backoff"
	"github.com/stretchr/testify/assert"

	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
)

func newClientset(nodes int) *fake.Clientset {
	clientset := fake.NewSimpleClientset()
	for i := 0; i < nodes; i++ {
		n := &v1.Node{}
		n.Name = fmt.Sprintf("node%d", i)
		clientset.CoreV1().Nodes().Create(n)
	}
	return clientset
}

func TestKindGroupAndAPIVersion(t *testing.T) {
	clientset := newClientset(3)

	config := DefaultConfig()
	config.K8sClient = clientset

	config.Name = "test-name.example.com"
	config.Version = "v1test1"
	config.Description = "Test Desc"

	tpr, err := New(config)
	assert.NoError(t, err, "New")

	assert.Equal(t, "TestName", tpr.Kind())
	assert.Equal(t, "example.com/v1test1", tpr.APIVersion())
	assert.Equal(t, config.Name, tpr.Name())
	assert.Equal(t, "example.com", tpr.Group())

	// Rest of tests should be covered in extractKindAndGroup tests.
}

func TestEndpoint(t *testing.T) {
	clientset := newClientset(3)

	config := DefaultConfig()
	config.K8sClient = clientset

	config.Name = "test-name.example.com"
	config.Version = "v1test1"
	config.Description = "Test Desc"

	tpr, err := New(config)
	assert.NoError(t, err, "New")

	tests := []struct {
		namespace        string
		expectedEndpoint string
	}{
		{
			namespace:        "default",
			expectedEndpoint: "/apis/example.com/v1test1/namespaces/default/testnames",
		},
		{
			namespace:        "",
			expectedEndpoint: "/apis/example.com/v1test1/testnames",
		},
	}
	for i, tc := range tests {
		endpoint := tpr.Endpoint(tc.namespace)
		assert.Equal(t, tc.expectedEndpoint, endpoint, "#%d", i)
	}
}

func TestWatchEndpoint(t *testing.T) {
	clientset := newClientset(3)

	config := DefaultConfig()
	config.K8sClient = clientset

	config.Name = "test-name.example.com"
	config.Version = "v1test1"
	config.Description = "Test Desc"

	tpr, err := New(config)
	assert.NoError(t, err, "New")

	tests := []struct {
		namespace             string
		expectedWatchEndpoint string
	}{
		{
			namespace:             "default",
			expectedWatchEndpoint: "/apis/example.com/v1test1/namespaces/default/watch/testnames",
		},
		{
			namespace:             "",
			expectedWatchEndpoint: "/apis/example.com/v1test1/watch/testnames",
		},
	}
	for i, tc := range tests {
		watchEndpoint := tpr.WatchEndpoint(tc.namespace)
		assert.Equal(t, tc.expectedWatchEndpoint, watchEndpoint, "#%d", i)
	}
}

func TestCreateTPR(t *testing.T) {
	clientset := newClientset(3)

	config := DefaultConfig()
	config.K8sClient = clientset

	config.Name = "test-name.example.com"
	config.Version = "v1test1"
	config.Description = "Test Desc"

	tpr, err := New(config)
	assert.NoError(t, err, "New")

	resp, err := clientset.ExtensionsV1beta1().ThirdPartyResources().List(apismetav1.ListOptions{})
	assert.Equal(t, 0, len(resp.Items))

	initBackOff := backoff.NewExponentialBackOff()
	initBackOff.MaxElapsedTime = 10

	err = tpr.create(initBackOff)
	assert.Nil(t, err)

	resp, err = clientset.ExtensionsV1beta1().ThirdPartyResources().List(apismetav1.ListOptions{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp.Items))
	assert.Equal(t, config.Name, resp.Items[0].Name)
	assert.Equal(t, 1, len(resp.Items[0].Versions))
	assert.Equal(t, "v1test1", resp.Items[0].Versions[0].Name)
	assert.Equal(t, "Test Desc", resp.Items[0].Description)
}
