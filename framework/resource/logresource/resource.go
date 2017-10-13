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
		return nil, microerror.Mask(err)
	}

	r.logger.Log("action", "end", "component", "operatorkit", "function", "GetCurrentState")

	return v, nil
}

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	r.logger.Log("action", "start", "component", "operatorkit", "function", "GetDesiredState")

	v, err := r.resource.GetDesiredState(ctx, obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("action", "end", "component", "operatorkit", "function", "GetDesiredState")

	return v, nil
}

func (r *Resource) GetCreateState(ctx context.Context, obj, cur, des interface{}) (interface{}, error) {
	r.logger.Log("action", "start", "component", "operatorkit", "function", "GetCreateState")

	v, err := r.resource.GetCreateState(ctx, obj, cur, des)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("action", "end", "component", "operatorkit", "function", "GetCreateState")

	return v, nil
}

func (r *Resource) GetDeleteState(ctx context.Context, obj, cur, des interface{}) (interface{}, error) {
	r.logger.Log("action", "start", "component", "operatorkit", "function", "GetDeleteState")

	v, err := r.resource.GetDeleteState(ctx, obj, cur, des)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("action", "end", "component", "operatorkit", "function", "GetDeleteState")

	return v, nil
}

func (r *Resource) GetUpdateState(ctx context.Context, obj, cur, des interface{}) (interface{}, interface{}, interface{}, error) {
	r.logger.Log("action", "start", "component", "operatorkit", "function", "GetUpdateState")

	createState, deleteState, updateState, err := r.resource.GetUpdateState(ctx, obj, cur, des)
	if err != nil {
		return nil, nil, nil, microerror.Mask(err)
	}

	r.logger.Log("action", "end", "component", "operatorkit", "function", "GetUpdateState")

	return createState, deleteState, updateState, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) ProcessCreateState(ctx context.Context, obj, cre interface{}) error {
	r.logger.Log("action", "start", "component", "operatorkit", "function", "ProcessCreateState")

	err := r.resource.ProcessCreateState(ctx, obj, cre)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Log("action", "end", "component", "operatorkit", "function", "ProcessCreateState")

	return nil
}

func (r *Resource) ProcessDeleteState(ctx context.Context, obj, del interface{}) error {
	r.logger.Log("action", "start", "component", "operatorkit", "function", "ProcessDeleteState")

	err := r.resource.ProcessDeleteState(ctx, obj, del)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Log("action", "end", "component", "operatorkit", "function", "ProcessDeleteState")

	return nil
}

func (r *Resource) ProcessUpdateState(ctx context.Context, obj, upd interface{}) error {
	r.logger.Log("action", "start", "component", "operatorkit", "function", "ProcessUpdateState")

	err := r.resource.ProcessUpdateState(ctx, obj, upd)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Log("action", "end", "component", "operatorkit", "function", "ProcessUpdateState")

	return nil
}

func (r *Resource) Underlying() framework.Resource {
	return r.resource.Underlying()
}
