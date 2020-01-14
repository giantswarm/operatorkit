package collector

import (
	"context"

	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

var (
	creationTimestampDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName("operatorkit", "controller", "creation_timestamp"),
		"CreationTimestamp of watched runtime objects.",
		[]string{
			"kind",
			"name",
			"namespace",
		},
		nil,
	)
	deletionTimestampDesc *prometheus.Desc = prometheus.NewDesc(
		prometheus.BuildFQName("operatorkit", "controller", "deletion_timestamp"),
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
	gvk.Kind += "List"
	list := &unstructured.UnstructuredList{}
	list.SetGroupVersionKind(gvk)

	err = t.k8sClient.CtrlClient().List(ctx, list)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, object := range list.Items {
		ch <- prometheus.MustNewConstMetric(
			creationTimestampDesc,
			prometheus.GaugeValue,
			float64(object.GetCreationTimestamp().Unix()),
			object.GetKind(),
			object.GetName(),
			object.GetNamespace(),
		)

		if object.GetDeletionTimestamp() != nil {
			ch <- prometheus.MustNewConstMetric(
				deletionTimestampDesc,
				prometheus.GaugeValue,
				float64(object.GetDeletionTimestamp().Unix()),
				object.GetKind(),
				object.GetName(),
				object.GetNamespace(),
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
