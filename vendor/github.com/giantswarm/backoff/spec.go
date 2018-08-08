package backoff

import "time"

// Interface describes how a backoff has to be implemented. Also see
// https://godoc.org/github.com/cenkalti/backoff#BackOff.
type Interface interface {
	NextBackOff() time.Duration
	Reset()
}
