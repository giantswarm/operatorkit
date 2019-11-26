package collector

import (
	"fmt"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	creationTimestampDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName("operatorkit", "informer", "creation_timestamp"),
		"CreationTimestamp of watched runtime objects.",
		[]string{
			"kind",
			"name",
			"namespace",
		},
		nil,
	)
	deletionTimestampDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName("operatorkit", "informer", "deletion_timestamp"),
		"DeletionTimestamp of watched runtime objects.",
		[]string{
			"kind",
			"name",
			"namespace",
		},
		nil,
	)
)

type TimestampConfig struct {
	Logger  micrologger.Logger
	Watcher Watcher

	ListOptions metav1.ListOptions
}

type Timestamp struct {
	logger  micrologger.Logger
	watcher Watcher

	listOptions metav1.ListOptions
}

func NewTimestamp(config TimestampConfig) (*Timestamp, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Watcher == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Watcher must not be empty", config)
	}

	t := &Timestamp{
		logger:  config.Logger,
		watcher: config.Watcher,

		listOptions: config.ListOptions,
	}

	return t, nil
}

func (t *Timestamp) Collect(ch chan<- prometheus.Metric) error {
	watcher, err := t.watcher.Watch(t.listOptions)
	if err != nil {
		return microerror.Mask(err)
	}

	defer watcher.Stop()

	alreadyEmitted := map[string]struct{}{}

	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				continue
			}

			m, err := meta.Accessor(event.Object)
			if err != nil {
				return microerror.Mask(err)
			}
			t, err := meta.TypeAccessor(event.Object)
			if err != nil {
				return microerror.Mask(err)
			}

			// In the case we have multiple events for the same runtime object, we
			// track which metrics we have emitted already, simply by indexing their
			// labels as map key, to not emit the same metric for the same runtime
			// object twice, which is an error for Prometheus.
			{
				k := fmt.Sprintf("creation-%s-%s-%s", t.GetKind(), m.GetName(), m.GetNamespace())
				_, ok := alreadyEmitted[k]
				if ok {
					continue
				}
				alreadyEmitted[k] = struct{}{}

				ch <- prometheus.MustNewConstMetric(
					creationTimestampDesc,
					prometheus.GaugeValue,
					float64(m.GetCreationTimestamp().Unix()),
					t.GetKind(),
					m.GetName(),
					m.GetNamespace(),
				)
			}

			// We do the same verification with delete events as with create and
			// update events above in order to prevent emitting duplicated metrics.
			if m.GetDeletionTimestamp() != nil {
				k := fmt.Sprintf("deletion-%s-%s-%s", t.GetKind(), m.GetName(), m.GetNamespace())
				_, ok := alreadyEmitted[k]
				if ok {
					continue
				}
				alreadyEmitted[k] = struct{}{}

				ch <- prometheus.MustNewConstMetric(
					deletionTimestampDesc,
					prometheus.GaugeValue,
					float64(m.GetDeletionTimestamp().Unix()),
					t.GetKind(),
					m.GetName(),
					m.GetNamespace(),
				)
			}

		case <-time.After(time.Second):
			return nil
		}
	}
}

func (t *Timestamp) Describe(ch chan<- *prometheus.Desc) error {
	ch <- creationTimestampDesc
	ch <- deletionTimestampDesc

	return nil
}
