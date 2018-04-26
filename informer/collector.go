package informer

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/watch"
)

var (
	description *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName("operatorkit", "informer", "deletion_timestamp"),
		"DeletionTimestamp of watched objects.",
		[]string{
			"name",
			"namespace",
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
	i.logger.Log("level", "debug", "message", "start collecting metrics")

	eventChan := make(chan watch.Event)
	ctx := context.Background()
	go func() {
		err := i.fillCache(ctx, eventChan)
		if err != nil {
			i.logger.Log("level", "error", "function", "Collect", "message", "could not list objects from cache", "stack", fmt.Sprintf("%#v", err))
			return
		}
		close(eventChan)
	}()

	for e := range eventChan {
		m, err := meta.Accessor(e.Object)
		if err != nil {
			i.logger.Log("level", "error", "function", "Collect", "message", "could not get accessor for object", "stack", fmt.Sprintf("%#v", err))
			return
		}
		if m.GetDeletionTimestamp() == nil {
			return // fall through, deletionTimestamp is not set yet
		}
		ch <- prometheus.MustNewConstMetric(
			description,
			prometheus.GaugeValue,
			float64(m.GetDeletionTimestamp().Unix()),
			m.GetName(),
			m.GetNamespace(),
		)
	}

	i.logger.Log("level", "debug", "message", "finished collecting metrics")
}
