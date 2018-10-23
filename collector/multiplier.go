package collector

import (
	"context"
	"sync"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	Collectors []Interface
	Logger     micrologger.Logger
}

type Multiplier struct {
	collectors []Interface
	logger     micrologger.Logger

	bootOnce sync.Once
}

func New(config Config) (*Multiplier, error) {
	if len(config.Collectors) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.Collectors must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	m := &Multiplier{
		collectors: config.Collectors,
		logger:     config.Logger,

		bootOnce: sync.Once{},
	}

	return m, nil
}

func (m *Multiplier) Boot(ctx context.Context) {
	m.logger.Log("level", "debug", "message", "booting collector")

	m.bootOnce.Do(func() {
		for _, c := range m.collectors {
			c.Boot(ctx)
		}
	})

	m.logger.Log("level", "debug", "message", "booted collector")
}

func (m *Multiplier) Collect(ch chan<- prometheus.Metric) {
	m.logger.Log("level", "debug", "message", "collecting metrics")

	for _, c := range m.collectors {
		c.Collect(ch)
	}

	m.logger.Log("level", "debug", "message", "collected metrics")
}

func (m *Multiplier) Describe(ch chan<- *prometheus.Desc) {
	m.logger.Log("level", "debug", "message", "describing metrics")

	for _, c := range m.collectors {
		c.Describe(ch)
	}

	m.logger.Log("level", "debug", "message", "described metrics")
}
