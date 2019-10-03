package backoff

import "time"

// BackOff describes how a backoff has to be implemented. Also see
// https://godoc.org/github.com/cenkalti/backoff#BackOff.
type BackOff interface {
	NextBackOff() time.Duration
	Reset()
}

// Interface is an alias for backward compatibility.
type Interface = BackOff
