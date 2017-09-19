package updatenecessarycontext

import (
	"context"
	"testing"
)

func Test_Framework_UpdateNecessaryContext(t *testing.T) {
	testCases := []struct {
		Ctx                       context.Context
		ExpectedIsUpdateNecessary bool
	}{
		{
			Ctx: context.TODO(),
			ExpectedIsUpdateNecessary: false,
		},
		{
			Ctx: NewContext(context.Background(), nil),
			ExpectedIsUpdateNecessary: false,
		},
		{
			Ctx: NewContext(context.Background(), make(chan struct{})),
			ExpectedIsUpdateNecessary: false,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background(), nil)
				SetUpdateNecessary(ctx)
				return ctx
			}(),
			ExpectedIsUpdateNecessary: false,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background(), make(chan struct{}))
				SetUpdateNecessary(ctx)
				return ctx
			}(),
			ExpectedIsUpdateNecessary: true,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background(), make(chan struct{}))
				SetUpdateNecessary(ctx)
				SetUpdateNecessary(ctx)
				SetUpdateNecessary(ctx)
				return ctx
			}(),
			ExpectedIsUpdateNecessary: true,
		},
	}

	for i, tc := range testCases {
		isUpdateNecessary := IsUpdateNecessary(tc.Ctx)
		if isUpdateNecessary != tc.ExpectedIsUpdateNecessary {
			t.Fatal("test", i+1, "expected", tc.ExpectedIsUpdateNecessary, "got", isUpdateNecessary)
		}
	}
}
