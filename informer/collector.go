package informer

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	clientmodel "github.com/prometheus/client_model/go"
	"k8s.io/apimachinery/pkg/api/meta"
)

var (
	CreationTimestampDescription *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName("operatorkit", "informer", "creation_timestamp"),
		"CreationTimestamp of watched objects.",
		[]string{
			"kind",
			"name",
			"namespace",
		},
		nil,
	)
	DeletionTimestampDescription *prometheus.Desc = prometheus.NewDesc(
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
	ch <- CreationTimestampDescription
	ch <- DeletionTimestampDescription
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

	// In the case we have multiple events for the same resource,
	// we track which metrics we have emitted to not emit the same metric for the
	// same resource twice, which is an error for Prometheus.
	emittedMetrics := []prometheus.Metric{}

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

			fmt.Printf("\n")
			fmt.Printf("----------------\n")
			fmt.Printf("event: %v %#v\n", m.GetName(), event)

			creationMetric := prometheus.MustNewConstMetric(
				CreationTimestampDescription,
				prometheus.GaugeValue,
				float64(m.GetCreationTimestamp().Unix()),
				t.GetKind(),
				m.GetName(),
				m.GetNamespace(),
			)

			if !metricEmitted(emittedMetrics, creationMetric) {
				ch <- creationMetric
				emittedMetrics = append(emittedMetrics, creationMetric)
			}

			if m.GetDeletionTimestamp() != nil {
				deletionMetric := prometheus.MustNewConstMetric(
					DeletionTimestampDescription,
					prometheus.GaugeValue,
					float64(m.GetDeletionTimestamp().Unix()),
					t.GetKind(),
					m.GetName(),
					m.GetNamespace(),
				)
				if !metricEmitted(emittedMetrics, deletionMetric) {
					ch <- deletionMetric
					emittedMetrics = append(emittedMetrics, deletionMetric)
				}
			}

			fmt.Printf("----------------\n")

		case <-time.After(time.Second):
			i.logger.Log("level", "debug", "message", "finished collecting metrics")
			return
		}
	}
}

func metricEmitted(emittedMetrics []prometheus.Metric, metric prometheus.Metric) bool {
	modelMetric := clientmodel.Metric{}
	metric.Write(&modelMetric)

	fmt.Printf("desc: %v\n", *metric.Desc())
	fmt.Printf("labels: %v\n", modelMetric.GetLabel())

	fmt.Printf("%v emitted metrics\n", len(emittedMetrics))

	for _, emittedMetric := range emittedMetrics {
		modelEmittedMetric := clientmodel.Metric{}
		emittedMetric.Write(&modelEmittedMetric)

		fmt.Printf("emitted desc: %v\n", *emittedMetric.Desc())
		fmt.Printf("emitted labels: %v\n", modelEmittedMetric.GetLabel())

		fmt.Printf("metric string: %v\n", modelMetric.String())
		fmt.Printf("emitted metric string: %v\n", modelEmittedMetric.String())

		if metric.Desc().String() == emittedMetric.Desc().String() && modelEmittedMetric.String() == modelMetric.String() {
			fmt.Printf("this is a metric that has been emitted already\n")
			return true
		}
	}

	fmt.Printf("this is a metric that has not been emitted already\n")
	return false
}
