package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"

	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
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
	Logger    micrologger.Logger
	K8sClient k8sclient.Interface
	CRD       *apiextensionsv1beta1.CustomResourceDefinition
}

type Timestamp struct {
	logger    micrologger.Logger
	K8sClient k8sclient.Interface
	crd       *apiextensionsv1beta1.CustomResourceDefinition
}

func NewTimestamp(config TimestampConfig) (*Timestamp, error) {

	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	t := &Timestamp{
		logger:    config.Logger,
		K8sClient: config.K8sClient,
		crd:       config.CRD,
	}

	return t, nil
}

func (t *Timestamp) Collect(ch chan<- prometheus.Metric) error {

	m, err := meta.Accessor(t.crd)
	if err != nil {
		return microerror.Mask(err)
	}
	k, err := meta.TypeAccessor(t.crd)
	if err != nil {
		return microerror.Mask(err)
	}

	ch <- prometheus.MustNewConstMetric(
		creationTimestampDesc,
		prometheus.GaugeValue,
		float64(m.GetCreationTimestamp().Unix()),
		k.GetKind(),
		m.GetName(),
		m.GetNamespace(),
	)

	if m.GetDeletionTimestamp() != nil {

		ch <- prometheus.MustNewConstMetric(
			deletionTimestampDesc,
			prometheus.GaugeValue,
			float64(m.GetDeletionTimestamp().Unix()),
			k.GetKind(),
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
