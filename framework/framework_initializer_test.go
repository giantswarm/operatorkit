package framework

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/micrologger/microloggertest"
)

func Test_Framework_Initializer_AddFunc(t *testing.T) {
	testCases := []struct {
		CustomObject  interface{}
		Initializer   func(ctx context.Context, obj interface{}) (context.Context, error)
		ExpectedOrder []string
	}{
		{
			CustomObject: nil,
			Initializer: func(ctx context.Context, obj interface{}) (context.Context, error) {
				return ctx, nil
			},
			ExpectedOrder: nil,
		},
		{
			CustomObject: nil,
			Initializer: func(ctx context.Context, obj interface{}) (context.Context, error) {
				ctx = testInitializerNewContext(ctx, "foo")
				return ctx, nil
			},
			ExpectedOrder: []string{
				"GetCurrentState",
				"GetDesiredState",
				"GetCreateState",
				"ProcessCreateState",
				"GetCurrentState",
				"GetDesiredState",
				"GetUpdateState",
				"ProcessCreateState",
				"ProcessDeleteState",
				"ProcessUpdateState",
			},
		},
	}

	for i, tc := range testCases {
		r := &testInitilizerResource{}

		var f *Framework
		{
			c := DefaultConfig()

			c.Initializer = tc.Initializer
			c.Logger = microloggertest.New()
			c.Resources = []Resource{
				r,
			}

			var err error
			f, err = New(c)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
		}

		f.AddFunc(tc.CustomObject)

		if !reflect.DeepEqual(tc.ExpectedOrder, r.Order) {
			t.Fatal("test", i+1, "expected", tc.ExpectedOrder, "got", r.Order)
		}
	}
}

func Test_Framework_Initializer_DeleteFunc(t *testing.T) {
	testCases := []struct {
		CustomObject  interface{}
		Initializer   func(ctx context.Context, obj interface{}) (context.Context, error)
		ExpectedOrder []string
	}{
		{
			CustomObject: nil,
			Initializer: func(ctx context.Context, obj interface{}) (context.Context, error) {
				return ctx, nil
			},
			ExpectedOrder: nil,
		},
		{
			CustomObject: nil,
			Initializer: func(ctx context.Context, obj interface{}) (context.Context, error) {
				ctx = testInitializerNewContext(ctx, "foo")
				return ctx, nil
			},
			ExpectedOrder: []string{
				"GetCurrentState",
				"GetDesiredState",
				"GetDeleteState",
				"ProcessDeleteState",
			},
		},
	}

	for i, tc := range testCases {
		r := &testInitilizerResource{}

		var f *Framework
		{
			c := DefaultConfig()

			c.Initializer = tc.Initializer
			c.Logger = microloggertest.New()
			c.Resources = []Resource{
				r,
			}

			var err error
			f, err = New(c)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
		}

		f.DeleteFunc(tc.CustomObject)

		if !reflect.DeepEqual(tc.ExpectedOrder, r.Order) {
			t.Fatal("test", i+1, "expected", tc.ExpectedOrder, "got", r.Order)
		}
	}
}

func Test_Framework_Initializer_UpdateFunc(t *testing.T) {
	testCases := []struct {
		CustomObject  interface{}
		Initializer   func(ctx context.Context, obj interface{}) (context.Context, error)
		ExpectedOrder []string
	}{
		{
			CustomObject: nil,
			Initializer: func(ctx context.Context, obj interface{}) (context.Context, error) {
				return ctx, nil
			},
			ExpectedOrder: nil,
		},
		{
			CustomObject: nil,
			Initializer: func(ctx context.Context, obj interface{}) (context.Context, error) {
				ctx = testInitializerNewContext(ctx, "foo")
				return ctx, nil
			},
			ExpectedOrder: []string{
				"GetCurrentState",
				"GetDesiredState",
				"GetCreateState",
				"ProcessCreateState",
				"GetCurrentState",
				"GetDesiredState",
				"GetUpdateState",
				"ProcessCreateState",
				"ProcessDeleteState",
				"ProcessUpdateState",
			},
		},
	}

	for i, tc := range testCases {
		r := &testInitilizerResource{}

		var f *Framework
		{
			c := DefaultConfig()

			c.Initializer = tc.Initializer
			c.Logger = microloggertest.New()
			c.Resources = []Resource{
				r,
			}

			var err error
			f, err = New(c)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
		}

		f.UpdateFunc(nil, tc.CustomObject)

		if !reflect.DeepEqual(tc.ExpectedOrder, r.Order) {
			t.Fatal("test", i+1, "expected", tc.ExpectedOrder, "got", r.Order)
		}
	}
}

type testInitilizerResource struct {
	Order []string
}

func (r *testInitilizerResource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	_, ok := testInitializerFromContext(ctx)
	if ok {
		m := "GetCurrentState"
		r.Order = append(r.Order, m)
	}

	return nil, nil
}

func (r *testInitilizerResource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	_, ok := testInitializerFromContext(ctx)
	if ok {
		m := "GetDesiredState"
		r.Order = append(r.Order, m)
	}

	return nil, nil
}

func (r *testInitilizerResource) GetCreateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	_, ok := testInitializerFromContext(ctx)
	if ok {
		m := "GetCreateState"
		r.Order = append(r.Order, m)
	}

	return nil, nil
}

func (r *testInitilizerResource) GetDeleteState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	_, ok := testInitializerFromContext(ctx)
	if ok {
		m := "GetDeleteState"
		r.Order = append(r.Order, m)
	}

	return nil, nil
}

func (r *testInitilizerResource) GetUpdateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, interface{}, interface{}, error) {
	_, ok := testInitializerFromContext(ctx)
	if ok {
		m := "GetUpdateState"
		r.Order = append(r.Order, m)
	}

	return nil, nil, nil, nil
}

func (r *testInitilizerResource) Name() string {
	return "testInitilizerResource"
}

func (r *testInitilizerResource) ProcessCreateState(ctx context.Context, obj, createState interface{}) error {
	_, ok := testInitializerFromContext(ctx)
	if ok {
		m := "ProcessCreateState"
		r.Order = append(r.Order, m)
	}

	return nil
}

func (r *testInitilizerResource) ProcessDeleteState(ctx context.Context, obj, deleteState interface{}) error {
	_, ok := testInitializerFromContext(ctx)
	if ok {
		m := "ProcessDeleteState"
		r.Order = append(r.Order, m)
	}

	return nil
}

func (r *testInitilizerResource) ProcessUpdateState(ctx context.Context, obj, updateState interface{}) error {
	_, ok := testInitializerFromContext(ctx)
	if ok {
		m := "ProcessUpdateState"
		r.Order = append(r.Order, m)
	}

	return nil
}

func (r *testInitilizerResource) Underlying() Resource {
	return r
}

type key string

var testInitializerKey key = "testinitiaqlizer"

func testInitializerNewContext(ctx context.Context, v interface{}) context.Context {
	return context.WithValue(ctx, testInitializerKey, v)
}

func testInitializerFromContext(ctx context.Context) (interface{}, bool) {
	v, ok := ctx.Value(testInitializerKey).(interface{})
	return v, ok
}
