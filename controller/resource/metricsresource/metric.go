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
		[]string{"resource", "operation"},
	)

	operationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "operation",
			Help:      "Time taken to process a single reconciliation operation.",
			Buckets:   []float64{1, 2, 5, 10, 20, 40, 60, 120},
		},
		[]string{"resource", "operation"},
	)
)

func init() {
	prometheus.MustRegister(operationCounter)
	prometheus.MustRegister(operationHistogram)
}
