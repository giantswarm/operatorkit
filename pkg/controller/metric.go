package controller

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	PrometheusNamespace = "operatorkit"
	PrometheusSubsystem = "controller"
)

var (
	// ReconcileErrors is a prometheus counter metrics which holds the total
	// number of errors from the Reconciler.
	reconcileErrors = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: PrometheusNamespace,
		Subsystem: PrometheusSubsystem,
		Name:      "errors_total",
		Help:      "Total number of reconciliation errors per controller",
	}, []string{"controller"})

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
	prometheus.MustRegister(reconcileErrors)
	prometheus.MustRegister(eventHistogram)
	prometheus.MustRegister(lastReconciledGauge)
}
