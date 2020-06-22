package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
)

type Default struct {
}

func (s *Default) Capture(ctx context.Context, err error) {
	sentry.CaptureException(err)
}
