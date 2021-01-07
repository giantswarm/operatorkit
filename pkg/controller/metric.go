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
	lastReconciledGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "last_reconciled",
			Help:      "Last reconciled Timestamp of watched runtime objects.",
		},
		[]string{"controller"},
	)
)

func init() {
	prometheus.MustRegister(errorGauge)
	prometheus.MustRegister(eventHistogram)
	prometheus.MustRegister(lastReconciledGauge)
}
