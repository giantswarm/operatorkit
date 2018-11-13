package collector

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/client_golang/prometheus"
)

type SetConfig struct {
	Collectors []Interface
	Logger     micrologger.Logger
}

type Set struct {
	collectors []Interface
	logger     micrologger.Logger

	bootedCounter uint32
	mutex         sync.Mutex
}

func NewSet(config SetConfig) (*Set, error) {
	if len(config.Collectors) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.Collectors must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	s := &Set{
		collectors: config.Collectors,
		logger:     config.Logger,

		bootedCounter: 0,
		mutex:         sync.Mutex{},
	}

	return s, nil
}

func (s *Set) Boot(ctx context.Context) error {
	s.logger.Log("level", "debug", "message", "booting collector")

	if s.isBooted() {
		return nil
	}

	{
		s.logger.LogCtx(ctx, "level", "debug", "message", "registering collector")

		err := prometheus.Register(s)
		if IsAlreadyRegisteredError(err) {
			s.logger.LogCtx(ctx, "level", "debug", "message", "collector already registered")
		} else if err != nil {
			s.logger.Log("level", "error", "message", "failed registering collector", "stack", fmt.Sprintf("%#v", microerror.Mask(err)))
		} else {
			s.logger.LogCtx(ctx, "level", "debug", "message", "registered collector")
		}
	}

	s.logger.Log("level", "debug", "message", "booted collector")

	return nil
}

func (s *Set) Collect(ch chan<- prometheus.Metric) {
	s.logger.Log("level", "debug", "message", "collecting metrics")

	for _, c := range s.collectors {
		err := c.Collect(ch)
		if err != nil {
			s.logger.Log("level", "error", "message", "failed collecting metrics", "stack", fmt.Sprintf("%#v", microerror.Mask(err)))
		}
	}

	s.logger.Log("level", "debug", "message", "collected metrics")
}

func (s *Set) Describe(ch chan<- *prometheus.Desc) {
	s.logger.Log("level", "debug", "message", "describing metrics")

	for _, c := range s.collectors {
		err := c.Describe(ch)
		if err != nil {
			s.logger.Log("level", "error", "message", "failed describing metrics", "stack", fmt.Sprintf("%#v", microerror.Mask(err)))
		}
	}

	s.logger.Log("level", "debug", "message", "described metrics")
}

func (s *Set) isBooted() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if atomic.LoadUint32(&s.bootedCounter) == 1 {
		return true
	}

	atomic.StoreUint32(&s.bootedCounter, 1)

	return false
}
