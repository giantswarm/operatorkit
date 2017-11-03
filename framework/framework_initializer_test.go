package framework

import (
	"context"
	"reflect"
	"testing"

	"github.com/cenk/backoff"
	"github.com/giantswarm/micrologger/microloggertest"
)

func Test_Framework_InitCtxFunc_AddFunc(t *testing.T) {
	testCases := []struct {
		CustomObject  interface{}
		InitCtxFunc   func(ctx context.Context, obj interface{}) (context.Context, error)
		ExpectedOrder []string
	}{
		{
			CustomObject: nil,
			InitCtxFunc: func(ctx context.Context, obj interface{}) (context.Context, error) {
				return ctx, nil
			},
			ExpectedOrder: nil,
		},
		{
			CustomObject: nil,
			InitCtxFunc: func(ctx context.Context, obj interface{}) (context.Context, error) {
				ctx = testInitCtxFuncNewContext(ctx, "foo")
				return ctx, nil
			},
			ExpectedOrder: []string{
				"GetCurrentState",
				"GetDesiredState",
				"NewUpdatePatch",
				"Create",
				"Delete",
				"Update",
			},
		},
	}

	for i, tc := range testCases {
		r := &testInitCtxFuncResource{}

		var f *Framework
		{
			c := DefaultConfig()

			c.BackOff = &backoff.StopBackOff{}
			c.InitCtxFunc = tc.InitCtxFunc
			c.Logger = microloggertest.New()
			c.ResourceRouter = NewDefaultResourceRouter([]Resource{
				r,
			})

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

func Test_Framework_InitCtxFunc_DeleteFunc(t *testing.T) {
	testCases := []struct {
		CustomObject  interface{}
		InitCtxFunc   func(ctx context.Context, obj interface{}) (context.Context, error)
		ExpectedOrder []string
	}{
		{
			CustomObject: nil,
			InitCtxFunc: func(ctx context.Context, obj interface{}) (context.Context, error) {
				return ctx, nil
			},
			ExpectedOrder: nil,
		},
		{
			CustomObject: nil,
			InitCtxFunc: func(ctx context.Context, obj interface{}) (context.Context, error) {
				ctx = testInitCtxFuncNewContext(ctx, "foo")
				return ctx, nil
			},
			ExpectedOrder: []string{
				"GetCurrentState",
				"GetDesiredState",
				"NewDeletePatch",
				"Create",
				"Delete",
				"Update",
			},
		},
	}

	for i, tc := range testCases {
		r := &testInitCtxFuncResource{}

		var f *Framework
		{
			c := DefaultConfig()

			c.BackOff = &backoff.StopBackOff{}
			c.InitCtxFunc = tc.InitCtxFunc
			c.Logger = microloggertest.New()
			c.ResourceRouter = NewDefaultResourceRouter([]Resource{
				r,
			})

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

func Test_Framework_InitCtxFunc_UpdateFunc(t *testing.T) {
	testCases := []struct {
		CustomObject  interface{}
		InitCtxFunc   func(ctx context.Context, obj interface{}) (context.Context, error)
		ExpectedOrder []string
	}{
		{
			CustomObject: nil,
			InitCtxFunc: func(ctx context.Context, obj interface{}) (context.Context, error) {
				return ctx, nil
			},
			ExpectedOrder: nil,
		},
		{
			CustomObject: nil,
			InitCtxFunc: func(ctx context.Context, obj interface{}) (context.Context, error) {
				ctx = testInitCtxFuncNewContext(ctx, "foo")
				return ctx, nil
			},
			ExpectedOrder: []string{
				"GetCurrentState",
				"GetDesiredState",
				"NewUpdatePatch",
				"Create",
				"Delete",
				"Update",
			},
		},
	}

	for i, tc := range testCases {
		r := &testInitCtxFuncResource{}

		var f *Framework
		{
			c := DefaultConfig()

			c.BackOff = &backoff.StopBackOff{}
			c.InitCtxFunc = tc.InitCtxFunc
			c.Logger = microloggertest.New()
			c.ResourceRouter = NewDefaultResourceRouter([]Resource{
				r,
			})

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

type testInitCtxFuncResource struct {
	Order []string
}

func (r *testInitCtxFuncResource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	_, ok := testInitCtxFuncFromContext(ctx)
	if ok {
		m := "GetCurrentState"
		r.Order = append(r.Order, m)
	}

	return nil, nil
}

func (r *testInitCtxFuncResource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	_, ok := testInitCtxFuncFromContext(ctx)
	if ok {
		m := "GetDesiredState"
		r.Order = append(r.Order, m)
	}

	return nil, nil
}

func (r *testInitCtxFuncResource) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error) {
	_, ok := testInitCtxFuncFromContext(ctx)
	if ok {
		m := "NewUpdatePatch"
		r.Order = append(r.Order, m)
	}

	p := NewPatch()
	p.SetCreateChange("test create state")
	p.SetUpdateChange("test update state")
	p.SetDeleteChange("test delete state")
	return p, nil
}

func (r *testInitCtxFuncResource) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*Patch, error) {
	_, ok := testInitCtxFuncFromContext(ctx)
	if ok {
		m := "NewDeletePatch"
		r.Order = append(r.Order, m)
	}

	p := NewPatch()
	p.SetCreateChange("test create state")
	p.SetUpdateChange("test update state")
	p.SetDeleteChange("test delete state")
	return p, nil
}

func (r *testInitCtxFuncResource) Name() string {
	return "testInitCtxFuncResource"
}

func (r *testInitCtxFuncResource) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	_, ok := testInitCtxFuncFromContext(ctx)
	if ok {
		m := "Create"
		r.Order = append(r.Order, m)
	}

	return nil
}

func (r *testInitCtxFuncResource) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	_, ok := testInitCtxFuncFromContext(ctx)
	if ok {
		m := "Delete"
		r.Order = append(r.Order, m)
	}

	return nil
}

func (r *testInitCtxFuncResource) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	_, ok := testInitCtxFuncFromContext(ctx)
	if ok {
		m := "Update"
		r.Order = append(r.Order, m)
	}

	return nil
}

func (r *testInitCtxFuncResource) Underlying() Resource {
	return r
}

type key string

var testInitCtxFuncKey key = "testinitiaqlizer"

func testInitCtxFuncNewContext(ctx context.Context, v interface{}) context.Context {
	return context.WithValue(ctx, testInitCtxFuncKey, v)
}

func testInitCtxFuncFromContext(ctx context.Context) (interface{}, bool) {
	v, ok := ctx.Value(testInitCtxFuncKey).(interface{})
	return v, ok
}
