// +build k8srequired

package testresource

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/kubernetes"
)

type Config struct {
	K8sClient kubernetes.Interface
	Logger    micrologger.Logger
	Name      string
}

type Resource struct {
	k8sClient kubernetes.Interface
	logger    micrologger.Logger

	createCount int
	deleteCount int
	name        string
	returnError bool
}

func New(config Config) (*Resource, error) {
	r := &Resource{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		createCount: 0,
		deleteCount: 0,
		name:        config.Name,
		returnError: false,
	}

	return r, nil
}

func (r *Resource) CreateCount() int {
	return r.createCount
}

func (r *Resource) DeleteCount() int {
	return r.deleteCount
}

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.incrementCreateCount()
	if r.returnError {
		return microerror.Mask(testError)
	}
	return nil
}

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.incrementDeleteCount()
	if r.returnError {
		return microerror.Mask(testError)
	}
	return nil
}

func (r *Resource) Name() string {
	return "testresource"
}

func (r *Resource) ReturnError(returnError bool) {
	r.returnError = returnError
}

func (r *Resource) incrementCreateCount() {
	r.createCount++
}

func (r *Resource) incrementDeleteCount() {
	r.deleteCount++
}
