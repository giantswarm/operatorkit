package informer

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	testDesc      *prometheus.Desc = prometheus.NewDesc("test", "test desc", []string{}, nil)
	labelTestDesc *prometheus.Desc = prometheus.NewDesc("label_test", "labelled test desc", []string{"a", "b"}, nil)
)

// Test_Informer_Metric_Emitted tests that the metricEmitted function
// correctly tracks which metrics have already been emitted.
func Test_Informer_Metric_Emitted(t *testing.T) {
	testCases := []struct {
		name           string
		emittedMetrics []prometheus.Metric
		metric         prometheus.Metric
		metricEmitted  bool
	}{
		{
			name:           "Test that with no metrics emitted a metric is not considered emitted already",
			emittedMetrics: []prometheus.Metric{},
			metric:         prometheus.MustNewConstMetric(testDesc, prometheus.GaugeValue, float64(1)),
			metricEmitted:  false,
		},

		{
			name: "Test that wih a metric emitted the same metric is considered emitted already",
			emittedMetrics: []prometheus.Metric{
				prometheus.MustNewConstMetric(testDesc, prometheus.GaugeValue, float64(1)),
			},
			metric:        prometheus.MustNewConstMetric(testDesc, prometheus.GaugeValue, float64(1)),
			metricEmitted: true,
		},

		{
			name: "Test X",
			emittedMetrics: []prometheus.Metric{
				prometheus.MustNewConstMetric(labelTestDesc, prometheus.GaugeValue, float64(1), "a", "b"),
				prometheus.MustNewConstMetric(labelTestDesc, prometheus.GaugeValue, float64(1), "c", "d"),
			},
			metric:        prometheus.MustNewConstMetric(labelTestDesc, prometheus.GaugeValue, float64(1), "c", "d"),
			metricEmitted: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			emitted := metricEmitted(tc.emittedMetrics, tc.metric)
			if emitted != tc.metricEmitted {
				t.Fatal("expected", tc.metricEmitted, "got", emitted)
			}
		})
	}
}
