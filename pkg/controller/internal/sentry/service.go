package sentry

import (
	"github.com/getsentry/sentry-go"
	"github.com/giantswarm/microerror"
)

type Config struct {
	DSN string
}

func New(config Config) (Interface, error) {
	if config.DSN == "" {
		return &Disabled{}, nil
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn: config.DSN,
	})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return &Default{}, nil
}
