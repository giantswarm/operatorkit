package metricsresource

import (
	"context"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/giantswarm/operatorkit/framework"
)

const (
	// Name is the identifier of the resource.
	Name = "metrics"
)

// Config represents the configuration used to create a new metrics resource.
type Config struct {
	// Dependencies.
	Resource framework.Resource

	// Settings.

	// Namespace is the Prometheus namespace used to create new vectors. The user
	// has to provide unique namespaces and subsystems. If these settings are not
	// properly configured and reused the registration of the Prometheus vectors
	// fails with a panic.
	Namespace string
	// Subsystem is the Prometheus subsystem used to create new vectors. The user
	// has to provide unique namespaces and subsystems. If these settings are not
	// properly configured and reused the registration of the Prometheus vectors
	// fails with a panic.
	Subsystem string
}

// DefaultConfig provides a default configuration to create a new metrics
// resource by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Resource: nil,

		// Settings.
		Namespace: "",
		Subsystem: "",
	}
}

// New creates a new configured metrics resource.
func New(config Config) (*Resource, error) {
	// Dependencies.
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Resource must not be empty")
	}

	// Settings.
	if config.Namespace == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Namespace must not be empty")
	}
	if config.Subsystem == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Subsystem must not be empty")
	}

	var operationDuration *prometheus.GaugeVec
	var operationTotal *prometheus.CounterVec
	{
		operationDuration = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: toCamelCase(config.Namespace),
				Subsystem: toCamelCase(config.Subsystem),
				Name:      "operatorkit_framework_operation_duration_milliseconds",
				Help:      "Time taken to process a single reconciliation operation.",
			},
			[]string{"operation"},
		)
		operationTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: toCamelCase(config.Namespace),
				Subsystem: toCamelCase(config.Subsystem),
				Name:      "operatorkit_framework_operation_total",
				Help:      "Number of processed reconciliation operations.",
			},
			[]string{"operation"},
		)

		prometheus.MustRegister(operationDuration)
		prometheus.MustRegister(operationTotal)
	}

	newResource := &Resource{
		// Dependencies.
		resource: config.Resource,

		// Internals.
		operationDuration: operationDuration,
		operationTotal:    operationTotal,
	}

	return newResource, nil
}

type Resource struct {
	// Dependencies.
	resource framework.Resource

	// Internals.
	operationDuration *prometheus.GaugeVec
	operationTotal    *prometheus.CounterVec
}

func (r *Resource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	defer r.updateMetrics("GetCurrentState", time.Now())

	v, err := r.resource.GetCurrentState(ctx, obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *Resource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	defer r.updateMetrics("GetDesiredState", time.Now())

	v, err := r.resource.GetDesiredState(ctx, obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *Resource) GetCreateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	defer r.updateMetrics("GetCreateState", time.Now())

	v, err := r.resource.GetCreateState(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *Resource) GetDeleteState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	defer r.updateMetrics("GetDeleteState", time.Now())

	v, err := r.resource.GetDeleteState(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return v, nil
}

func (r *Resource) GetUpdateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, interface{}, interface{}, error) {
	defer r.updateMetrics("GetUpdateState", time.Now())

	createState, deleteState, updateState, err := r.resource.GetUpdateState(ctx, obj, currentState, desiredState)
	if err != nil {
		return nil, nil, nil, microerror.Mask(err)
	}

	return createState, deleteState, updateState, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) ProcessCreateState(ctx context.Context, obj, createState interface{}) error {
	defer r.updateMetrics("ProcessCreateState", time.Now())

	err := r.resource.ProcessCreateState(ctx, obj, createState)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) ProcessDeleteState(ctx context.Context, obj, deleteState interface{}) error {
	defer r.updateMetrics("ProcessDeleteState", time.Now())

	err := r.resource.ProcessDeleteState(ctx, obj, deleteState)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) ProcessUpdateState(ctx context.Context, obj, updateState interface{}) error {
	defer r.updateMetrics("ProcessUpdateState", time.Now())

	err := r.resource.ProcessUpdateState(ctx, obj, updateState)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) Underlying() framework.Resource {
	return r.resource.Underlying()
}

func (r *Resource) updateMetrics(operation string, startTime time.Time) {
	r.operationDuration.WithLabelValues(operation).Set(float64(time.Since(startTime) / time.Millisecond))
	r.operationTotal.WithLabelValues(operation).Inc()
}
