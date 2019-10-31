package controller

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger/microloggertest"

	"github.com/giantswarm/operatorkit/controller/context/reconciliationcanceledcontext"
	"github.com/giantswarm/operatorkit/controller/context/resourcecanceledcontext"
	"github.com/giantswarm/operatorkit/resource"
	"github.com/giantswarm/operatorkit/resource/crud"
)

func Test_ProcessDelete_CRUD(t *testing.T) {
	testCases := []struct {
		Resources     []resource.Interface
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
			Resources: []resource.Interface{
				newTestCRUDResource("r0"),
			},
			ExpectedOrder: []string{
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewDeletePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",
			},
			ErrorMatcher: nil,
		},

		// Test 2 ensures ProcessDelete executes all the resources in
		// the expected order.
		{
			Resources: []resource.Interface{
				newTestCRUDResource("r0"),
				newTestCRUDResource("r1"),
			},
			ExpectedOrder: []string{
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewDeletePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",

				"r1.GetCurrentState",
				"r1.GetDesiredState",
				"r1.NewDeletePatch",
				"r1.ApplyCreateChange",
				"r1.ApplyDeleteChange",
				"r1.ApplyUpdateChange",
			},
			ErrorMatcher: nil,
		},

		// Test 3 ensures ProcessDelete executes resources in the
		// expected order until the reconciliation gets canceled.
		{
			Resources: []resource.Interface{
				newTestCRUDResource("r0"),
				newTestCRUDResource("r1"),
				newTestCRUDResource("r2").CancelReconciliationAt("GetDesiredState"),
				newTestCRUDResource("r3"),
				newTestCRUDResource("r4"),
			},
			ExpectedOrder: []string{
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewDeletePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",

				"r1.GetCurrentState",
				"r1.GetDesiredState",
				"r1.NewDeletePatch",
				"r1.ApplyCreateChange",
				"r1.ApplyDeleteChange",
				"r1.ApplyUpdateChange",

				"r2.GetCurrentState",
				"r2.GetDesiredState",
			},
			ErrorMatcher: nil,
		},
		// Test 4 ensures ProcessDelete executes next resource after
		// resourcecanceledcontext is cancelled.
		{
			Resources: []resource.Interface{
				newTestCRUDResource("r0").CancelResourceAt("ApplyDeleteChange"),
				newTestCRUDResource("r1"),
				newTestCRUDResource("r2").CancelResourceAt("GetDesiredState"),
				newTestCRUDResource("r3"),
				newTestCRUDResource("r4"),
			},
			ExpectedOrder: []string{
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewDeletePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",

				"r1.GetCurrentState",
				"r1.GetDesiredState",
				"r1.NewDeletePatch",
				"r1.ApplyCreateChange",
				"r1.ApplyDeleteChange",
				"r1.ApplyUpdateChange",

				"r2.GetCurrentState",
				"r2.GetDesiredState",

				"r3.GetCurrentState",
				"r3.GetDesiredState",
				"r3.NewDeletePatch",
				"r3.ApplyCreateChange",
				"r3.ApplyDeleteChange",
				"r3.ApplyUpdateChange",

				"r4.GetCurrentState",
				"r4.GetDesiredState",
				"r4.NewDeletePatch",
				"r4.ApplyCreateChange",
				"r4.ApplyDeleteChange",
				"r4.ApplyUpdateChange",
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
				order = append(order, r.(*testCRUDResource).Order()...)
			}

			if !reflect.DeepEqual(tc.ExpectedOrder, order) {
				t.Fatal("test", i, "expected", tc.ExpectedOrder, "got", order)
			}
		}
	}
}

