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
			"deletion_timestamp",
		},
		nil,
	)
)

func (i *Informer) Describe(ch chan<- *prometheus.Desc) {

	ch <- description

}

func (i *Informer) Collect(ch chan<- prometheus.Metric) {
	eventChan := make(chan watch.Event)
	ctx := context.Background()
	go func() {
		err := i.fillCache(ctx, eventChan)
		if err != nil {
			return
		}
		go func() {
			for {
				select {
				case <-ctx.Done():
					close(eventChan)
					return
				}
			}
		}()
	}()

	for e := range eventChan {
		m, err := meta.Accessor(e.Object)
		if err != nil {
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
