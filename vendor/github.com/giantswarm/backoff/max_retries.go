package backoff

import "github.com/cenkalti/backoff"

func WithMaxRetries(max uint64) Interface {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 0
	return backoff.WithMaxRetries(b, max)
}
