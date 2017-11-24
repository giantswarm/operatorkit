package framework

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	frameworkEventSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "operatorkit",
			Subsystem:  "framework",
			Name:       "event",
			Help:       "Summary for events within the operatorkit framework.",
			Objectives: map[float64]float64{},
		},
		[]string{"event"},
	)
)

func init() {
	prometheus.MustRegister(frameworkEventSummary)
}
