package informer

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	PrometheusNamespace = "operatorkit"
	PrometheusSubsystem = "informer"
)

var (
	cacheLastUpdatedGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "cache_last_updated",
			Help:      "A gauge metric expressing the time of the last cache update.",
		},
	)

	cacheSizeGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "cache_size",
			Help:      "A gauge metric expressing the number of events being cached in memory.",
		},
	)

	watcherCloseCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "watcher_close_total",
			Help:      "A counter metric expressing the number of events when the watcher's watch got closed.",
		},
	)

	watchEventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: PrometheusNamespace,
			Subsystem: PrometheusSubsystem,
			Name:      "watch_event_total",
			Help:      "A counter metric expressing the number of event kinds happening.",
		},
		[]string{"kind"},
	)
)

func init() {
	prometheus.MustRegister(cacheLastUpdatedGauge)
	prometheus.MustRegister(cacheSizeGauge)
	prometheus.MustRegister(watcherCloseCounter)
	prometheus.MustRegister(watchEventCounter)
}
