package k8sversion

import (
	"github.com/giantswarm/microerror"
)

var invalidKubeVersionError = &microerror.Error{
	Kind: "invalidKubeVersionError",
}

// IsInvalidKubeVersion asserts invalidKubeVersionError.
func IsInvalidKubeVersion(err error) bool {
	return microerror.Cause(err) == invalidKubeVersionError
}
