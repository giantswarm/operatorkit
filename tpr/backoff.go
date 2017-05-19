package tpr

import (
	"time"

	"github.com/cenkalti/backoff"
)

type BackOff struct {
	retries  int
	interval time.Duration
}

func (b *BackOff) Reset() {}

func (b *BackOff) NextBackOff() time.Duration {
	if b.retries < 1 {
		return backoff.Stop
	}
	return b.interval
}

func NewBackOff(interval time.Duration, maxRetries int) *BackOff {
	return &BackOff{interval: interval, retries: maxRetries}
}
