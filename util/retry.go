package util

import (
	"fmt"
	"time"
)

type ReadyFunc func() (bool, error)

// Retry retries f every interval until after maxRetries.
// The interval won't be affected by how long f takes.
// For example, if interval is 3s, f takes 1s, another f will be called 2s later.
// However, if f takes longer than interval, it will be delayed.
func Retry(interval time.Duration, maxRetries int, f ReadyFunc) error {
	tick := time.NewTicker(interval)
	defer tick.Stop()

	for i := 0; i < maxRetries; i++ {
		ok, err := f()
		if err != nil {
			return fmt.Errorf("util.Retry: #%d: %+v", i, err)
		}
		if ok {
			return nil
		}
		if i < maxRetries-1 {
			<-tick.C
		}
	}
	return fmt.Errorf("util.Retry: not ready after %d retry", maxRetries)
}
