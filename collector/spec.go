package collector

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

// Interface defines how a collector implementation should look like.
type Interface interface {
	// Boot is used to initial and register the collector. Boot must be allowed to
	// be called multiple times and thus must be idempotent. Usual implementations
	// could make use of sync.Once. See also
	// https://godoc.org/github.com/prometheus/client_golang/prometheus#Register.
	Boot(ctx context.Context)
	// Collect should align with the monitoring system's implementation
	// requirements. In this case Prometheus. See also
	// https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector.
	Collect(ch chan<- prometheus.Metric)
	// Describe should align with the monitoring system's implementation
	// requirements. In this case Prometheus. See also
	// https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector.
	Describe(ch chan<- *prometheus.Desc)
}
