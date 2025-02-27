package collector

import (
	"context"
	"fmt"

	"github.com/giantswarm/k8sclient/v8/pkg/k8sclient"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

type TimestampConfig struct {
	Logger               micrologger.Logger
	K8sClient            k8sclient.Interface
	NewRuntimeObjectFunc func() client.Object
	Selector             labels.Selector

	Controller string
}

type Timestamp struct {
	logger               micrologger.Logger
	ctrlClient           client.Client
	newRuntimeObjectFunc func() client.Object
	selector             labels.Selector
	scheme               *runtime.Scheme

	controller string
}

func NewTimestamp(config TimestampConfig) (*Timestamp, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Controller == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Controller must not be empty", config)
	}
	if config.NewRuntimeObjectFunc == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.NewRuntimeObjectFunc must not be empty", config)
	}
	if config.Selector == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Selector must not be empty", config)
	}

	t := &Timestamp{
		logger:               config.Logger,
		ctrlClient:           config.K8sClient.CtrlClient(),
		scheme:               config.K8sClient.Scheme(),
		newRuntimeObjectFunc: config.NewRuntimeObjectFunc,
		selector:             config.Selector,

		controller: config.Controller,
	}

	return t, nil
}

func (t *Timestamp) Collect(ch chan<- prometheus.Metric) error {
	var list unstructured.UnstructuredList
	{
		gvk, err := apiutil.GVKForObject(t.newRuntimeObjectFunc(), t.scheme)
		if err != nil {
			return microerror.Mask(err)
		}
		gvk.Kind = fmt.Sprintf("%sList", gvk.Kind)
		list.SetGroupVersionKind(gvk)
	}

	err := t.ctrlClient.List(context.Background(), &list, &client.ListOptions{
		LabelSelector: t.selector,
	})
	if err != nil {
		return microerror.Mask(err)
	}

	for _, object := range list.Items {
		ch <- prometheus.MustNewConstMetric(
			t.creationTimestampDesc(),
			prometheus.GaugeValue,
			float64(object.GetCreationTimestamp().Unix()),
			object.GetKind(),
			object.GetName(),
			object.GetNamespace(),
		)

		if object.GetDeletionTimestamp() != nil {
			ch <- prometheus.MustNewConstMetric(
				t.deletionTimestampDesc(),
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
	ch <- t.creationTimestampDesc()
	ch <- t.deletionTimestampDesc()

	return nil
}

// creationTimestampDesc must use the controller name as contant labels in order
// to keep the metrics unique for Prometheus registration.
func (t *Timestamp) creationTimestampDesc() *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName("operatorkit", "controller", "creation_timestamp"),
		"CreationTimestamp of watched runtime objects.",
		[]string{
			"kind",
			"name",
			"namespace",
		},
		map[string]string{
			"controller": t.controller,
		},
	)
}

// deletionTimestampDesc must use the controller name as contant labels in order
// to keep the metrics unique for Prometheus registration.
func (t *Timestamp) deletionTimestampDesc() *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName("operatorkit", "controller", "deletion_timestamp"),
		"DeletionTimestamp of watched runtime objects.",
		[]string{
			"kind",
			"name",
			"namespace",
		},
		map[string]string{
			"controller": t.controller,
		},
	)
}
