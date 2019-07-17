package retryresource

import (
	"testing"

	"github.com/giantswarm/operatorkit/resource/wrapper/internal"
)

func Test_resourceWrapper_Wrapper(t *testing.T) {
	// This won't compile if the *resourceWrapper doesn't implement Wrapper
	// interface.
	var _ internal.Wrapper = &resourceWrapper{}
}
