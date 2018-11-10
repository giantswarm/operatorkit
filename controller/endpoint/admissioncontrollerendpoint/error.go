package admissioncontrollerendpoint

import (
	"github.com/giantswarm/microerror"
)

var decodeFailedError = &microerror.Error{
	Kind: "decodeFailedError",
}

// IsDecodeFailed asserts deleteFailedError.
func IsDecodeFailed(err error) bool {
	return microerror.Cause(err) == decodeFailedError
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidRequestError = &microerror.Error{
	Kind: "invalidRequestError",
}

// IsInvalidRequest asserts invalidRequestError.
func IsInvalidRequest(err error) bool {
	return microerror.Cause(err) == invalidRequestError
}
