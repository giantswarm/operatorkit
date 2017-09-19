package updateallowedcontext

import (
	"context"
	"testing"
)

func Test_Framework_UpdateAllowedContext(t *testing.T) {
	testCases := []struct {
		Ctx                     context.Context
		ExpectedIsUpdateAllowed bool
	}{
		{
			Ctx: context.TODO(),
			ExpectedIsUpdateAllowed: false,
		},
		{
			Ctx: NewContext(context.Background(), nil),
			ExpectedIsUpdateAllowed: false,
		},
		{
			Ctx: NewContext(context.Background(), make(chan struct{})),
			ExpectedIsUpdateAllowed: false,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background(), nil)
				SetUpdateAllowed(ctx)
				return ctx
			}(),
			ExpectedIsUpdateAllowed: false,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background(), make(chan struct{}))
				SetUpdateAllowed(ctx)
				return ctx
			}(),
			ExpectedIsUpdateAllowed: true,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background(), make(chan struct{}))
				SetUpdateAllowed(ctx)
				SetUpdateAllowed(ctx)
				SetUpdateAllowed(ctx)
				return ctx
			}(),
			ExpectedIsUpdateAllowed: true,
		},
	}

	for i, tc := range testCases {
		isUpdateAllowed := IsUpdateAllowed(tc.Ctx)
		if isUpdateAllowed != tc.ExpectedIsUpdateAllowed {
			t.Fatal("test", i+1, "expected", tc.ExpectedIsUpdateAllowed, "got", isUpdateAllowed)
		}
	}
}
