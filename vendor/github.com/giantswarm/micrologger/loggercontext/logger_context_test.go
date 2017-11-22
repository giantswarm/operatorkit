package loggercontext

import (
	"context"
	"testing"
)

func Test_LoggerContext_NewContainer(t *testing.T) {
	ctx := context.Background()

	_, ok := FromContext(ctx)
	if ok {
		t.Fatalf("expected %#v got %#v", false, true)
	}

	ctx = NewContext(ctx, NewContainer())

	c1, ok := FromContext(ctx)
	if !ok {
		t.Fatalf("expected %#v got %#v", true, false)
	}
	l1 := len(c1.KeyVals)
	if l1 != 0 {
		t.Fatalf("expected %#v got %#v", 0, l1)
	}

	c1.KeyVals["foo"] = "bar"

	c2, ok := FromContext(ctx)
	if !ok {
		t.Fatalf("expected %#v got %#v", true, false)
	}
	l2 := len(c2.KeyVals)
	if l2 != 1 {
		t.Fatalf("expected %#v got %#v", 1, l2)
	}
	v2, ok := c2.KeyVals["foo"]
	if !ok {
		t.Fatalf("expected %#v got %#v", true, false)
	}
	if v2 != "bar" {
		t.Fatalf("expected %#v got %#v", "bar", v2)
	}
}
