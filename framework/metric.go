package framework

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	PrometheusNamespace = "operatorkit"
	PrometheusSubsystem = "framework"
)

var (
	frameworkHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "event",
			Help:      "Histogram for events within the operatorkit framework.",
		},
		[]string{"event"},
	)
)

func init() {
	prometheus.MustRegister(frameworkHistogram)
}
