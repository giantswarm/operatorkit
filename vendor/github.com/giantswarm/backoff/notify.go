package backoff

import (
	"time"

	"github.com/cenkalti/backoff"
)

// Notify is a notify-on-error function. It receives an operation error and
// backoff delay if the operation failed (with an error).
//
// NOTE that if the backoff policy stated to stop retrying, the notify function
// isn't called.
type Notify func(error, time.Duration)

func (n Notify) toCenkalti() backoff.Notify {
	var f func(error, time.Duration)
	f = n
	return f
}
