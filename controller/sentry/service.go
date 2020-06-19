package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/giantswarm/microerror"
)

type Config struct {
	Dsn string
}

type Sentry struct {
	enabled bool
}

func New(config Config) (*Sentry, error) {
	disabled := Sentry{enabled: false}
	if config.Dsn == "" {
		return &disabled, nil
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn: config.Dsn,
	})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	s := &Sentry{
		enabled: true,
	}

	return s, nil
}

func (s *Sentry) Capture(ctx context.Context, err error) {
	if s.enabled {
		sentry.CaptureException(err)
	}
}
