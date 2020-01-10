package collector

import (
	"context"

	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
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
	Logger               micrologger.Logger
	K8sClient            k8sclient.Interface
	NewRuntimeObjectFunc func() runtime.Object
}

type Timestamp struct {
	logger               micrologger.Logger
	k8sClient            k8sclient.Interface
	newRuntimeObjectFunc func() runtime.Object
}

func NewTimestamp(config TimestampConfig) (*Timestamp, error) {

	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	t := &Timestamp{
		logger:               config.Logger,
		k8sClient:            config.K8sClient,
		newRuntimeObjectFunc: config.NewRuntimeObjectFunc,
	}

	return t, nil
}

func (t *Timestamp) Collect(ch chan<- prometheus.Metric) error {

	ctx := context.Background()

	gvk, err := apiutil.GVKForObject(t.newRuntimeObjectFunc(), t.k8sClient.Scheme())
	if err != nil {
		return microerror.Mask(err)
	}
	list := &unstructured.UnstructuredList{}
	list.SetGroupVersionKind(gvk)
	err = t.k8sClient.CtrlClient().List(ctx, list)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, object := range list.Items {

		m, err := meta.Accessor(object)
		if err != nil {
			return microerror.Mask(err)
		}
		k, err := meta.TypeAccessor(object)
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
	}

	return nil

}

func (t *Timestamp) Describe(ch chan<- *prometheus.Desc) error {
	ch <- creationTimestampDesc
	ch <- deletionTimestampDesc

	return nil
}
