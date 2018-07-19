// +build k8srequired

package testresource

import (
	"context"
	"fmt"
	"time"

	"github.com/giantswarm/microerror"
)

type Config struct {
}

type Resource struct {
	createCount int
	deleteCount int
	returnError bool
}

func New(config Config) (*Resource, error) {
	r := &Resource{
		createCount: 0,
		deleteCount: 0,
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
	r.createCount++
}

func (r *Resource) incrementDeleteCount() {
	r.deleteCount++
}
