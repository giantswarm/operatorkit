// package sentry provides integration for error handling with sentry.io.
package sentry

import (
	"context"

	raven "github.com/getsentry/raven-go"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

type Config struct {
	Logger micrologger.Logger

	DSN         string
	Environment string
	Release     string
}

type Service struct {
	client *raven.Client
	logger micrologger.Logger
}

func New(config Config) (*Service, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.DSN == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.DSN must not be empty", config)
	}
	if config.Environment == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Environment must not be empty", config)
	}
	if config.Release == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Release must not be empty", config)
	}

	client, err := raven.New(config.DSN)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	client.SetEnvironment(config.Environment)
	client.SetRelease(config.Release)

	s := &Service{
		client: client,
		logger: config.Logger,
	}

	return s, nil
}

func (s *Service) Capture(ctx context.Context, err error) {
	s.logger.LogCtx(ctx, "level", "info", "message", "capturing error to Sentry")
	eventID := s.client.CaptureErrorAndWait(err, map[string]string{})
	s.logger.LogCtx(ctx, "level", "info", "message", "captured error to Sentry", "eventID", eventID)
}
