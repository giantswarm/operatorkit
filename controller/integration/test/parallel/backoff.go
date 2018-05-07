// +build k8srequired

package parallel

import (
	"time"

	"github.com/cenkalti/backoff"
)

func newConstantBackoff(maxRetries uint64) backoff.BackOff {
	return backoff.WithMaxTries(backoff.NewConstantBackOff(1*time.Second), uint64(maxRetries))
}
