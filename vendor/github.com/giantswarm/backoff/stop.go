package backoff

import (
	"time"

	"github.com/cenkalti/backoff"
)

const Stop = backoff.Stop

type stopBackOff struct{}

func NewStop() BackOff {
	return &stopBackOff{}
}

func (b *stopBackOff) Reset() {}

func (b *stopBackOff) NextBackOff() time.Duration { return Stop }
