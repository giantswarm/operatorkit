package credential

import "github.com/giantswarm/microerror"

var invalidConfig = &microerror.Error{
	Kind: "invalidConfig",
}

// IsInvalidConfigFoundError asserts invalidConfig.
func IsInvalidConfigFoundError(err error) bool {
	return microerror.Cause(err) == invalidConfig
}

var missingValueError = &microerror.Error{
	Kind: "missingValueError",
}

// IsMissingValue asserts missingValueError.
func IsMissingValue(err error) bool {
	return microerror.Cause(err) == missingValueError
}
