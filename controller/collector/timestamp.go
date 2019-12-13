package collector

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/api/meta"
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
	Logger micrologger.Logger
}

type Timestamp struct {
	logger micrologger.Logger
}

func NewTimestamp(config TimestampConfig) (*Timestamp, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	t := &Timestamp{
		logger: config.Logger,
	}

	return t, nil
}

func (t *Timestamp) Collect(ch chan<- prometheus.Metric) error {

	ch <- prometheus.MustNewConstMetric(
		creationTimestampDesc,
		prometheus.GaugeValue,
		float64(m.GetCreationTimestamp().Unix()),
		t.GetKind(),
		m.GetName(),
		m.GetNamespace(),
	)

	if m.GetDeletionTimestamp() != nil {

		ch <- prometheus.MustNewConstMetric(
			deletionTimestampDesc,
			prometheus.GaugeValue,
			float64(m.GetDeletionTimestamp().Unix()),
			t.GetKind(),
			m.GetName(),
			m.GetNamespace(),
		)
	}

	return nil

}

func (t *Timestamp) Describe(ch chan<- *prometheus.Desc) error {
	ch <- creationTimestampDesc
	ch <- deletionTimestampDesc

	return nil
}
