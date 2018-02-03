package framework

import (
	"context"

	"github.com/cenkalti/backoff"
)

func DefaultBackOffFactory() func() backoff.BackOff {
	return func() backoff.BackOff {
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = 0
		return backoff.WithMaxTries(b, 7)
	}
}

func DefaultResourceSetResourceFunc(rs []Resource) func(ctx context.Context, obj interface{}) ([]Resource, error) {
	return func(ctx context.Context, obj interface{}) ([]Resource, error) {
		return rs, nil
	}
}
