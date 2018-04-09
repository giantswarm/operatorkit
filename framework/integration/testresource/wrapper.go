// +build k8srequired

package testresource

import (
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
	"k8s.io/client-go/kubernetes"
)

var (
	createCount = 0
	deleteCount = 0
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
	Name      string
}

type Wrapper struct {
	Resource framework.Resource
}

func New(config Config) (*Wrapper, error) {
	r := &Resource{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		name: config.Name,
	}

	tr := &Wrapper{
		Resource: r,
	}

	return tr, nil
}

func (r *Wrapper) GetCreateCount() int {
	return createCount
}

func (r *Wrapper) GetDeleteCount() int {
	return deleteCount
}
