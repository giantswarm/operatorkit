package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/go-errors/errors"
)

type Default struct {
}

func (s *Default) Capture(ctx context.Context, err error) {
	sentry.CaptureException(convertErrToGoError(err))
}

func convertErrToGoError(err error) *errors.Error {
	return errors.New(err)
}
