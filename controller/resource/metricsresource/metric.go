package metricsresource

import "github.com/prometheus/client_golang/prometheus"

const (
	PrometheusNamespace = "operatorkit"
	PrometheusSubsystem = "controller"
)

var (
	operationCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "operation_total",
			Help:      "Number of processed reconciliation operations.",
		},
		[]string{"service", "resource", "operation"},
	)

	operationErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "operation_error_total",
			Help:      "Number of operation errors.",
		},
		[]string{"service", "resource", "operation"},
	)

	operationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "operation",
			Help:      "Time taken to process a single reconciliation operation.",
		},
		[]string{"service", "resource", "operation"},
	)
)

func init() {
	prometheus.MustRegister(operationCounter)
	prometheus.MustRegister(operationErrorCounter)
	prometheus.MustRegister(operationHistogram)
}
