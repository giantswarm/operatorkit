package sentry

import (
	"context"
	"fmt"
	"reflect"

	"github.com/getsentry/sentry-go"
)

type Default struct {
}

func (s *Default) Capture(ctx context.Context, err error) {
	fmt.Printf("type: %s\n", reflect.TypeOf(err).String())
	sentry.CaptureException(err)
}
