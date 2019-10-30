package backoff

import (
	"github.com/cenkalti/backoff"
)

func Permanent(err error) error {
	return backoff.Permanent(err)
}
