package deletionallowedcontext

import (
	"context"
	"testing"
)

func Test_Framework_DeletionAllowedContext(t *testing.T) {
	testCases := []struct {
		Ctx                       context.Context
		ExpectedIsDeletionAllowed bool
	}{
		{
			Ctx: context.TODO(),
			ExpectedIsDeletionAllowed: false,
		},
		{
			Ctx: NewContext(context.Background(), nil),
			ExpectedIsDeletionAllowed: false,
		},
		{
			Ctx: NewContext(context.Background(), make(chan struct{})),
			ExpectedIsDeletionAllowed: false,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background(), nil)
				SetDeletionAllowed(ctx)
				return ctx
			}(),
			ExpectedIsDeletionAllowed: false,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background(), make(chan struct{}))
				SetDeletionAllowed(ctx)
				return ctx
			}(),
			ExpectedIsDeletionAllowed: true,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background(), make(chan struct{}))
				SetDeletionAllowed(ctx)
				SetDeletionAllowed(ctx)
				SetDeletionAllowed(ctx)
				return ctx
			}(),
			ExpectedIsDeletionAllowed: true,
		},
	}

	for i, tc := range testCases {
		isDeletionAllowed := IsDeletionAllowed(tc.Ctx)
		if isDeletionAllowed != tc.ExpectedIsDeletionAllowed {
			t.Fatal("test", i+1, "expected", tc.ExpectedIsDeletionAllowed, "got", isDeletionAllowed)
		}
	}
}
