package metricsresource

import (
	"testing"

	"github.com/giantswarm/operatorkit/resource/wrapper/internal"
)

func Test_crudResourceWrapper_Wrapper(t *testing.T) {
	// This won't compile if the *crudResourceWrapper doesn't implement Wrapper
	// interface.
	var _ internal.Wrapper = &crudResourceWrapper{}
}
