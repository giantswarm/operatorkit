package framework

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/operatorkit/framework/context/reconciliationcanceledcontext"
)

func Test_processDelete(t *testing.T) {
	testCases := []struct {
		Ctx           context.Context
		Resources     []Resource
		ExpectedOrder []string
		ErrorMatcher  func(err error) bool
	}{
		// Test 0 ensures processDelete returns an error in case no resources are
		// provided.
		{
			Ctx:           context.Background(),
			Resources:     nil,
			ExpectedOrder: nil,
			ErrorMatcher:  IsExecutionFailed,
		},

		// Test 1 ensures processDelete calls EnsureDeleted method of the resource.
		{
			Ctx: context.Background(),
			Resources: []Resource{
				newTestResource("r0"),
			},
			ExpectedOrder: []string{
				"r0.EnsureDeleted",
			},
			ErrorMatcher: nil,
		},

		// Test 2 ensures processDelete executes multiple resources in
		// the expected order.
		{
			Ctx: context.Background(),
			Resources: []Resource{
				newTestResource("r0"),
				newTestResource("r1"),
			},
			ExpectedOrder: []string{
				"r0.EnsureDeleted",
				"r1.EnsureDeleted",
			},
			ErrorMatcher: nil,
		},

		// Test 3 ensures processDelete executes resources in the
		// expected order until the reconciliation gets canceled.
		{
			Ctx: reconciliationcanceledcontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				newTestResource("r0"),
				newTestResource("r1"),
				newTestResource("r2").SetReconcilationCancelledAt("EnsureDeleted"),
				newTestResource("r3"),
				newTestResource("r4"),
			},
			ExpectedOrder: []string{
				"r0.EnsureDeleted",
				"r1.EnsureDeleted",
				"r2.EnsureDeleted",
			},
			ErrorMatcher: nil,
		},

		// Test 4 ensures processDelete executes resources in the
		// expected order even if the reconciliation is canceled while
		// the given context does not contain a canceler.
		{
			Ctx: context.Background(),
			Resources: []Resource{
				newTestResource("r0"),
				newTestResource("r1"),
				newTestResource("r2").SetReconcilationCancelledAt("EnsureDeleted"),
				newTestResource("r3"),
				newTestResource("r4"),
			},
			ExpectedOrder: []string{
				"r0.EnsureDeleted",
				"r1.EnsureDeleted",
				"r2.EnsureDeleted",
				"r3.EnsureDeleted",
				"r4.EnsureDeleted",
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		err := processDelete(tc.Ctx, nil, tc.Resources)
		if err != nil {
			if tc.ErrorMatcher == nil {
				t.Fatal("test", i, "expected", nil, "got", err)
			} else if !tc.ErrorMatcher(err) {
				t.Fatal("test", i, "expected", true, "got", false)
			}
		} else {
			var order []string
			for _, r := range tc.Resources {
				order := append(order, r.(*testResource).Order...)
			}

			if !reflect.DeepEqual(tc.ExpectedOrder, rder) {
				t.Fatal("test", i, "expected", tc.ExpectedOrder, "got", order)
			}
		}
	}
}

func Test_processUpdate(t *testing.T) {
	testCases := []struct {
		Ctx           context.Context
		Resources     []Resource
		ExpectedOrder []string
		ErrorMatcher  func(err error) bool
	}{
		// Test 0 ensures processUpdate returns an error in case no resources are
		// provided.
		{
			Ctx:           context.Background(),
			Resources:     nil,
			ExpectedOrder: nil,
			ErrorMatcher:  IsExecutionFailed,
		},

		// Test 1 ensures processUpdate calls EnsureUpdated method of the resource.
		{
			Ctx: context.Background(),
			Resources: []Resource{
				newTestResource("r0"),
			},
			ExpectedOrder: []string{
				"r0.EnsureUpdated",
			},
			ErrorMatcher: nil,
		},

		// Test 2 ensures processUpdate executes multiple resources in
		// the expected order.
		{
			Ctx: context.Background(),
			Resources: []Resource{
				newTestResource("r0"),
				newTestResource("r1"),
			},
			ExpectedOrder: []string{
				"r0.EnsureUpdated",
				"r1.EnsureUpdated",
			},
			ErrorMatcher: nil,
		},

		// Test 3 ensures processUpdate executes resources in the
		// expected order until the reconciliation gets canceled.
		{
			Ctx: reconciliationcanceledcontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				newTestResource("r0"),
				newTestResource("r1"),
				newTestResource("r2").SetReconcilationCancelledAt("EnsureUpdated"),
				newTestResource("r3"),
				newTestResource("r4"),
			},
			ExpectedOrder: []string{
				"r0.EnsureUpdated",
				"r1.EnsureUpdated",
				"r2.EnsureUpdated",
			},
			ErrorMatcher: nil,
		},

		// Test 4 ensures processUpdate executes resources in the
		// expected order even if the reconciliation is canceled while
		// the given context does not contain a canceler.
		{
			Ctx: context.Background(),
			Resources: []Resource{
				newTestResource("r0"),
				newTestResource("r1"),
				newTestResource("r2").SetReconcilationCancelledAt("EnsureUpdated"),
				newTestResource("r3"),
				newTestResource("r4"),
			},
			ExpectedOrder: []string{
				"r0.EnsureUpdated",
				"r1.EnsureUpdated",
				"r2.EnsureUpdated",
				"r3.EnsureUpdated",
				"r4.EnsureUpdated",
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		err := processUpdate(tc.Ctx, nil, tc.Resources)
		if err != nil {
			if tc.ErrorMatcher == nil {
				t.Fatal("test", i, "expected", nil, "got", err)
			} else if !tc.ErrorMatcher(err) {
				t.Fatal("test", i, "expected", true, "got", false)
			}
		} else {
			var order []string
			for _, r := range tc.Resources {
				order := append(order, r.(*testResource).Order...)
			}

			if !reflect.DeepEqual(tc.ExpectedOrder, rder) {
				t.Fatal("test", i, "expected", tc.ExpectedOrder, "got", order)
			}
		}
	}
}

type testResource struct {
	Name                       string
	ReconciliationCanceledStep string
	Order                      []string

	errorCount int
}

func newTestResource(name string) *testResource {
	return &testResource{
		Name: name,
	}
}

func (r *testResource) SetReconcilationCancelledAt(method String) *testResource {
	r.ReconciliationCanceledStep = method
	return r
}

func (r *testResource) EnsureCreated(ctx context.Context, obj interface{}) error {
	m := "EnsureCreated"
	r.Order = append(r.Order, Name+"."+m)

	if r.ReconciliationCanceledStep == m {
		reconciliationcanceledcontext.SetCanceled(ctx)
	}

	return nil
}

func (r *testResource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	m := "EnsureDeleted"
	r.Order = append(r.Order, Name+"."+m)

	if r.ReconciliationCanceledStep == m {
		reconciliationcanceledcontext.SetCanceled(ctx)
	}

	return nil
}

func (r *testResource) Name() string {
	return r.Name
}
