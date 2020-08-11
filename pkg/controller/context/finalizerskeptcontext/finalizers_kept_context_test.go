package finalizerskeptcontext

import (
	"context"
	"testing"
)

func Test_Controller_FinalizersKeptContext(t *testing.T) {
	testCases := []struct {
		Ctx            context.Context
		ExpectedIsKept bool
	}{
		{
			Ctx:            context.TODO(),
			ExpectedIsKept: false,
		},
		{
			Ctx:            NewContext(context.Background(), nil),
			ExpectedIsKept: false,
		},
		{
			Ctx:            NewContext(context.Background(), make(chan struct{})),
			ExpectedIsKept: false,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background(), nil)
				SetKept(ctx)
				return ctx
			}(),
			ExpectedIsKept: false,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background(), make(chan struct{}))
				SetKept(ctx)
				return ctx
			}(),
			ExpectedIsKept: true,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background(), make(chan struct{}))
				SetKept(ctx)
				SetKept(ctx)
				SetKept(ctx)
				return ctx
			}(),
			ExpectedIsKept: true,
		},
	}

	for i, tc := range testCases {
		isKept := IsKept(tc.Ctx)
		if isKept != tc.ExpectedIsKept {
			t.Fatal("test", i+1, "expected", tc.ExpectedIsKept, "got", isKept)
		}
	}
}
