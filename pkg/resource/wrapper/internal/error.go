package internal

import "github.com/giantswarm/microerror"

var loopDetectedError = &microerror.Error{
	Kind: "loopDetectedError",
}

// IsLoopDetected asserts loopDetectedError.
func IsLoopDetected(err error) bool {
	return microerror.Cause(err) == loopDetectedError
}
