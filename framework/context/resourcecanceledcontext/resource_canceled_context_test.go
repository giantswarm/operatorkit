package resourcecanceledcontext

import (
	"context"
	"testing"
)

func Test_Framework_ResourceCanceledContext(t *testing.T) {
	testCases := []struct {
		Ctx                context.Context
		ExpectedIsCanceled bool
	}{
		{
			Ctx:                context.TODO(),
			ExpectedIsCanceled: false,
		},
		{
			Ctx:                NewContext(context.Background(), nil),
			ExpectedIsCanceled: false,
		},
		{
			Ctx:                NewContext(context.Background(), make(chan struct{})),
			ExpectedIsCanceled: false,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background(), nil)
				SetCanceled(ctx)
				return ctx
			}(),
			ExpectedIsCanceled: false,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background(), make(chan struct{}))
				SetCanceled(ctx)
				return ctx
			}(),
			ExpectedIsCanceled: true,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background(), make(chan struct{}))
				SetCanceled(ctx)
				SetCanceled(ctx)
				SetCanceled(ctx)
				return ctx
			}(),
			ExpectedIsCanceled: true,
		},
	}

	for i, tc := range testCases {
		isCanceled := IsCanceled(tc.Ctx)
		if isCanceled != tc.ExpectedIsCanceled {
			t.Fatal("test", i+1, "expected", tc.ExpectedIsCanceled, "got", isCanceled)
		}
	}
}
