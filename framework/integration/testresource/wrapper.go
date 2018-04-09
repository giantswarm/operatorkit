// +build k8srequired

package testresource

import (
	"context"

	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
	Name      string
}

type Wrapper struct {
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	createCount int
	deleteCount int
	name        string
}

func New(config Config) (*Wrapper, error) {
	w := &Wrapper{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		createCount: 0,
		deleteCount: 0,
		name:        config.Name,
	}

	return w, nil
}

func (w *Wrapper) EnsureCreated(ctx context.Context, obj interface{}) error {
	w.IncrementCreateCount()
	return nil
}

func (w *Wrapper) EnsureDeleted(ctx context.Context, obj interface{}) error {
	w.IncrementDeleteCount()
	return nil
}

func (w *Wrapper) Name() string {
	return "testresource"
}

func (w *Wrapper) GetCreateCount() int {
	return w.createCount
}

func (w *Wrapper) GetDeleteCount() int {
	return w.deleteCount
}

func (w *Wrapper) IncrementCreateCount() {
	w.createCount++
}

func (w *Wrapper) IncrementDeleteCount() {
	w.deleteCount++
}
