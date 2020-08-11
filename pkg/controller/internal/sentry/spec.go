package sentry

import (
	"context"
)

type Interface interface {
	// Capture is used to send the error to Sentry API.
	Capture(ctx context.Context, err error)
}
