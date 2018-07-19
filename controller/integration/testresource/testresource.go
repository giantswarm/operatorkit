// +build k8srequired

package testresource

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/giantswarm/microerror"
)

type Config struct {
}

type Resource struct {
	createCount int
	deleteCount int
	mutex       sync.Mutex
	returnError bool
}

func New(config Config) (*Resource, error) {
	r := &Resource{
		createCount: 0,
		deleteCount: 0,
		mutex:       sync.Mutex{},
		returnError: false,
	}

	return r, nil
}

func (r *Resource) CreateCount() int {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.createCount
}

func (r *Resource) DeleteCount() int {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.deleteCount
}

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	fmt.Printf("\n")
	fmt.Printf("%s: EnsureCreated of test resource executed\n", time.Now())
	fmt.Printf("\n")
	r.incrementCreateCount()
	if r.returnError {
		return microerror.Mask(testError)
	}
	return nil
}

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	fmt.Printf("\n")
	fmt.Printf("%s: EnsureDeleted of test resource executed\n", time.Now())
	fmt.Printf("\n")
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
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.createCount++
}

func (r *Resource) incrementDeleteCount() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.deleteCount++
}
