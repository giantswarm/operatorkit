package informer

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/api/meta"
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

	watcher, err := i.watcher.Watch(i.listOptions)
	if err != nil {
		i.logger.Log("level", "error", "function", "Collect", "message", "could not start watch", "stack", fmt.Sprintf("%#v", err))
		return
	}

	defer watcher.Stop()

	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if ok {
				m, err := meta.Accessor(event.Object)
				if err != nil {
					i.logger.Log("level", "error", "function", "Collect", "message", "could not get accessor for object", "stack", fmt.Sprintf("%#v", err))
					break
				}
				if m.GetDeletionTimestamp() != nil {
					ch <- prometheus.MustNewConstMetric(
						description,
						prometheus.GaugeValue,
						float64(m.GetDeletionTimestamp().Unix()),
						m.GetName(),
						m.GetNamespace(),
					)
				}
			}
		case <-time.After(time.Second):
			i.logger.Log("level", "debug", "message", "finished collecting metrics")
			return
		}
	}
}
