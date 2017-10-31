package logresource

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

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
	r.logger.Log("action", "start", "component", "operatorkit", "function", "GetCurrentState")

	v, err := r.resource.GetCurrentState(ctx, obj)
	if err != nil {
		r.logger.Log("action", "error", "component", "operatorkit", "function", "GetCurrentState")
		return nil, microerror.Mask(err)
	}

	r.logger.Log("action", "end", "component", "operatorkit", "function", "GetCurrentState")

	return v, nil
}

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	r.logger.Log("action", "start", "component", "operatorkit", "function", "GetDesiredState")

	v, err := r.resource.GetDesiredState(ctx, obj)
	if err != nil {
		r.logger.Log("action", "error", "component", "operatorkit", "function", "GetDesiredState")
		return nil, microerror.Mask(err)
	}

	r.logger.Log("action", "end", "component", "operatorkit", "function", "GetDesiredState")

	return v, nil
}

func (r *Resource) NewUpdatePatch(ctx context.Context, obj, cur, des interface{}) (*framework.Patch, error) {
	r.logger.Log("action", "start", "component", "operatorkit", "function", "NewUpdatePatch")

	v, err := r.resource.NewUpdatePatch(ctx, obj, cur, des)
	if err != nil {
		r.logger.Log("action", "error", "component", "operatorkit", "function", "NewUpdatePatch")
		return nil, microerror.Mask(err)
	}

	r.logger.Log("action", "end", "component", "operatorkit", "function", "NewUpdatePatch")

	return v, nil
}

func (r *Resource) NewDeletePatch(ctx context.Context, obj, cur, des interface{}) (*framework.Patch, error) {
	r.logger.Log("action", "start", "component", "operatorkit", "function", "NewDeletePatch")

	v, err := r.resource.NewDeletePatch(ctx, obj, cur, des)
	if err != nil {
		r.logger.Log("action", "error", "component", "operatorkit", "function", "NewDeletePatch")
		return nil, microerror.Mask(err)
	}

	r.logger.Log("action", "end", "component", "operatorkit", "function", "NewDeletePatch")

	return v, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) ApplyCreateChange(ctx context.Context, obj, cre interface{}) error {
	r.logger.Log("action", "start", "component", "operatorkit", "function", "ApplyCreatePatch")

	err := r.resource.ApplyCreateChange(ctx, obj, cre)
	if err != nil {
		r.logger.Log("action", "error", "component", "operatorkit", "function", "ApplyCreatePatch")
		return microerror.Mask(err)
	}

	r.logger.Log("action", "end", "component", "operatorkit", "function", "ApplyCreatePatch")

	return nil
}

func (r *Resource) ApplyDeleteChange(ctx context.Context, obj, del interface{}) error {
	r.logger.Log("action", "start", "component", "operatorkit", "function", "ApplyDeletePatch")

	err := r.resource.ApplyDeleteChange(ctx, obj, del)
	if err != nil {
		r.logger.Log("action", "error", "component", "operatorkit", "function", "ApplyDeletePatch")
		return microerror.Mask(err)
	}

	r.logger.Log("action", "end", "component", "operatorkit", "function", "ApplyDeletePatch")

	return nil
}

func (r *Resource) ApplyUpdateChange(ctx context.Context, obj, upd interface{}) error {
	r.logger.Log("action", "start", "component", "operatorkit", "function", "ApplyUpdatePatch")

	err := r.resource.ApplyUpdateChange(ctx, obj, upd)
	if err != nil {
		r.logger.Log("action", "error", "component", "operatorkit", "function", "ApplyUpdatePatch")
		return microerror.Mask(err)
	}

	r.logger.Log("action", "end", "component", "operatorkit", "function", "ApplyUpdatePatch")

	return nil
}

func (r *Resource) Underlying() framework.Resource {
	return r.resource.Underlying()
}
