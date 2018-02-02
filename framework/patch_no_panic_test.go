package framework

import (
	"context"
	"testing"
)

func Test_Framework_Resource_PatchNoPanic(t *testing.T) {
	testCases := []struct {
		ProcessMethod func(ctx context.Context, obj interface{}, rs []Resource) error
		Resources     []Resource
		ErrorMatcher  func(err error) bool
	}{
		// Test 0 ensures ProcessDelete returns an error in case no resources are
		// provided.
		{
			ProcessMethod: ProcessDelete,
			Resources:     nil,
			ErrorMatcher:  IsExecutionFailed,
		},

		// Test 1 ensures ProcessDelete does not panic when executing a single
		// resource.
		{
			ProcessMethod: ProcessDelete,
			Resources: []Resource{
				&testResourcePatchNoPanic{},
			},
			ErrorMatcher: nil,
		},

		// Test 2 ensures ProcessDelete does not panic when executing two resources.
		{
			ProcessMethod: ProcessDelete,
			Resources: []Resource{
				&testResourcePatchNoPanic{},
				&testResourcePatchNoPanic{},
			},
			ErrorMatcher: nil,
		},

		// Test 3 ensures ProcessUpdate returns an error in case no resources are
		// provided.
		{
			ProcessMethod: ProcessUpdate,
			Resources:     nil,
			ErrorMatcher:  IsExecutionFailed,
		},

		// Test 4 ensures ProcessUpdate does not panic when executing a single
		// resource.
		{
			ProcessMethod: ProcessUpdate,
			Resources: []Resource{
				&testResourcePatchNoPanic{},
			},
			ErrorMatcher: nil,
		},

		// Test 5 ensures ProcessUpdate does not panic when executing two resources.
		{
			ProcessMethod: ProcessUpdate,
			Resources: []Resource{
				&testResourcePatchNoPanic{},
				&testResourcePatchNoPanic{},
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		err := tc.ProcessMethod(context.TODO(), nil, tc.Resources)
		if err != nil {
			if tc.ErrorMatcher == nil {
				t.Fatal("test", i, "expected", nil, "got", err)
			} else if !tc.ErrorMatcher(err) {
				t.Fatal("test", i, "expected", true, "got", false)
			}
		}
	}
}

type testResourcePatchNoPanic struct {
}

func (r *testResourcePatchNoPanic) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

func (r *testResourcePatchNoPanic) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

// NewUpdatePatch returns nil for the *Patch return value, thus making sure the
// resource reconciliation does still work and e.g. not panic.
func (r *testResourcePatchNoPanic) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error) {
	return nil, nil
}

// NewDeletePatch returns nil for the *Patch return value, thus making sure the
// resource reconciliation does still work and e.g. not panic.
func (r *testResourcePatchNoPanic) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error) {
	return nil, nil
}

func (r *testResourcePatchNoPanic) Name() string {
	return "testResourcePatchNoPanic"
}

func (r *testResourcePatchNoPanic) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	return nil
}

func (r *testResourcePatchNoPanic) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	return nil
}

func (r *testResourcePatchNoPanic) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	return nil
}

func (r *testResourcePatchNoPanic) Underlying() Resource {
	return r
}
