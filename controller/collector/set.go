package collector

import (
	"github.com/giantswarm/exporterkit/collector"
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

type SetConfig struct {
	Logger    micrologger.Logger
	K8sClient k8sclient.Interface
}

// Set is basically only a wrapper for the informer's collector implementations.
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
			Logger:    config.Logger,
			K8sClient: config.K8sClient,
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