func Test_ProcessUpdate_CRUD(t *testing.T) {
	testCases := []struct {
		Resources     []resource.Interface
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

		// Test 1 ensures ProcessUpdate calls EnsureDeleted method of the resource.
		{
			Resources: []resource.Interface{
				newTestCRUDResource("r0"),
			},
			ExpectedOrder: []string{
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewUpdatePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",
			},
			ErrorMatcher: nil,
		},

		// Test 2 ensures ProcessUpdate executes all the resources in
		// the expected order.
		{
			Resources: []resource.Interface{
				newTestCRUDResource("r0"),
				newTestCRUDResource("r1"),
			},
			ExpectedOrder: []string{
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewUpdatePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",

				"r1.GetCurrentState",
				"r1.GetDesiredState",
				"r1.NewUpdatePatch",
				"r1.ApplyCreateChange",
				"r1.ApplyDeleteChange",
				"r1.ApplyUpdateChange",
			},
			ErrorMatcher: nil,
		},

		// Test 3 ensures ProcessUpdate executes resources in the
		// expected order until the reconciliation gets canceled.
		{
			Resources: []resource.Interface{
				newTestCRUDResource("r0"),
				newTestCRUDResource("r1"),
				newTestCRUDResource("r2").CancelReconciliationAt("GetDesiredState"),
				newTestCRUDResource("r3"),
				newTestCRUDResource("r4"),
			},
			ExpectedOrder: []string{
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewUpdatePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",

				"r1.GetCurrentState",
				"r1.GetDesiredState",
				"r1.NewUpdatePatch",
				"r1.ApplyCreateChange",
				"r1.ApplyDeleteChange",
				"r1.ApplyUpdateChange",

				"r2.GetCurrentState",
				"r2.GetDesiredState",
			},
			ErrorMatcher: nil,
		},
		// Test 4 ensures ProcessUpdate executes next resource after
		// resourcecanceledcontext is cancelled.
		{
			Resources: []resource.Interface{
				newTestCRUDResource("r0").CancelResourceAt("ApplyDeleteChange"),
				newTestCRUDResource("r1"),
				newTestCRUDResource("r2").CancelResourceAt("GetDesiredState"),
				newTestCRUDResource("r3"),
				newTestCRUDResource("r4"),
			},
			ExpectedOrder: []string{
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewUpdatePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",

				"r1.GetCurrentState",
				"r1.GetDesiredState",
				"r1.NewUpdatePatch",
				"r1.ApplyCreateChange",
				"r1.ApplyDeleteChange",
				"r1.ApplyUpdateChange",

				"r2.GetCurrentState",
				"r2.GetDesiredState",

				"r3.GetCurrentState",
				"r3.GetDesiredState",
				"r3.NewUpdatePatch",
				"r3.ApplyCreateChange",
				"r3.ApplyDeleteChange",
				"r3.ApplyUpdateChange",

				"r4.GetCurrentState",
				"r4.GetDesiredState",
				"r4.NewUpdatePatch",
				"r4.ApplyCreateChange",
				"r4.ApplyDeleteChange",
				"r4.ApplyUpdateChange",
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
				order = append(order, r.(*testCRUDResource).Order()...)
			}

			if !reflect.DeepEqual(tc.ExpectedOrder, order) {
				t.Fatal("test", i, "expected", tc.ExpectedOrder, "got", order)
			}
		}
	}
}

type testCRUDResource struct {
	*crud.Resource
	ops *testCRUDResourceOps
}

func newTestCRUDResource(name string) *testCRUDResource {
	var err error

	ops := newTestResourceOps(name)

	var crudResource *crud.Resource
	{
		c := crud.ResourceConfig{
			CRUD:   ops,
			Logger: microloggertest.New(),
		}

		crudResource, err = crud.NewResource(c)
		if err != nil {
			panic(fmt.Sprintf("%#v", microerror.Mask(err)))
		}
	}

	return &testCRUDResource{
		Resource: crudResource,
		ops:      ops,
	}
}

func (r *testCRUDResource) CancelReconciliationAt(method string) *testCRUDResource {
	r.ops.CancelReconciliationAt(method)
	return r
}

func (r *testCRUDResource) CancelResourceAt(method string) *testCRUDResource {
	r.ops.CancelResourceAt(method)
	return r
}

func (r *testCRUDResource) Order() []string {
	return append([]string{}, r.ops.Order...)
}

type testCRUDResourceOps struct {
	name                       string
	reconciliationCanceledStep string
	resourceCanceledStep       string

	Order []string
}

func newTestResourceOps(name string) *testCRUDResourceOps {
	return &testCRUDResourceOps{
		name: name,

		Order: []string{},
	}
}

func (o *testCRUDResourceOps) Name() string {
	return o.name
}

func (o *testCRUDResourceOps) CancelReconciliationAt(method string) {
	o.reconciliationCanceledStep = method
}

func (o *testCRUDResourceOps) CancelResourceAt(method string) {
	o.resourceCanceledStep = method
}

func (o *testCRUDResourceOps) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	o.executeMethod(ctx, "GetCurrentState")
	return nil, nil
}

func (o *testCRUDResourceOps) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	o.executeMethod(ctx, "GetDesiredState")
	return nil, nil
}

func (o *testCRUDResourceOps) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*crud.Patch, error) {
	o.executeMethod(ctx, "NewUpdatePatch")
	return newFullPatch(), nil
}

func (o *testCRUDResourceOps) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*crud.Patch, error) {
	o.executeMethod(ctx, "NewDeletePatch")
	return newFullPatch(), nil
}

func (o *testCRUDResourceOps) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	o.executeMethod(ctx, "ApplyCreateChange")
	return nil
}

func (o *testCRUDResourceOps) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	o.executeMethod(ctx, "ApplyDeleteChange")
	return nil
}

func (o *testCRUDResourceOps) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	o.executeMethod(ctx, "ApplyUpdateChange")
	return nil
}

func (o *testCRUDResourceOps) executeMethod(ctx context.Context, method string) {
	o.Order = append(o.Order, o.name+"."+method)

	if o.reconciliationCanceledStep == method {
		reconciliationcanceledcontext.SetCanceled(ctx)
	}
	if o.resourceCanceledStep == method {
		resourcecanceledcontext.SetCanceled(ctx)
	}
}

// newFullPatch returns Patch filled with nil's so all Apply*Change methods are
// executed.
func newFullPatch() *crud.Patch {
	p := crud.NewPatch()
	p.SetCreateChange(nil)
	p.SetDeleteChange(nil)
	p.SetUpdateChange(nil)

	return p
}
