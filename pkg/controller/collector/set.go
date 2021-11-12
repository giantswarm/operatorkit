package collector

import (
	"github.com/giantswarm/exporterkit/collector"
	"github.com/giantswarm/k8sclient/v6/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SetConfig struct {
	Logger     micrologger.Logger
	K8sClient  k8sclient.Interface
	Controller string

	NewRuntimeObjectFunc func() client.Object
	Selector             labels.Selector
}

// Set is basically only a wrapper for the collector implementations.
// It eases the initialization and prevents some weird import mess so we do not
// have to alias packages.
type Set struct {
	*collector.Set
}

func NewSet(config SetConfig) (*Set, error) {
	var err error

	var timestampCollector *Timestamp
	{
		c := TimestampConfig{
			Logger:     config.Logger,
			K8sClient:  config.K8sClient,
			Controller: config.Controller,

			NewRuntimeObjectFunc: config.NewRuntimeObjectFunc,
			Selector:             config.Selector,
		}

		timestampCollector, err = NewTimestamp(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var collectorSet *collector.Set
	{
		c := collector.SetConfig{
			Collectors: []collector.Interface{
				timestampCollector,
			},
			Logger: config.Logger,
		}

		collectorSet, err = collector.NewSet(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	s := &Set{
		Set: collectorSet,
	}

	return s, nil
}
