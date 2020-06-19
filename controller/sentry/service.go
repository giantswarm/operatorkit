package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/giantswarm/microerror"
)

type Config struct {
	Dsn string
}

type Service struct {
	enabled bool
}

func New(config Config) (*Service, error) {
	disabled := Service{enabled: false}
	if config.Dsn == "" {
		return &disabled, nil
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn: config.Dsn,
	})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	s := &Service{
		enabled: true,
	}

	return s, nil
}

func (s *Service) Capture(ctx context.Context, err error) {
	if s.enabled {
		sentry.CaptureException(err)
	}
}
