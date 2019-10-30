package backoff

import "github.com/cenkalti/backoff"

// An Operation is executing by Retry() or RetryNotify().
// The operation will be retried using a backoff policy if it returns an error.
type Operation func() error

func (o Operation) toCenkalti() backoff.Operation {
	var f func() error
	f = o
	return f
}
