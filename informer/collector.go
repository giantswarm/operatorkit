package informer

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/watch"
)

var (
	description *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName("operatorkit", "informer", "deletion_timestamp"),
		"DeletionTimestamp of watched objects.",
		[]string{
			"exported_name",
			"exported_namespace",
		},
		nil,
	)
)

// Describe is used to describe metrics which are exported to Prometheus.
func (i *Informer) Describe(ch chan<- *prometheus.Desc) {
	ch <- description
}

// Collect is called by the Prometheus registry when collecting metrics from
// the informer.
func (i *Informer) Collect(ch chan<- prometheus.Metric) {
	eventChan := make(chan watch.Event)
	ctx := context.Background()
	go func() {
		err := i.fillCache(ctx, eventChan)
		if err != nil {
			return
		}
		close(eventChan)
	}()

	for e := range eventChan {
		m, err := meta.Accessor(e.Object)
		if err != nil || m.GetDeletionTimestamp() == nil {
			return
		}
		ch <- prometheus.MustNewConstMetric(
			description,
			prometheus.GaugeValue,
			float64(m.GetDeletionTimestamp().Unix()),
			m.GetName(),
			m.GetNamespace(),
		)
	}
}
