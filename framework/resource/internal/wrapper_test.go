package internal

import (
	"context"
	"testing"

	"github.com/giantswarm/operatorkit/framework"
)

// Test_Underlying_Wrapped tests resource unwrapping by Underlying.
func Test_Underlying_Wrapped(t *testing.T) {
	testCases := []struct {
		Resource framework.Resource
	}{
		// Test 0.
		{
			Resource: &testWrappingResource{
				resource: &testNoopResource{},
			},
		},
		// Test 1.
		{
			Resource: &testWrappingResource{
				resource: &testWrappingResource{
					resource: &testNoopResource{},
				},
			},
		},
		// Test 2.
		{
			Resource: &testWrappingResource{
				resource: &testWrappingResource{
					resource: &testWrappingResource{
						resource: &testNoopResource{},
					},
				},
			},
		},
		// Test 3.
		{
			Resource: &testWrappingResource{
				resource: &testWrappingResource{
					resource: &testWrappingResource{
						resource: &testWrappingResource{
							resource: &testNoopResource{},
						},
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		if tc.Resource.Name() != "testWrappingResource" {
			t.Fatalf("test %d: expected %q resource, got %q", i, "testWrappingResource", tc.Resource.Name())
		}

		underlying, err := Underlying(tc.Resource)
		if err != nil {
			t.Fatalf("test %d: unexpected error = %#v", i, err)
		}

		if underlying.Name() != "testNoopResource" {
			t.Fatalf("test %d: expected %q resource, got %q", i, "testNoopResource", underlying.Name())
		}
	}
}

// Test_Underlying_NonWrapped tests Underlying returns the same resource when
// it isn't a Wrapper.
func Test_Underlying_NonWrapped(t *testing.T) {
	r := &testNoopResource{}

	underlying, err := Underlying(r)
	if err != nil {
		t.Fatalf("unexpected error = %#v", err)
	}

	if underlying != r {
		t.Fatalf("expected %#v, got %#v", r, underlying)
	}
}

// Test_Underlying_Loop tests if Underlying returns an error when there is
// an infinite loop.
func Test_Underlying_Loop(t *testing.T) {
	testCases := []struct {
		ResourceFunc func() framework.Resource
	}{
		// Test 0. r1 -> r1.
		{
			ResourceFunc: func() framework.Resource {
				r1 := &testWrappingResource{}

				r1.resource = r1

				return r1
			},
		},
		// Test 1. r1 -> r2 -> r3 -> r1.
		{
			ResourceFunc: func() framework.Resource {
				r1 := &testWrappingResource{}
				r2 := &testWrappingResource{}
				r3 := &testWrappingResource{}

				r1.resource = r2
				r2.resource = r3
				r3.resource = r1

				return r1
			},
		},
		// Test 1. r1 -> r2 -> r3 -> r4 -> r2.
		{
			ResourceFunc: func() framework.Resource {
				r1 := &testWrappingResource{}
				r2 := &testWrappingResource{}
				r3 := &testWrappingResource{}
				r4 := &testWrappingResource{}

				r1.resource = r2
				r2.resource = r3
				r3.resource = r4
				r4.resource = r2

				return r1
			},
		},
	}

	for i, tc := range testCases {
		_, err := Underlying(tc.ResourceFunc())

		if !IsLoopDetected(err) {
			t.Fatalf("test %d: expected %v, got %v", i, loopDetectedError, err)
		}
	}
}

type testNoopResource struct{}

func (r *testNoopResource) Name() string {
	return "testNoopResource"
}

func (r *testNoopResource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

func (r *testNoopResource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

func (r *testNoopResource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	return nil, nil
}

func (r *testNoopResource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	return nil, nil
}

func (r *testNoopResource) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	return nil
}

func (r *testNoopResource) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	return nil
}

func (r *testNoopResource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	return nil
}

type testWrappingResource struct {
	resource framework.Resource
}

func (r *testWrappingResource) Wrapped() framework.Resource {
	return r.resource
}

func (r *testWrappingResource) Name() string {
	return "testWrappingResource"
}

func (r *testWrappingResource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

func (r *testWrappingResource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	return nil, nil
}

func (r *testWrappingResource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	return nil, nil
}

func (r *testWrappingResource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
	return nil, nil
}

func (r *testWrappingResource) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	return nil
}

func (r *testWrappingResource) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	return nil
}

func (r *testWrappingResource) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	return nil
}
