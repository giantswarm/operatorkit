package testresource

import (
	"context"
	"sync"

	"github.com/giantswarm/microerror"
)

type Config struct {
	Name            string
	ReturnErrorFunc func(obj interface{}) error
}

type Resource struct {
	createCount     int
	deleteCount     int
	mutex           sync.Mutex
	name            string
	returnErrorFunc func(obj interface{}) error
}

func New(config Config) (*Resource, error) {
	if config.Name == "" {
		config.Name = "test-resource"
	}

	r := &Resource{
		createCount:     0,
		deleteCount:     0,
		mutex:           sync.Mutex{},
		name:            config.Name,
		returnErrorFunc: config.ReturnErrorFunc,
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
	r.incrementCreateCount()

	r.mutex.Lock()
	errFunc := r.returnErrorFunc
	r.mutex.Unlock()

	if errFunc != nil {
		err := errFunc(obj)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.incrementDeleteCount()

	r.mutex.Lock()
	errFunc := r.returnErrorFunc
	r.mutex.Unlock()

	if errFunc != nil {
		err := errFunc(obj)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (r *Resource) Name() string {
	return r.name
}

func (r *Resource) SetReturnErrorFunc(f func(obj interface{}) error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.returnErrorFunc = f
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
