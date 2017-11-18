package tpr

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	prometheusNamespace = "operatorkit"
)

var (
	tpoCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: prometheusNamespace,
			Name:      "tpo_total",
			Help:      "Number of TPOs in the cluster, labelled by kind, api version, name, and group.",
		},
		[]string{"kind", "api_version", "name", "group"},
	)
)

func init() {
	prometheus.MustRegister(tpoCount)
}
