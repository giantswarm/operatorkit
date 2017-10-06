package metricsresource

import "github.com/prometheus/client_golang/prometheus"

var (
	errorTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "operatorkit",
			Subsystem: "framework",
			Name:      "error_total",
			Help:      "Number of operation errors.",
		},
		[]string{"service", "resource", "operation"},
	)

	operationDuration = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "operatorkit",
			Subsystem: "framework",
			Name:      "operation_duration_milliseconds",
			Help:      "Time taken to process a single reconciliation operation.",
		},
		[]string{"service", "resource", "operation"},
	)

	operationTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "operatorkit",
			Subsystem: "framework",
			Name:      "operation_total",
			Help:      "Number of processed reconciliation operations.",
		},
		[]string{"service", "resource", "operation"},
	)
)

func init() {
	prometheus.MustRegister(errorTotal)
	prometheus.MustRegister(operationDuration)
	prometheus.MustRegister(operationTotal)
}
