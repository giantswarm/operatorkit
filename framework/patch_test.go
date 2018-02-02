package framework

import (
	"context"
	"testing"
)

func Test_Framework_ResourcePatchNoPanic(t *testing.T) {
	testCases := []struct {
		ProcessMethod func(ctx context.Context, obj interface{}, rs []Resource) error
		Ctx           context.Context
		Resources     []Resource
		ErrorMatcher  func(err error) bool
	}{
		// Test 0 ensures ProcessDelete returns an error in case no resources are
		// provided.
		{
			ProcessMethod: ProcessDelete,
			Ctx:           context.TODO(),
			Resources:     nil,
			ErrorMatcher:  IsExecutionFailed,
		},

		// Test 1 ensures ProcessDelete executes the steps of a single resource in
		// the expected order.
		{
			ProcessMethod: ProcessDelete,
			Ctx:           context.TODO(),
			Resources: []Resource{
				&testPatchResource{},
			},
			ErrorMatcher: nil,
		},

		// Test 2 ensures ProcessDelete executes the steps of multile resources in
		// the expected order.
		{
			ProcessMethod: ProcessDelete,
			Ctx:           context.TODO(),
			Resources: []Resource{
				&testPatchResource{},
				&testPatchResource{},
			},
			ErrorMatcher: nil,
		},

		// Test 3 ensures ProcessDelete returns an error in case no resources are
		// provided.
		{
			ProcessMethod: ProcessDelete,
			Ctx:           context.TODO(),
			Resources:     nil,
			ErrorMatcher:  IsExecutionFailed,
		},

		// Test 4 ensures ProcessDelete executes the steps of a single resource in
		// the expected order.
		{
			ProcessMethod: ProcessDelete,
			Ctx:           context.TODO(),
			Resources: []Resource{
				&testPatchResource{},
			},
			ErrorMatcher: nil,
		},

		// Test 5 ensures ProcessDelete executes the steps of multile resources in
		// the expected order.
		{
			ProcessMethod: ProcessDelete,
			Ctx:           context.TODO(),
			Resources: []Resource{
				&testPatchResource{},
				&testPatchResource{},
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		err := tc.ProcessMethod(tc.Ctx, nil, tc.Resources)
		if err != nil {
			if tc.ErrorMatcher == nil {
				t.Fatal("test", i+1, "expected", nil, "got", err)
			} else if !tc.ErrorMatcher(err) {
				t.Fatal("test", i+1, "expected", true, "got", false)
			}
		}
	}
}

type testPatchResource struct {
	SetupPatchFunc func(p *Patch)
}

func (r *testPatchResource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

func (r *testPatchResource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

// NewUpdatePatch returns nil for the *Patch return value, thus making sure the
// resource reconciliation does still work and e.g. not panic.
func (r *testPatchResource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error) {
	return nil, nil
}

// NewDeletePatch returns nil for the *Patch return value, thus making sure the
// resource reconciliation does still work and e.g. not panic.
func (r *testPatchResource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error) {
	return nil, nil
}

func (r *testPatchResource) Name() string {
	return "testPatchResource"
}

func (r *testPatchResource) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	return nil
}

func (r *testPatchResource) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	return nil
}

func (r *testPatchResource) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	return nil
}

func (r *testPatchResource) Underlying() Resource {
	return r
}
