package backoff

import (
	"context"
	"fmt"
	"time"

	"github.com/giantswarm/micrologger"
)

func NewNotifier(l micrologger.Logger, ctx context.Context) func(error, time.Duration) {
	return func(err error, d time.Duration) {
		l.LogCtx(ctx, "level", "warning", "message", fmt.Sprintf("retrying backoff in '%s' due to error", d.String()), "stack", fmt.Sprintf("%#v", err))
	}
}
