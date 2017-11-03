package framework

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/giantswarm/operatorkit/framework/context/canceledcontext"
)

// Test_Framework_ResourceCallOrder ensures the resource's methods are
// executed as expected when creating resources.
func Test_Framework_ResourceCallOrder(t *testing.T) {
	testCases := []struct {
		ProcessMethod  func(ctx context.Context, obj interface{}, rs []Resource) error
		Ctx            context.Context
		Resources      []Resource
		ExpectedOrders [][]string
		ErrorMatcher   func(err error) bool
	}{
		// Test 1 ensures ProcessCreate returns an error in case no resources are
		// provided.
		{
			ProcessMethod:  ProcessCreate,
			Ctx:            context.TODO(),
			Resources:      nil,
			ExpectedOrders: nil,
			ErrorMatcher:   IsExecutionFailed,
		},

		// Test 2 ensures ProcessCreate executes the steps of a single resource in
		// the expected order.
		{
			ProcessMethod: ProcessCreate,
			Ctx:           context.TODO(),
			Resources: []Resource{
				&testResource{},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
					"ApplyUpdatePatch",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 3 ensures ProcessCreate executes the steps of multile resources in
		// the expected order.
		{
			ProcessMethod: ProcessCreate,
			Ctx:           context.TODO(),
			Resources: []Resource{
				&testResource{},
				&testResource{},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
					"ApplyUpdatePatch",
				},
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
					"ApplyUpdatePatch",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 4 ensures ProcessCreate executes the steps of a single resource in
		// the expected order until it gets canceled.
		{
			ProcessMethod: ProcessCreate,
			Ctx:           canceledcontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{
					CancelingStep: "GetCurrentState",
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=false)",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 5 ensures ProcessCreate executes the steps of a single resource in
		// the expected order even if the resource is canceled while the given
		// context does not contain a canceler.
		{
			ProcessMethod: ProcessCreate,
			Ctx:           context.TODO(),
			Resources: []Resource{
				&testResource{
					CancelingStep: "GetCurrentState",
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
					"ApplyUpdatePatch",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 6 ensures ProcessCreate executes the steps of the first resource in
		// the expected order in case multile resources are given, until the first
		// resource gets canceled.
		{
			ProcessMethod: ProcessCreate,
			Ctx:           canceledcontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{
					CancelingStep: "NewPatch",
				},
				&testResource{},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
				},
				nil,
			},
			ErrorMatcher: nil,
		},

		// Test 7 ensures ProcessCreate executes the steps of the first and second resource in
		// the expected order in case multile resources are given, until the second
		// resource gets canceled.
		{
			ProcessMethod: ProcessCreate,
			Ctx:           canceledcontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{},
				&testResource{
					CancelingStep: "NewPatch",
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
					"ApplyUpdatePatch",
				},
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 8 ensures ProcessDelete returns an error in case no resources are
		// provided.
		{
			ProcessMethod:  ProcessDelete,
			Ctx:            context.TODO(),
			Resources:      nil,
			ExpectedOrders: nil,
			ErrorMatcher:   IsExecutionFailed,
		},

		// Test 9 ensures ProcessDelete executes the steps of a single resource in
		// the expected order.
		{
			ProcessMethod: ProcessDelete,
			Ctx:           context.TODO(),
			Resources: []Resource{
				&testResource{},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=true)",
					"GetDesiredState(deleted=true)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
					"ApplyUpdatePatch",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 10 ensures ProcessDelete executes the steps of multile resources in
		// the expected order.
		{
			ProcessMethod: ProcessDelete,
			Ctx:           context.TODO(),
			Resources: []Resource{
				&testResource{},
				&testResource{},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=true)",
					"GetDesiredState(deleted=true)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
					"ApplyUpdatePatch",
				},
				{
					"GetCurrentState(deleted=true)",
					"GetDesiredState(deleted=true)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
					"ApplyUpdatePatch",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 11 ensures ProcessDelete executes the steps of a single resource in
		// the expected order until it gets canceled.
		{
			ProcessMethod: ProcessDelete,
			Ctx:           canceledcontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{
					CancelingStep: "GetCurrentState",
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=true)",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 12 ensures ProcessDelete executes the steps of a single resource in
		// the expected order even if the resource is canceled while the given
		// context does not contain a canceler.
		{
			ProcessMethod: ProcessDelete,
			Ctx:           context.TODO(),
			Resources: []Resource{
				&testResource{
					CancelingStep: "GetCurrentState",
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=true)",
					"GetDesiredState(deleted=true)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
					"ApplyUpdatePatch",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 13 ensures ProcessDelete executes the steps of the first resource in
		// the expected order in case multile resources are given, until the first
		// resource gets canceled.
		{
			ProcessMethod: ProcessDelete,
			Ctx:           canceledcontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{
					CancelingStep: "NewPatch",
				},
				&testResource{},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=true)",
					"GetDesiredState(deleted=true)",
					"NewPatch",
				},
				nil,
			},
			ErrorMatcher: nil,
		},

		// Test 14 ensures ProcessDelete executes the steps of the first and second resource in
		// the expected order in case multile resources are given, until the second
		// resource gets canceled.
		{
			ProcessMethod: ProcessDelete,
			Ctx:           canceledcontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{},
				&testResource{
					CancelingStep: "NewPatch",
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=true)",
					"GetDesiredState(deleted=true)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
					"ApplyUpdatePatch",
				},
				{
					"GetCurrentState(deleted=true)",
					"GetDesiredState(deleted=true)",
					"NewPatch",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 15 ensures ProcessUpdate returns an error in case no resources are
		// provided.
		{
			ProcessMethod:  ProcessUpdate,
			Ctx:            context.TODO(),
			Resources:      nil,
			ExpectedOrders: nil,
			ErrorMatcher:   IsExecutionFailed,
		},

		// Test 16 ensures ProcessUpdate executes the steps of a single resource in
		// the expected order.
		{
			ProcessMethod: ProcessUpdate,
			Ctx:           context.TODO(),
			Resources: []Resource{
				&testResource{},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
					"ApplyUpdatePatch",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 17 ensures ProcessUpdate executes the steps of multile resources in
		// the expected order.
		{
			ProcessMethod: ProcessUpdate,
			Ctx:           context.TODO(),
			Resources: []Resource{
				&testResource{},
				&testResource{},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
					"ApplyUpdatePatch",
				},
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
					"ApplyUpdatePatch",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 18 ensures ProcessUpdate executes the steps of a single resource in
		// the expected order until it gets canceled.
		{
			ProcessMethod: ProcessUpdate,
			Ctx:           canceledcontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{
					CancelingStep: "GetCurrentState",
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=false)",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 19 ensures ProcessUpdate executes the steps of a single resource in
		// the expected order even if the resource is canceled while the given
		// context does not contain a canceler.
		{
			ProcessMethod: ProcessUpdate,
			Ctx:           context.TODO(),
			Resources: []Resource{
				&testResource{
					CancelingStep: "GetCurrentState",
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
					"ApplyUpdatePatch",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 20 ensures ProcessUpdate executes the steps of the first resource in
		// the expected order in case multile resources are given, until the first
		// resource gets canceled.
		{
			ProcessMethod: ProcessUpdate,
			Ctx:           canceledcontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{
					CancelingStep: "NewPatch",
				},
				&testResource{},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
				},
				nil,
			},
			ErrorMatcher: nil,
		},

		// Test 21 ensures ProcessUpdate executes the steps of the first and second resource in
		// the expected order in case multile resources are given, until the second
		// resource gets canceled.
		{
			ProcessMethod: ProcessUpdate,
			Ctx:           canceledcontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{},
				&testResource{
					CancelingStep: "NewPatch",
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
					"ApplyUpdatePatch",
				},
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 22 ensures ProcessCreate calls Resource.Apply*Patch
		// only when Patch has corresponding part set.
		{
			ProcessMethod: ProcessCreate,
			Ctx:           canceledcontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{
					SetupPatchFunc: func(p *Patch) {
					},
				},
				&testResource{
					SetupPatchFunc: func(p *Patch) {
						p.SetCreateChange("test create data")
					},
				},
				&testResource{
					SetupPatchFunc: func(p *Patch) {
						p.SetDeleteChange("test delete data")
					},
				},
				&testResource{
					SetupPatchFunc: func(p *Patch) {
						p.SetUpdateChange("test update data")
					},
				},
				&testResource{
					SetupPatchFunc: func(p *Patch) {
						p.SetCreateChange("test create data")
						p.SetDeleteChange("test delete data")
					},
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
				},
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyCreatePatch",
				},
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyDeletePatch",
				},
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyUpdatePatch",
				},
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 23 ensures ProcessDelete calls Resource.Apply*Patch
		// only when Patch has corresponding part set.
		{
			ProcessMethod: ProcessDelete,
			Ctx:           canceledcontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{
					SetupPatchFunc: func(p *Patch) {
					},
				},
				&testResource{
					SetupPatchFunc: func(p *Patch) {
						p.SetCreateChange("test create data")
					},
				},
				&testResource{
					SetupPatchFunc: func(p *Patch) {
						p.SetDeleteChange("test delete data")
					},
				},
				&testResource{
					SetupPatchFunc: func(p *Patch) {
						p.SetUpdateChange("test update data")
					},
				},
				&testResource{
					SetupPatchFunc: func(p *Patch) {
						p.SetCreateChange("test create data")
						p.SetDeleteChange("test delete data")
					},
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=true)",
					"GetDesiredState(deleted=true)",
					"NewPatch",
				},
				{
					"GetCurrentState(deleted=true)",
					"GetDesiredState(deleted=true)",
					"NewPatch",
					"ApplyCreatePatch",
				},
				{
					"GetCurrentState(deleted=true)",
					"GetDesiredState(deleted=true)",
					"NewPatch",
					"ApplyDeletePatch",
				},
				{
					"GetCurrentState(deleted=true)",
					"GetDesiredState(deleted=true)",
					"NewPatch",
					"ApplyUpdatePatch",
				},
				{
					"GetCurrentState(deleted=true)",
					"GetDesiredState(deleted=true)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 24 ensures ProcessUpdate calls Resource.Apply*Patch
		// only when Patch has corresponding part set.
		{
			ProcessMethod: ProcessUpdate,
			Ctx:           canceledcontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{
					SetupPatchFunc: func(p *Patch) {
					},
				},
				&testResource{
					SetupPatchFunc: func(p *Patch) {
						p.SetCreateChange("test create data")
					},
				},
				&testResource{
					SetupPatchFunc: func(p *Patch) {
						p.SetDeleteChange("test delete data")
					},
				},
				&testResource{
					SetupPatchFunc: func(p *Patch) {
						p.SetUpdateChange("test update data")
					},
				},
				&testResource{
					SetupPatchFunc: func(p *Patch) {
						p.SetCreateChange("test create data")
						p.SetDeleteChange("test delete data")
					},
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
				},
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyCreatePatch",
				},
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyDeletePatch",
				},
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyUpdatePatch",
				},
				{
					"GetCurrentState(deleted=false)",
					"GetDesiredState(deleted=false)",
					"NewPatch",
					"ApplyCreatePatch",
					"ApplyDeletePatch",
				},
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		err := tc.ProcessMethod(tc.Ctx, nil, tc.Resources)
		if err != nil {
			if tc.ErrorMatcher == nil {
				t.Fatal("test", i+1, "expected", "error matcher", "got", nil)
			} else if !tc.ErrorMatcher(err) {
				t.Fatal("test", i+1, "expected", true, "got", false)
			}
		} else {
			if len(tc.Resources) != len(tc.ExpectedOrders) {
				t.Fatal("test", i+1, "expected", len(tc.ExpectedOrders), "got", len(tc.ExpectedOrders))
			}

			for j, r := range tc.Resources {
				if !reflect.DeepEqual(tc.ExpectedOrders[j], r.(*testResource).Order) {
					t.Fatal("test", i+1, "expected", tc.ExpectedOrders[j], "got", r.(*testResource).Order)
				}
			}
		}
	}
}

type testResource struct {
	CancelingStep  string
	Error          error
	ErrorCount     int
	ErrorMethod    string
	Order          []string
	SetupPatchFunc func(p *Patch)

	errorCount int
}

func (r *testResource) GetCurrentState(ctx context.Context, obj interface{}, deleted bool) (interface{}, error) {
	m := "GetCurrentState"
	r.Order = append(r.Order, fmt.Sprintf("%s(deleted=%t)", m, deleted))

	if r.CancelingStep == m {
		canceledcontext.SetCanceled(ctx)
		if canceledcontext.IsCanceled(ctx) {
			return nil, nil
		}
	}

	if r.returnErrorFor(m) {
		return nil, r.Error
	}

	return nil, nil
}

func (r *testResource) GetDesiredState(ctx context.Context, obj interface{}, deleted bool) (interface{}, error) {
	m := "GetDesiredState"
	r.Order = append(r.Order, fmt.Sprintf("%s(deleted=%t)", m, deleted))

	if r.CancelingStep == m {
		canceledcontext.SetCanceled(ctx)
		if canceledcontext.IsCanceled(ctx) {
			return nil, nil
		}
	}

	if r.returnErrorFor(m) {
		return nil, r.Error
	}

	return nil, nil
}

func (r *testResource) NewPatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error) {
	m := "NewPatch"
	r.Order = append(r.Order, m)

	if r.CancelingStep == m {
		canceledcontext.SetCanceled(ctx)
		if canceledcontext.IsCanceled(ctx) {
			return nil, nil
		}
	}

	if r.returnErrorFor(m) {
		return nil, r.Error
	}

	p := NewPatch()
	if r.SetupPatchFunc != nil {
		r.SetupPatchFunc(p)
	} else {
		p.SetCreateChange("test create data")
		p.SetUpdateChange("test update data")
		p.SetDeleteChange("test delete data")
	}
	return p, nil
}

func (r *testResource) Name() string {
	return "testResource"
}

func (r *testResource) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	m := "ApplyCreatePatch"
	r.Order = append(r.Order, m)

	if r.CancelingStep == m {
		canceledcontext.SetCanceled(ctx)
		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
	}

	if r.returnErrorFor(m) {
		return r.Error
	}

	return nil
}

func (r *testResource) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	m := "ApplyDeletePatch"
	r.Order = append(r.Order, m)

	if r.CancelingStep == m {
		canceledcontext.SetCanceled(ctx)
		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
	}

	if r.returnErrorFor(m) {
		return r.Error
	}

	return nil
}

func (r *testResource) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	m := "ApplyUpdatePatch"
	r.Order = append(r.Order, m)

	if r.CancelingStep == m {
		canceledcontext.SetCanceled(ctx)
		if canceledcontext.IsCanceled(ctx) {
			return nil
		}
	}

	if r.returnErrorFor(m) {
		return r.Error
	}

	return nil
}

func (r *testResource) Underlying() Resource {
	return r
}

func (r *testResource) returnErrorFor(errorMethod string) bool {
	ok := r.Error != nil && r.ErrorCount > r.errorCount && r.ErrorMethod == errorMethod

	if ok {
		r.errorCount++
		return true
	}

	return false
}
