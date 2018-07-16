package informer

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/api/meta"
)

var (
	creationTimestampDescription *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName("operatorkit", "informer", "creation_timestamp"),
		"DeletionTimestamp of watched objects.",
		[]string{
			"kind",
			"name",
			"namespace",
		},
		nil,
	)
	deletionTimestampDescription *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName("operatorkit", "informer", "deletion_timestamp"),
		"DeletionTimestamp of watched objects.",
		[]string{
			"kind",
			"name",
			"namespace",
		},
		nil,
	)
)

// Describe is used to describe metrics which are exported to Prometheus.
func (i *Informer) Describe(ch chan<- *prometheus.Desc) {
	ch <- creationTimestampDescription
	ch <- deletionTimestampDescription
}

// Collect is called by the Prometheus registry when collecting metrics from
// the informer.
func (i *Informer) Collect(ch chan<- prometheus.Metric) {
	i.logger.Log("level", "debug", "message", "start collecting metrics")

	watcher, err := i.watcher.Watch(i.listOptions)
	if err != nil {
		i.logger.Log("level", "error", "message", "could not start watch", "stack", fmt.Sprintf("%#v", err))
		return
	}

	defer watcher.Stop()

	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				continue
			}

			m, err := meta.Accessor(event.Object)
			if err != nil {
				i.logger.Log("level", "error", "message", "could not get meta accessor for object", "stack", fmt.Sprintf("%#v", err))
				break
			}
			t, err := meta.TypeAccessor(event.Object)
			if err != nil {
				i.logger.Log("level", "error", "message", "could not get type accessor for object", "stack", fmt.Sprintf("%#v", err))
				break
			}

			ch <- prometheus.MustNewConstMetric(
				creationTimestampDescription,
				prometheus.GaugeValue,
				float64(m.GetCreationTimestamp().Unix()),
				t.GetKind(),
				m.GetName(),
				m.GetNamespace(),
			)

			if m.GetDeletionTimestamp() != nil {
				ch <- prometheus.MustNewConstMetric(
					deletionTimestampDescription,
					prometheus.GaugeValue,
					float64(m.GetDeletionTimestamp().Unix()),
					t.GetKind(),
					m.GetName(),
					m.GetNamespace(),
				)
			}
		case <-time.After(time.Second):
			i.logger.Log("level", "debug", "message", "finished collecting metrics")
			return
		}
	}
}
