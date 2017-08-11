package logresource

import (
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

func (r *Resource) GetCurrentState(obj interface{}) (interface{}, error) {
	r.logger.Log("debug", PreOperationMessage, "operation", "GetCurrentState")

	v, err := r.resource.GetCurrentState(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("debug", PostOperationMessage, "operation", "GetCurrentState")

	return v, nil
}

func (r *Resource) GetDesiredState(obj interface{}) (interface{}, error) {
	r.logger.Log("debug", PreOperationMessage, "operation", "GetDesiredState")

	v, err := r.resource.GetDesiredState(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("debug", PostOperationMessage, "operation", "GetDesiredState")

	return v, nil
}

func (r *Resource) GetCreateState(obj, cur, des interface{}) (interface{}, error) {
	r.logger.Log("debug", PreOperationMessage, "operation", "GetCreateState")

	v, err := r.resource.GetCreateState(obj, cur, des)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("debug", PostOperationMessage, "operation", "GetCreateState")

	return v, nil
}

func (r *Resource) GetDeleteState(obj, cur, des interface{}) (interface{}, error) {
	r.logger.Log("debug", PreOperationMessage, "operation", "GetDeleteState")

	v, err := r.resource.GetDeleteState(obj, cur, des)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r.logger.Log("debug", PostOperationMessage, "operation", "GetDeleteState")

	return v, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) ProcessCreateState(obj, cre interface{}) error {
	r.logger.Log("debug", PreOperationMessage, "operation", "ProcessCreateState")

	err := r.resource.ProcessCreateState(obj, cre)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Log("debug", PostOperationMessage, "operation", "ProcessCreateState")

	return nil
}

func (r *Resource) ProcessDeleteState(obj, del interface{}) error {
	r.logger.Log("debug", PreOperationMessage, "operation", "ProcessDeleteState")

	err := r.resource.ProcessDeleteState(obj, del)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Log("debug", PostOperationMessage, "operation", "ProcessDeleteState")

	return nil
}

func (r *Resource) Underlying() framework.Resource {
	return r.resource.Underlying()
}
