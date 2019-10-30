package backoff

import (
	"time"

	"github.com/cenkalti/backoff"
)

func NewExponential(maxWait, maxInterval time.Duration) BackOff {
	b := &backoff.ExponentialBackOff{
		InitialInterval:     backoff.DefaultInitialInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         maxInterval,
		MaxElapsedTime:      maxWait,
		Clock:               backoff.SystemClock,
	}

	b.Reset()

	return b
}
