package framework

//import (
//	"context"
//	"testing"
//
//	"github.com/giantswarm/operatorkit/framework/context/reconciliationcanceledcontext"
//	"github.com/giantswarm/operatorkit/framework/context/resourcecanceledcontext"
//)
//
//// Test_Framework_ResourceCallOrder ensures the resource's methods are
//// executed as expected when creating resources.
//func Test_Framework_ResourceCallOrder(t *testing.T) {
//	testCases := []struct {
//		ProcessMethod  func(ctx context.Context, obj interface{}, rs []Resource) error
//		Ctx            context.Context
//		Resources      []Resource
//		ExpectedOrders [][]string
//		ErrorMatcher   func(err error) bool
//	}{
//		// Test 0 ensures ProcessDelete returns an error in case no resources are
//		// provided.
//		{
//			ProcessMethod:  ProcessDelete,
//			Ctx:            context.TODO(),
//			Resources:      nil,
//			ExpectedOrders: nil,
//			ErrorMatcher:   IsExecutionFailed,
//		},
//
//		// Test 1 ensures ProcessDelete executes the steps of a single resource in
//		// the expected order.
//		{
//			ProcessMethod: ProcessDelete,
//			Ctx:           context.TODO(),
//			Resources: []Resource{
//				&testCRUDResourceOps{},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 2 ensures ProcessDelete executes the steps of multile resources in
//		// the expected order.
//		{
//			ProcessMethod: ProcessDelete,
//			Ctx:           context.TODO(),
//			Resources: []Resource{
//				&testCRUDResourceOps{},
//				&testCRUDResourceOps{},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 3 ensures ProcessDelete executes the steps of a single resource in
//		// the expected order until the reconciliation gets canceled.
//		{
//			ProcessMethod: ProcessDelete,
//			Ctx:           reconciliationcanceledcontext.NewContext(context.Background(), make(chan struct{})),
//			Resources: []Resource{
//				&testCRUDResourceOps{
//					ReconciliationCanceledStep: "GetCurrentState",
//				},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 4 ensures ProcessDelete executes the steps of a single resource in
//		// the expected order until the resource gets canceled.
//		{
//			ProcessMethod: ProcessDelete,
//			Ctx:           resourcecanceledcontext.NewContext(context.Background(), make(chan struct{})),
//			Resources: []Resource{
//				&testCRUDResourceOps{
//					ResourceCanceledStep: "GetCurrentState",
//				},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 5 ensures ProcessDelete executes the steps of a single resource in
//		// the expected order even if the reconciliation is canceled while the given
//		// context does not contain a canceler.
//		{
//			ProcessMethod: ProcessDelete,
//			Ctx:           context.TODO(),
//			Resources: []Resource{
//				&testCRUDResourceOps{
//					ReconciliationCanceledStep: "GetCurrentState",
//				},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 6 ensures ProcessDelete executes the steps of a single resource in
//		// the expected order even if the resource is canceled while the given
//		// context does not contain a canceler.
//		{
//			ProcessMethod: ProcessDelete,
//			Ctx:           context.TODO(),
//			Resources: []Resource{
//				&testCRUDResourceOps{
//					ResourceCanceledStep: "GetCurrentState",
//				},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 7 ensures ProcessDelete executes the steps of the first resource in
//		// the expected order in case multile resources are given, until the
//		// reconciliation of the first resource gets canceled.
//		{
//			ProcessMethod: ProcessDelete,
//			Ctx:           reconciliationcanceledcontext.NewContext(context.Background(), make(chan struct{})),
//			Resources: []Resource{
//				&testCRUDResourceOps{
//					ReconciliationCanceledStep: "NewDeletePatch",
//				},
//				&testCRUDResourceOps{},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//				},
//				nil,
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 8 ensures ProcessDelete executes the steps of the first resource in
//		// the expected order in case multile resources are given, until the
//		// resource gets canceled.
//		{
//			ProcessMethod: ProcessDelete,
//			Ctx:           resourcecanceledcontext.NewContext(context.Background(), make(chan struct{})),
//			Resources: []Resource{
//				&testCRUDResourceOps{
//					ResourceCanceledStep: "NewDeletePatch",
//				},
//				&testCRUDResourceOps{},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 9 ensures ProcessDelete executes the steps of the first and second
//		// resource in the expected order in case multile resources are given, until
//		// the reconciliation of the second resource gets canceled.
//		{
//			ProcessMethod: ProcessDelete,
//			Ctx:           reconciliationcanceledcontext.NewContext(context.Background(), make(chan struct{})),
//			Resources: []Resource{
//				&testCRUDResourceOps{},
//				&testCRUDResourceOps{
//					ReconciliationCanceledStep: "NewDeletePatch",
//				},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 10 ensures ProcessDelete executes the steps of the first and second
//		// resource in the expected order in case multile resources are given, until
//		// the second resource gets canceled.
//		{
//			ProcessMethod: ProcessDelete,
//			Ctx:           resourcecanceledcontext.NewContext(context.Background(), make(chan struct{})),
//			Resources: []Resource{
//				&testCRUDResourceOps{},
//				&testCRUDResourceOps{
//					ResourceCanceledStep: "NewDeletePatch",
//				},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 11 ensures processUpdate returns an error in case no resources are
//		// provided.
//		{
//			ProcessMethod:  processUpdate,
//			Ctx:            context.TODO(),
//			Resources:      nil,
//			ExpectedOrders: nil,
//			ErrorMatcher:   IsExecutionFailed,
//		},
//
//		// Test 12 ensures processUpdate executes the steps of a single resource in
//		// the expected order.
//		{
//			ProcessMethod: processUpdate,
//			Ctx:           context.TODO(),
//			Resources: []Resource{
//				&testCRUDResourceOps{},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 13 ensures processUpdate executes the steps of multile resources in
//		// the expected order.
//		{
//			ProcessMethod: processUpdate,
//			Ctx:           context.TODO(),
//			Resources: []Resource{
//				&testCRUDResourceOps{},
//				&testCRUDResourceOps{},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 14 ensures processUpdate executes the steps of a single resource in
//		// the expected order until the reconciliation gets canceled.
//		{
//			ProcessMethod: processUpdate,
//			Ctx:           reconciliationcanceledcontext.NewContext(context.Background(), make(chan struct{})),
//			Resources: []Resource{
//				&testCRUDResourceOps{
//					ReconciliationCanceledStep: "GetCurrentState",
//				},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 15 ensures processUpdate executes the steps of a single resource in
//		// the expected order until the resource gets canceled.
//		{
//			ProcessMethod: processUpdate,
//			Ctx:           resourcecanceledcontext.NewContext(context.Background(), make(chan struct{})),
//			Resources: []Resource{
//				&testCRUDResourceOps{
//					ResourceCanceledStep: "GetCurrentState",
//				},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 16 ensures processUpdate executes the steps of a single resource in
//		// the expected order even if the reconciliation is canceled while the given
//		// context does not contain a canceler.
//		{
//			ProcessMethod: processUpdate,
//			Ctx:           context.TODO(),
//			Resources: []Resource{
//				&testCRUDResourceOps{
//					ReconciliationCanceledStep: "GetCurrentState",
//				},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 17 ensures processUpdate executes the steps of a single resource in
//		// the expected order even if the resource is canceled while the given
//		// context does not contain a canceler.
//		{
//			ProcessMethod: processUpdate,
//			Ctx:           context.TODO(),
//			Resources: []Resource{
//				&testCRUDResourceOps{
//					ResourceCanceledStep: "GetCurrentState",
//				},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 18 ensures processUpdate executes the steps of the first resource in
//		// the expected order in case multile resources are given, until the
//		// reconciliation of the first resource gets canceled.
//		{
//			ProcessMethod: processUpdate,
//			Ctx:           reconciliationcanceledcontext.NewContext(context.Background(), make(chan struct{})),
//			Resources: []Resource{
//				&testCRUDResourceOps{
//					ReconciliationCanceledStep: "NewUpdatePatch",
//				},
//				&testCRUDResourceOps{},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//				},
//				nil,
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 19 ensures processUpdate executes the steps of the first resource in
//		// the expected order in case multile resources are given, until the first
//		// resource gets canceled.
//		{
//			ProcessMethod: processUpdate,
//			Ctx:           resourcecanceledcontext.NewContext(context.Background(), make(chan struct{})),
//			Resources: []Resource{
//				&testCRUDResourceOps{
//					ResourceCanceledStep: "NewUpdatePatch",
//				},
//				&testCRUDResourceOps{},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 20 ensures processUpdate executes the steps of the first and second
//		// resource in the expected order in case multile resources are given, until
//		// the reconciliation of the second resource gets canceled.
//		{
//			ProcessMethod: processUpdate,
//			Ctx:           reconciliationcanceledcontext.NewContext(context.Background(), make(chan struct{})),
//			Resources: []Resource{
//				&testCRUDResourceOps{},
//				&testCRUDResourceOps{
//					ReconciliationCanceledStep: "NewUpdatePatch",
//				},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 21 ensures processUpdate executes the steps of the first and second
//		// resource in the expected order in case multile resources are given, until
//		// the second resource gets canceled.
//		{
//			ProcessMethod: processUpdate,
//			Ctx:           resourcecanceledcontext.NewContext(context.Background(), make(chan struct{})),
//			Resources: []Resource{
//				&testCRUDResourceOps{},
//				&testCRUDResourceOps{
//					ResourceCanceledStep: "NewUpdatePatch",
//				},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//					"ApplyUpdatePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 22 ensures ProcessDelete calls Resource.Apply*Patch
//		// only when Patch has corresponding part set.
//		{
//			ProcessMethod: ProcessDelete,
//			Ctx:           context.TODO(),
//			Resources: []Resource{
//				&testCRUDResourceOps{
//					SetupPatchFunc: func(p *Patch) {
//					},
//				},
//				&testCRUDResourceOps{
//					SetupPatchFunc: func(p *Patch) {
//						p.SetCreateChange("test create data")
//					},
//				},
//				&testCRUDResourceOps{
//					SetupPatchFunc: func(p *Patch) {
//						p.SetDeleteChange("test delete data")
//					},
//				},
//				&testCRUDResourceOps{
//					SetupPatchFunc: func(p *Patch) {
//						p.SetUpdateChange("test update data")
//					},
//				},
//				&testCRUDResourceOps{
//					SetupPatchFunc: func(p *Patch) {
//						p.SetCreateChange("test create data")
//						p.SetDeleteChange("test delete data")
//					},
//				},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//					"ApplyCreatePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//					"ApplyDeletePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//					"ApplyUpdatePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewDeletePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//
//		// Test 23 ensures processUpdate calls Resource.Apply*Patch
//		// only when Patch has corresponding part set.
//		{
//			ProcessMethod: processUpdate,
//			Ctx:           context.TODO(),
//			Resources: []Resource{
//				&testCRUDResourceOps{
//					SetupPatchFunc: func(p *Patch) {
//					},
//				},
//				&testCRUDResourceOps{
//					SetupPatchFunc: func(p *Patch) {
//						p.SetCreateChange("test create data")
//					},
//				},
//				&testCRUDResourceOps{
//					SetupPatchFunc: func(p *Patch) {
//						p.SetDeleteChange("test delete data")
//					},
//				},
//				&testCRUDResourceOps{
//					SetupPatchFunc: func(p *Patch) {
//						p.SetUpdateChange("test update data")
//					},
//				},
//				&testCRUDResourceOps{
//					SetupPatchFunc: func(p *Patch) {
//						p.SetCreateChange("test create data")
//						p.SetDeleteChange("test delete data")
//					},
//				},
//			},
//			ExpectedOrders: [][]string{
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//					"ApplyCreatePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//					"ApplyDeletePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//					"ApplyUpdatePatch",
//				},
//				{
//					"GetCurrentState",
//					"GetDesiredState",
//					"NewUpdatePatch",
//					"ApplyCreatePatch",
//					"ApplyDeletePatch",
//				},
//			},
//			ErrorMatcher: nil,
//		},
//	}
//
//	for i, tc := range testCases {
//		err := tc.ProcessMethod(tc.Ctx, nil, tc.Resources)
//		if err != nil {
//			if tc.ErrorMatcher == nil {
//				t.Fatal("test", i, "expected", nil, "got", err)
//			} else if !tc.ErrorMatcher(err) {
//				t.Fatal("test", i, "expected", true, "got", false)
//			}
//		} else {
//			if len(tc.Resources) != len(tc.ExpectedOrders) {
//				t.Fatal("test", i, "expected", len(tc.ExpectedOrders), "got", len(tc.ExpectedOrders))
//			}
//
//			for j, r := range tc.Resources {
//				if !reflect.DeepEqual(tc.ExpectedOrders[j], r.(*testCRUDResourceOps).Order) {
//					t.Fatal("test", i, "expected", tc.ExpectedOrders[j], "got", r.(*testCRUDResourceOps).Order)
//				}
//			}
//		}
//	}
//}
//
//type testCRUDResourceOps struct {
//	ReconciliationCanceledStep string
//	ResourceCanceledStep       string
//	Error                      error
//	ErrorCount                 int
//	ErrorMethod                string
//	Order                      []string
//	SetupPatchFunc             func(p *Patch)
//
//	errorCount int
//}
//
//func (r *testCRUDResourceOps) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
//	m := "GetCurrentState"
//	r.Order = append(r.Order, m)
//
//	if r.ReconciliationCanceledStep == m {
//		reconciliationcanceledcontext.SetCanceled(ctx)
//		if reconciliationcanceledcontext.IsCanceled(ctx) {
//			return nil, nil
//		}
//	}
//	if r.ResourceCanceledStep == m {
//		resourcecanceledcontext.SetCanceled(ctx)
//		if resourcecanceledcontext.IsCanceled(ctx) {
//			return nil, nil
//		}
//	}
//
//	if r.returnErrorFor(m) {
//		return nil, r.Error
//	}
//
//	return nil, nil
//}
//
//func (r *testCRUDResourceOps) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
//	m := "GetDesiredState"
//	r.Order = append(r.Order, m)
//
//	if r.ReconciliationCanceledStep == m {
//		reconciliationcanceledcontext.SetCanceled(ctx)
//		if reconciliationcanceledcontext.IsCanceled(ctx) {
//			return nil, nil
//		}
//	}
//	if r.ResourceCanceledStep == m {
//		resourcecanceledcontext.SetCanceled(ctx)
//		if resourcecanceledcontext.IsCanceled(ctx) {
//			return nil, nil
//		}
//	}
//
//	if r.returnErrorFor(m) {
//		return nil, r.Error
//	}
//
//	return nil, nil
//}
//
//func (r *testCRUDResourceOps) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error) {
//	m := "NewUpdatePatch"
//	r.Order = append(r.Order, m)
//
//	if r.ReconciliationCanceledStep == m {
//		reconciliationcanceledcontext.SetCanceled(ctx)
//		if reconciliationcanceledcontext.IsCanceled(ctx) {
//			return NewPatch(), nil
//		}
//	}
//	if r.ResourceCanceledStep == m {
//		resourcecanceledcontext.SetCanceled(ctx)
//		if resourcecanceledcontext.IsCanceled(ctx) {
//			return NewPatch(), nil
//		}
//	}
//
//	if r.returnErrorFor(m) {
//		return nil, r.Error
//	}
//
//	p := NewPatch()
//	if r.SetupPatchFunc != nil {
//		r.SetupPatchFunc(p)
//	} else {
//		p.SetCreateChange("test create data")
//		p.SetUpdateChange("test update data")
//		p.SetDeleteChange("test delete data")
//	}
//	return p, nil
//}
//
//func (r *testCRUDResourceOps) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error) {
//	m := "NewDeletePatch"
//	r.Order = append(r.Order, m)
//
//	if r.ReconciliationCanceledStep == m {
//		reconciliationcanceledcontext.SetCanceled(ctx)
//		if reconciliationcanceledcontext.IsCanceled(ctx) {
//			return NewPatch(), nil
//		}
//	}
//	if r.ResourceCanceledStep == m {
//		resourcecanceledcontext.SetCanceled(ctx)
//		if resourcecanceledcontext.IsCanceled(ctx) {
//			return NewPatch(), nil
//		}
//	}
//
//	if r.returnErrorFor(m) {
//		return nil, r.Error
//	}
//
//	p := NewPatch()
//	if r.SetupPatchFunc != nil {
//		r.SetupPatchFunc(p)
//	} else {
//		p.SetCreateChange("test create data")
//		p.SetUpdateChange("test update data")
//		p.SetDeleteChange("test delete data")
//	}
//	return p, nil
//}
//
//func (r *testCRUDResourceOps) Name() string {
//	return "testCRUDResourceOps"
//}
//
//func (r *testCRUDResourceOps) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
//	m := "ApplyCreatePatch"
//	r.Order = append(r.Order, m)
//
//	if r.ReconciliationCanceledStep == m {
//		reconciliationcanceledcontext.SetCanceled(ctx)
//		if reconciliationcanceledcontext.IsCanceled(ctx) {
//			return nil
//		}
//	}
//	if r.ResourceCanceledStep == m {
//		resourcecanceledcontext.SetCanceled(ctx)
//		if resourcecanceledcontext.IsCanceled(ctx) {
//			return nil
//		}
//	}
//
//	if r.returnErrorFor(m) {
//		return r.Error
//	}
//
//	return nil
//}
//
//func (r *testCRUDResourceOps) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
//	m := "ApplyDeletePatch"
//	r.Order = append(r.Order, m)
//
//	if r.ReconciliationCanceledStep == m {
//		reconciliationcanceledcontext.SetCanceled(ctx)
//		if reconciliationcanceledcontext.IsCanceled(ctx) {
//			return nil
//		}
//	}
//	if r.ResourceCanceledStep == m {
//		resourcecanceledcontext.SetCanceled(ctx)
//		if resourcecanceledcontext.IsCanceled(ctx) {
//			return nil
//		}
//	}
//
//	if r.returnErrorFor(m) {
//		return r.Error
//	}
//
//	return nil
//}
//
//func (r *testCRUDResourceOps) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
//	m := "ApplyUpdatePatch"
//	r.Order = append(r.Order, m)
//
//	if r.ReconciliationCanceledStep == m {
//		reconciliationcanceledcontext.SetCanceled(ctx)
//		if reconciliationcanceledcontext.IsCanceled(ctx) {
//			return nil
//		}
//	}
//	if r.ResourceCanceledStep == m {
//		resourcecanceledcontext.SetCanceled(ctx)
//		if resourcecanceledcontext.IsCanceled(ctx) {
//			return nil
//		}
//	}
//
//	if r.returnErrorFor(m) {
//		return r.Error
//	}
//
//	return nil
//}
//
//func (r *testCRUDResourceOps) Underlying() Resource {
//	return r
//}
//
//func (r *testCRUDResourceOps) returnErrorFor(errorMethod string) bool {
//	ok := r.Error != nil && r.ErrorCount > r.errorCount && r.ErrorMethod == errorMethod
//
//	if ok {
//		r.errorCount++
//		return true
//	}
//
//	return false
//}
