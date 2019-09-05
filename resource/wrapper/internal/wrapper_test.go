package internal

import (
	"context"
	"testing"

	"github.com/giantswarm/operatorkit/resource"
)

// Test_Underlying_Wrapped tests resource unwrapping by Underlying.
func Test_Underlying_Wrapped(t *testing.T) {
	testCases := []struct {
		Resource resource.Interface
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
		ResourceFunc func() resource.Interface
	}{
		// Test 0. r1 -> r1.
		{
			ResourceFunc: func() resource.Interface {
				r1 := &testWrappingResource{}

				r1.resource = r1

				return r1
			},
		},
		// Test 1. r1 -> r2 -> r3 -> r1.
		{
			ResourceFunc: func() resource.Interface {
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
			ResourceFunc: func() resource.Interface {
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

func (r *testNoopResource) EnsureCreated(ctx context.Context, obj interface{}) error {
	return nil
}

func (r *testNoopResource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	return nil
}

type testWrappingResource struct {
	resource resource.Interface
}

func (r *testWrappingResource) Wrapped() resource.Interface {
	return r.resource
}

func (r *testWrappingResource) Name() string {
	return "testWrappingResource"
}

func (r *testWrappingResource) EnsureCreated(ctx context.Context, obj interface{}) error {
	return nil
}

func (r *testWrappingResource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	return nil
}
