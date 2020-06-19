package sentry

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

type Config struct {
	Dsn    string
	Logger micrologger.Logger
}

type Service struct {
	enabled bool
	logger  micrologger.Logger
}

func New(config Config) (*Service, error) {
	disabled := Service{enabled: false}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Dsn == "" {
		config.Logger.Log("level", "warning", "Sentry DSN is not set.")
		return &disabled, nil
	}
	config.Logger.Log("level", "debug", fmt.Sprintf("Setting up Sentry with DSN %s", config.Dsn))
	err := sentry.Init(sentry.ClientOptions{
		Dsn: config.Dsn,
	})
	if err != nil {
		return &disabled, microerror.Mask(err)
	}

	svc := Service{
		enabled: true,
		logger:  config.Logger,
	}

	return &svc, nil
}

func (s *Service) Capture(ctx context.Context, err error) {
	if s.enabled {
		sentry.CaptureException(err)
	}
}
