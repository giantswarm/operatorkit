package resource

import "github.com/giantswarm/microerror"

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidDesiredSateError = &microerror.Error{
	Kind: "invalidDesiredSateError",
}

// IsInvalidDesiredSate asserts invalidDesiredSateError.
func IsInvalidDesiredSate(err error) bool {
	return microerror.Cause(err) == invalidDesiredSateError
}
