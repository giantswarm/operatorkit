package sentry

import (
	"context"
)

type Disabled struct {
}

func (d *Disabled) Capture(ctx context.Context, err error) {
	// This implementation does nothing.
}
