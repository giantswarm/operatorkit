package updateallowedcontext

import (
	"context"
	"testing"
)

func Test_Controller_UpdateAllowedContext(t *testing.T) {
	testCases := []struct {
		Ctx                     context.Context
		ExpectedIsUpdateAllowed bool
	}{
		{
			Ctx:                     context.TODO(),
			ExpectedIsUpdateAllowed: false,
		},
		{
			Ctx:                     NewContext(context.Background()),
			ExpectedIsUpdateAllowed: false,
		},
		{
			Ctx:                     NewContext(context.Background()),
			ExpectedIsUpdateAllowed: false,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background())
				SetUpdateAllowed(ctx)
				return ctx
			}(),
			ExpectedIsUpdateAllowed: true,
		},
		{
			Ctx: func() context.Context {
				ctx := NewContext(context.Background())
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
