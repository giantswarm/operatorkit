package backoff

import (
	"github.com/cenkalti/backoff"
	"github.com/giantswarm/microerror"
)

// Retry retries the operation o until it does not return error or BackOff
// stops. See https://godoc.org/github.com/cenkalti/backoff#Retry for details.
func Retry(o Operation, b BackOff) error {
	err := backoff.Retry(o.toCenkalti(), b)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

// RetryNotify does what Retry do with notification between each try.
func RetryNotify(o Operation, b BackOff, n Notify) error {
	err := backoff.RetryNotify(o.toCenkalti(), b, n.toCenkalti())
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
