package controller

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	PrometheusNamespace = "operatorkit"
	PrometheusSubsystem = "controller"
)

var (
	errorGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "error_total",
			Help:      "Number of reconciliation errors.",
		},
	)
	eventHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "event",
			Help:      "Histogram for events within the operatorkit controller.",
		},
		[]string{"event"},
	)
	creationTimestampGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "creation_timestamp",
			Help:      "CreationTimestamp of watched runtime objects.",
		},
		[]string{
			"kind",
			"name",
			"namespace",
		},
	)
	deletionTimestampGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "deletion_timestamp",
			Help:      "DeletionTimestamp of watched runtime objects.",
		},
		[]string{
			"kind",
			"name",
			"namespace",
		},
	)
)

func init() {
	prometheus.MustRegister(errorGauge)
	prometheus.MustRegister(eventHistogram)
	prometheus.MustRegister(creationTimestampGauge)
	prometheus.MustRegister(deletionTimestampGauge)
}
