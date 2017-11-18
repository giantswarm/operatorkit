package framework

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	frameworkHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "operatorkit",
			Subsystem: "framework",
			Name:      "event",
			Help:      "Histogram for events within the operatorkit framework.",
		},
		[]string{"event"},
	)
)

func init() {
	prometheus.MustRegister(frameworkHistogram)
}
