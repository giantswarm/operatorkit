// +build k8srequired

package testresource

import (
	"context"

	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"
)

type Resource struct {
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	name string
}

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	createCount++
	return nil
}

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	deleteCount++
	return nil
}

func (r *Resource) Name() string {
	return "testresource"
}
