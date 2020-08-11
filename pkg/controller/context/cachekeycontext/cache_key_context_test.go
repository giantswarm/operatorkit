package cachekeycontext

import (
	"context"
	"testing"
)

func Test_Controller_CacheKeyContext(t *testing.T) {
	s, _ := FromContext(NewContext(context.Background(), "test"))
	if s != "test" {
		t.Fatalf("expected %#q and %#q to be different", "test", s)
	}
}
