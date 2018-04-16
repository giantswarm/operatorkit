package controller

import (
	"context"
	"testing"
)

func Test_Controller_Resource_PatchNoPanic(t *testing.T) {
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

func (r *testResourcePatchNoPanic) Name() string {
	return "testResourcePatchNoPanic"
}

func (r *testResourcePatchNoPanic) EnsureCreated(ctx context.Context, obj interface{}) error {
	return nil
}

func (r *testResourcePatchNoPanic) EnsureDeleted(ctx context.Context, obj interface{}) error {
	return nil
}
