package reconciliationcanceledcontext

import (
	"context"
	"testing"
)

func Test_Controller_ReconciliationCanceledContext(t *testing.T) {
	testCases := []struct {
		Ctx                context.Context
		ExpectedIsCanceled bool
	}{
		{
			Ctx:                context.TODO(),
			ExpectedIsCanceled: false,
		},
		{
			Ctx:                NewContext(context.Background()),
			ExpectedIsCanceled: false,
		},
		{
			Ctx:                NewContext(context.Background()),
			ExpectedIsCanceled: false,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background())
				SetCanceled(ctx)
				return ctx
			}(),
			ExpectedIsCanceled: true,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background())
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
			t.Fatal("test", i, "expected", tc.ExpectedIsCanceled, "got", isCanceled)
		}
	}
}
