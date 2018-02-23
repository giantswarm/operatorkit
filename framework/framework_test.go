package framework

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/operatorkit/framework/context/reconciliationcanceledcontext"
)

func Test_processDelete(t *testing.T) {
	testCases := []struct {
		Resources     []Resource
		ExpectedOrder []string
		ErrorMatcher  func(err error) bool
	}{
		// Test 0 ensures processDelete returns an error in case no resources are
		// provided.
		{
			Resources:     nil,
			ExpectedOrder: nil,
			ErrorMatcher:  IsExecutionFailed,
		},

		// Test 1 ensures processDelete calls EnsureDeleted method of the resource.
		{
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
	}

	for i, tc := range testCases {
		err := processDelete(context.Background(), nil, tc.Resources)
		if err != nil {
			if tc.ErrorMatcher == nil {
				t.Fatal("test", i, "expected", nil, "got", err)
			} else if !tc.ErrorMatcher(err) {
				t.Fatal("test", i, "expected", true, "got", false)
			}
		} else {
			var order []string
			for _, r := range tc.Resources {
				order = append(order, r.(*testResource).Order...)
			}

			if !reflect.DeepEqual(tc.ExpectedOrder, order) {
				t.Fatal("test", i, "expected", tc.ExpectedOrder, "got", order)
			}
		}
	}
}

func Test_processUpdate(t *testing.T) {
	testCases := []struct {
		Resources     []Resource
		ExpectedOrder []string
		ErrorMatcher  func(err error) bool
	}{
		// Test 0 ensures processUpdate returns an error in case no resources are
		// provided.
		{
			Resources:     nil,
			ExpectedOrder: nil,
			ErrorMatcher:  IsExecutionFailed,
		},

		// Test 1 ensures processUpdate calls EnsureCreated method of the resource.
		{
			Resources: []Resource{
				newTestResource("r0"),
			},
			ExpectedOrder: []string{
				"r0.EnsureCreated",
			},
			ErrorMatcher: nil,
		},

		// Test 2 ensures processUpdate executes multiple resources in
		// the expected order.
		{
			Resources: []Resource{
				newTestResource("r0"),
				newTestResource("r1"),
			},
			ExpectedOrder: []string{
				"r0.EnsureCreated",
				"r1.EnsureCreated",
			},
			ErrorMatcher: nil,
		},

		// Test 3 ensures processUpdate executes resources in the
		// expected order until the reconciliation gets canceled.
		{
			Resources: []Resource{
				newTestResource("r0"),
				newTestResource("r1"),
				newTestResource("r2").SetReconcilationCancelledAt("EnsureCreated"),
				newTestResource("r3"),
				newTestResource("r4"),
			},
			ExpectedOrder: []string{
				"r0.EnsureCreated",
				"r1.EnsureCreated",
				"r2.EnsureCreated",
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		err := processUpdate(context.Background(), nil, tc.Resources)
		if err != nil {
			if tc.ErrorMatcher == nil {
				t.Fatal("test", i, "expected", nil, "got", err)
			} else if !tc.ErrorMatcher(err) {
				t.Fatal("test", i, "expected", true, "got", false)
			}
		} else {
			var order []string
			for _, r := range tc.Resources {
				order = append(order, r.(*testResource).Order...)
			}

			if !reflect.DeepEqual(tc.ExpectedOrder, order) {
				t.Fatal("test", i, "expected", tc.ExpectedOrder, "got", order)
			}
		}
	}
}

type testResource struct {
	name                       string
	reconciliationCanceledStep string

	Order []string
}

func newTestResource(name string) *testResource {
	return &testResource{
		name: name,
	}
}

func (r *testResource) SetReconcilationCancelledAt(method string) *testResource {
	r.reconciliationCanceledStep = method
	return r
}

func (r *testResource) EnsureCreated(ctx context.Context, obj interface{}) error {
	m := "EnsureCreated"
	r.Order = append(r.Order, r.name+"."+m)

	if r.reconciliationCanceledStep == m {
		reconciliationcanceledcontext.SetCanceled(ctx)
	}

	return nil
}

func (r *testResource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	m := "EnsureDeleted"
	r.Order = append(r.Order, r.name+"."+m)

	if r.reconciliationCanceledStep == m {
		reconciliationcanceledcontext.SetCanceled(ctx)
	}

	return nil
}

func (r *testResource) Name() string {
	return r.name
}
