package controller

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/operatorkit/controller/context/reconciliationcanceledcontext"
	"github.com/giantswarm/operatorkit/controller/context/resourcecanceledcontext"
)

func Test_ProcessDelete(t *testing.T) {
	testCases := []struct {
		Resources     []Resource
		ExpectedOrder []string
		ErrorMatcher  func(err error) bool
	}{
		// Test 0 ensures ProcessDelete returns an error in case no resources are
		// provided.
		{
			Resources:     nil,
			ExpectedOrder: nil,
			ErrorMatcher:  IsExecutionFailed,
		},

		// Test 1 ensures ProcessDelete calls EnsureDeleted method of the resource.
		{
			Resources: []Resource{
				newTestResource("r0"),
			},
			ExpectedOrder: []string{
				"r0.EnsureDeleted",
			},
			ErrorMatcher: nil,
		},

		// Test 2 ensures ProcessDelete executes all the resources in
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

		// Test 3 ensures ProcessDelete executes resources in the
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
		// Test 4 ensures ProcessDelete executes next resource after
		// resourcecanceledcontext is cancelled.
		{
			Resources: []Resource{
				newTestResource("r0"),
				newTestResource("r1"),
				newTestResource("r2").CancelResourceAt("EnsureDeleted"),
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
		err := ProcessDelete(context.Background(), nil, tc.Resources)
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

func Test_ProcessUpdate(t *testing.T) {
	testCases := []struct {
		Resources     []Resource
		ExpectedOrder []string
		ErrorMatcher  func(err error) bool
	}{
		// Test 0 ensures ProcessUpdate returns an error in case no resources are
		// provided.
		{
			Resources:     nil,
			ExpectedOrder: nil,
			ErrorMatcher:  IsExecutionFailed,
		},

		// Test 1 ensures ProcessUpdate calls EnsureCreated method of the resource.
		{
			Resources: []Resource{
				newTestResource("r0"),
			},
			ExpectedOrder: []string{
				"r0.EnsureCreated",
			},
			ErrorMatcher: nil,
		},

		// Test 2 ensures ProcessUpdate executes all resources in the
		// expected order.
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

		// Test 3 ensures ProcessUpdate executes resources in the
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
		// Test 4 ensures ProcessUpdate executes next resource after
		// resourcecanceledcontext is cancelled.
		{
			Resources: []Resource{
				newTestResource("r0"),
				newTestResource("r1"),
				newTestResource("r2").CancelResourceAt("EnsureCreated"),
				newTestResource("r3"),
				newTestResource("r4"),
			},
			ExpectedOrder: []string{
				"r0.EnsureCreated",
				"r1.EnsureCreated",
				"r2.EnsureCreated",
				"r3.EnsureCreated",
				"r4.EnsureCreated",
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		err := ProcessUpdate(context.Background(), nil, tc.Resources)
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
	resourceCanceledStep       string

	Order []string
}

func newTestResource(name string) *testResource {
	return &testResource{
		name: name,
	}
}

func (r *testResource) Name() string {
	return r.name
}

func (r *testResource) SetReconcilationCancelledAt(method string) *testResource {
	r.reconciliationCanceledStep = method
	return r
}

func (r *testResource) CancelResourceAt(method string) *testResource {
	r.resourceCanceledStep = method
	return r
}

func (r *testResource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.executeMethod(ctx, "EnsureCreated")
	return nil
}

func (r *testResource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.executeMethod(ctx, "EnsureDeleted")
	return nil
}

func (r *testResource) executeMethod(ctx context.Context, method string) {
	r.Order = append(r.Order, r.name+"."+method)

	if r.reconciliationCanceledStep == method {
		reconciliationcanceledcontext.SetCanceled(ctx)
	}
	if r.resourceCanceledStep == method {
		resourcecanceledcontext.SetCanceled(ctx)
	}
}
