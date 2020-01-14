package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Interface defines how a collector implementation should look like.
type Interface interface {
	// Collect should align with the monitoring system's implementation
	// requirements. In this case Prometheus. See also
	// https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector.
	// The difference here is that this specific interface provides additional
	// error handling capabilities.
	Collect(ch chan<- prometheus.Metric) error
	// Describe should align with the monitoring system's implementation
	// requirements. In this case Prometheus. See also
	// https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector.
	// The difference here is that this specific interface provides additional
	// error handling capabilities.
	Describe(ch chan<- *prometheus.Desc) error
}
