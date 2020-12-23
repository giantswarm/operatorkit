package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/giantswarm/microerror"
)

type Default struct {
}

func (s *Default) Capture(ctx context.Context, err error) {
	sentry.CaptureException(microerror.Mask(err))
}
