package logresource

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/micrologger/loggercontext"

	"github.com/giantswarm/operatorkit/framework"
)

const (
	// Name is the identifier of the resource.
	Name                 = "log"
	PostOperationMessage = "executed resource operation without errors"
	PreOperationMessage  = "start to execute resource operation"
)

// Config represents the configuration used to create a new log resource.
type Config struct {
	// Dependencies.
	Logger   micrologger.Logger
	Resource framework.Resource
}

// DefaultConfig provides a default configuration to create a new log resource
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger:   nil,
		Resource: nil,
	}
}

type Resource struct {
	// Dependencies.
	logger   micrologger.Logger
	resource framework.Resource
}

// New creates a new configured log resource.
func New(config Config) (*Resource, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Resource must not be empty")
	}

	newResource := &Resource{
		// Dependencies.
		logger: config.Logger.With(
			"underlyingResource", config.Resource.Underlying().Name(),
		),
		resource: config.Resource,
	}

	return newResource, nil
}

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	container, ok := loggercontext.FromContext(ctx)
	if ok {
		container.KeyVals["function"] = "GetCurrentState"
		defer delete(container.KeyVals, "function")
	}

	v, err := r.resource.GetCurrentState(ctx, obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	container, ok := loggercontext.FromContext(ctx)
	if ok {
		container.KeyVals["function"] = "GetDesiredState"
		defer delete(container.KeyVals, "function")
	}

	v, err := r.resource.GetDesiredState(ctx, obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *Resource) NewUpdatePatch(ctx context.Context, obj, cur, des interface{}) (*framework.Patch, error) {
	container, ok := loggercontext.FromContext(ctx)
	if ok {
		container.KeyVals["function"] = "NewUpdatePatch"
		defer delete(container.KeyVals, "function")
	}

	v, err := r.resource.NewUpdatePatch(ctx, obj, cur, des)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *Resource) NewDeletePatch(ctx context.Context, obj, cur, des interface{}) (*framework.Patch, error) {
	container, ok := loggercontext.FromContext(ctx)
	if ok {
		container.KeyVals["function"] = "NewDeletePatch"
		defer delete(container.KeyVals, "function")
	}

	v, err := r.resource.NewDeletePatch(ctx, obj, cur, des)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) ApplyCreateChange(ctx context.Context, obj, cre interface{}) error {
	container, ok := loggercontext.FromContext(ctx)
	if ok {
		container.KeyVals["function"] = "ApplyCreateChange"
		defer delete(container.KeyVals, "function")
	}

	err := r.resource.ApplyCreateChange(ctx, obj, cre)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, del interface{}) error {
	container, ok := loggercontext.FromContext(ctx)
	if ok {
		container.KeyVals["function"] = "ApplyDeleteChange"
		defer delete(container.KeyVals, "function")
	}

	err := r.resource.ApplyDeleteChange(ctx, obj, del)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, upd interface{}) error {
	container, ok := loggercontext.FromContext(ctx)
	if ok {
		container.KeyVals["function"] = "ApplyUpdateChange"
		defer delete(container.KeyVals, "function")
	}

	err := r.resource.ApplyUpdateChange(ctx, obj, upd)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) Underlying() framework.Resource {
	return r.resource.Underlying()
}
