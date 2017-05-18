package tpr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

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

func TestCreateTPR(t *testing.T) {
	clientset := newClientset(3)

	config := Config{
		Clientset: clientset,

		Name:        "testname",
		Domain:      "example.com",
		Version:     "v1test1",
		Description: "Test Desc",
	}

	tpr, err := New(config)
	assert.NoError(t, err, "New")

	resp, err := clientset.ExtensionsV1beta1().ThirdPartyResources().List(v1.ListOptions{})
	assert.Equal(t, 0, len(resp.Items))

	err = tpr.create()
	assert.Nil(t, err)

	resp, err = clientset.ExtensionsV1beta1().ThirdPartyResources().List(v1.ListOptions{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp.Items))
	assert.Equal(t, config.Name+"."+config.Domain, resp.Items[0].Name)
	assert.Equal(t, 1, len(resp.Items[0].Versions))
	assert.Equal(t, "v1test1", resp.Items[0].Versions[0].Name)
	assert.Equal(t, "Test Desc", resp.Items[0].Description)
}
