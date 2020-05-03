package cachekeycontext

import (
	"context"
	"testing"
	"time"
)

func Test_Controller_CacheKeyContext_PseudoUnique(t *testing.T) {
	a, _ := FromContext(NewContext(context.Background()))

	time.Sleep(1 * time.Millisecond)

	b, _ := FromContext(NewContext(context.Background()))

	if a == b {
		t.Fatalf("expected %#q and %#q to be different", a, b)
	}
}
