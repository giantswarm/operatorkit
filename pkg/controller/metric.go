package controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	PrometheusNamespace = "operatorkit"
	PrometheusSubsystem = "controller"
)

var (
	errorCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "error_total",
			Help:      "Number of reconciliation errors.",
		},
		[]string{"cr_name", "cluster_id", "release_version"},
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
)

func init() {
	prometheus.MustRegister(errorCounterVec)
	prometheus.MustRegister(eventHistogram)
}

func errorMetricLabels(obj runtime.Object) prometheus.Labels {
	labels := make(map[string]string)

	t, err := meta.TypeAccessor(obj)
	if err != nil {
		return labels
	}

	m, err := meta.Accessor(obj)
	if err != nil {
		return labels
	}

	labels["cr_kind"] = t.GetKind()
	labels["cr_name"] = m.GetName()

	if v, exists := m.GetLabels()["cluster.x-k8s.io/cluster-name"]; exists {
		labels["cluster_id"] = v
	} else if v, exists := m.GetLabels()["giantswarm.io/cluster"]; exists {
		labels["cluster_id"] = v
	} else {
		labels["cluster_id"] = ""
	}

	if v, exists := m.GetLabels()["release.giantswarm.io/version"]; exists {
		labels["release_version"] = v
	} else {
		labels["release_version"] = ""
	}

	return labels
}
