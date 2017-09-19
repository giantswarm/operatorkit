package framework

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/operatorkit/framework/cancelercontext"
)

// Test_Framework_ProcessCreate_ResourceOrder ensures the resource's methods are
// executed as expected when creating resources.
func Test_Framework_ProcessCreate_ResourceOrder(t *testing.T) {
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
					"GetCurrentState",
					"GetDesiredState",
					"GetCreateState",
					"ProcessCreateState",
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
					"GetCurrentState",
					"GetDesiredState",
					"GetCreateState",
					"ProcessCreateState",
				},
				{
					"GetCurrentState",
					"GetDesiredState",
					"GetCreateState",
					"ProcessCreateState",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 4 ensures ProcessCreate executes the steps of a single resource in
		// the expected order until it gets canceled.
		{
			ProcessMethod: ProcessCreate,
			Ctx:           cancelercontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{
					CancelingStep: "GetCurrentState",
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState",
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
					"GetCurrentState",
					"GetDesiredState",
					"GetCreateState",
					"ProcessCreateState",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 6 ensures ProcessCreate executes the steps of the first resource in
		// the expected order in case multile resources are given, until the first
		// resource gets canceled.
		{
			ProcessMethod: ProcessCreate,
			Ctx:           cancelercontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{
					CancelingStep: "GetCreateState",
				},
				&testResource{},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState",
					"GetDesiredState",
					"GetCreateState",
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
			Ctx:           cancelercontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{},
				&testResource{
					CancelingStep: "GetCreateState",
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState",
					"GetDesiredState",
					"GetCreateState",
					"ProcessCreateState",
				},
				{
					"GetCurrentState",
					"GetDesiredState",
					"GetCreateState",
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
					"GetCurrentState",
					"GetDesiredState",
					"GetDeleteState",
					"ProcessDeleteState",
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
					"GetCurrentState",
					"GetDesiredState",
					"GetDeleteState",
					"ProcessDeleteState",
				},
				{
					"GetCurrentState",
					"GetDesiredState",
					"GetDeleteState",
					"ProcessDeleteState",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 11 ensures ProcessDelete executes the steps of a single resource in
		// the expected order until it gets canceled.
		{
			ProcessMethod: ProcessDelete,
			Ctx:           cancelercontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{
					CancelingStep: "GetCurrentState",
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState",
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
					"GetCurrentState",
					"GetDesiredState",
					"GetDeleteState",
					"ProcessDeleteState",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 13 ensures ProcessDelete executes the steps of the first resource in
		// the expected order in case multile resources are given, until the first
		// resource gets canceled.
		{
			ProcessMethod: ProcessDelete,
			Ctx:           cancelercontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{
					CancelingStep: "GetDeleteState",
				},
				&testResource{},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState",
					"GetDesiredState",
					"GetDeleteState",
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
			Ctx:           cancelercontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{},
				&testResource{
					CancelingStep: "GetDeleteState",
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState",
					"GetDesiredState",
					"GetDeleteState",
					"ProcessDeleteState",
				},
				{
					"GetCurrentState",
					"GetDesiredState",
					"GetDeleteState",
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
					"GetCurrentState",
					"GetDesiredState",
					"GetUpdateState",
					"ProcessCreateState",
					"ProcessDeleteState",
					"ProcessUpdateState",
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
					"GetCurrentState",
					"GetDesiredState",
					"GetUpdateState",
					"ProcessCreateState",
					"ProcessDeleteState",
					"ProcessUpdateState",
				},
				{
					"GetCurrentState",
					"GetDesiredState",
					"GetUpdateState",
					"ProcessCreateState",
					"ProcessDeleteState",
					"ProcessUpdateState",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 18 ensures ProcessUpdate executes the steps of a single resource in
		// the expected order until it gets canceled.
		{
			ProcessMethod: ProcessUpdate,
			Ctx:           cancelercontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{
					CancelingStep: "GetCurrentState",
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState",
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
					"GetCurrentState",
					"GetDesiredState",
					"GetUpdateState",
					"ProcessCreateState",
					"ProcessDeleteState",
					"ProcessUpdateState",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 20 ensures ProcessUpdate executes the steps of the first resource in
		// the expected order in case multile resources are given, until the first
		// resource gets canceled.
		{
			ProcessMethod: ProcessUpdate,
			Ctx:           cancelercontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{
					CancelingStep: "GetUpdateState",
				},
				&testResource{},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState",
					"GetDesiredState",
					"GetUpdateState",
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
			Ctx:           cancelercontext.NewContext(context.Background(), make(chan struct{})),
			Resources: []Resource{
				&testResource{},
				&testResource{
					CancelingStep: "GetUpdateState",
				},
			},
			ExpectedOrders: [][]string{
				{
					"GetCurrentState",
					"GetDesiredState",
					"GetUpdateState",
					"ProcessCreateState",
					"ProcessDeleteState",
					"ProcessUpdateState",
				},
				{
					"GetCurrentState",
					"GetDesiredState",
					"GetUpdateState",
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
	CancelingStep string
	Error         error
	ErrorCount    int
	ErrorMethod   string
	Order         []string

	errorCount int
}

func (r *testResource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	m := "GetCurrentState"
	r.Order = append(r.Order, m)

	if r.CancelingStep == m {
		canceler, exists := cancelercontext.FromContext(ctx)
		if exists {
			close(canceler)
			return nil, nil
		}
	}

	if r.returnErrorFor(m) {
		return nil, r.Error
	}

	return nil, nil
}

func (r *testResource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	m := "GetDesiredState"
	r.Order = append(r.Order, m)

	if r.CancelingStep == m {
		canceler, exists := cancelercontext.FromContext(ctx)
		if exists {
			close(canceler)
			return nil, nil
		}
	}

	if r.returnErrorFor(m) {
		return nil, r.Error
	}

	return nil, nil
}

func (r *testResource) GetCreateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	m := "GetCreateState"
	r.Order = append(r.Order, m)

	if r.CancelingStep == m {
		canceler, exists := cancelercontext.FromContext(ctx)
		if exists {
			close(canceler)
			return nil, nil
		}
	}

	if r.returnErrorFor(m) {
		return nil, r.Error
	}

	return nil, nil
}

func (r *testResource) GetDeleteState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	m := "GetDeleteState"
	r.Order = append(r.Order, m)

	if r.CancelingStep == m {
		canceler, exists := cancelercontext.FromContext(ctx)
		if exists {
			close(canceler)
			return nil, nil
		}
	}

	if r.returnErrorFor(m) {
		return nil, r.Error
	}

	return nil, nil
}

func (r *testResource) GetUpdateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, interface{}, interface{}, error) {
	m := "GetUpdateState"
	r.Order = append(r.Order, m)

	if r.CancelingStep == m {
		canceler, exists := cancelercontext.FromContext(ctx)
		if exists {
			close(canceler)
			return nil, nil, nil, nil
		}
	}

	if r.returnErrorFor(m) {
		return nil, nil, nil, r.Error
	}

	return nil, nil, nil, nil
}

func (r *testResource) Name() string {
	return "testResource"
}

func (r *testResource) ProcessCreateState(ctx context.Context, obj, createState interface{}) error {
	m := "ProcessCreateState"
	r.Order = append(r.Order, m)

	if r.CancelingStep == m {
		canceler, exists := cancelercontext.FromContext(ctx)
		if exists {
			close(canceler)
			return nil
		}
	}

	if r.returnErrorFor(m) {
		return r.Error
	}

	return nil
}

func (r *testResource) ProcessDeleteState(ctx context.Context, obj, deleteState interface{}) error {
	m := "ProcessDeleteState"
	r.Order = append(r.Order, m)

	if r.CancelingStep == m {
		canceler, exists := cancelercontext.FromContext(ctx)
		if exists {
			close(canceler)
			return nil
		}
	}

	if r.returnErrorFor(m) {
		return r.Error
	}

	return nil
}

func (r *testResource) ProcessUpdateState(ctx context.Context, obj, updateState interface{}) error {
	m := "ProcessUpdateState"
	r.Order = append(r.Order, m)

	if r.CancelingStep == m {
		canceler, exists := cancelercontext.FromContext(ctx)
		if exists {
			close(canceler)
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
