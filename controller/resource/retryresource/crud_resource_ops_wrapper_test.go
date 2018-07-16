package retryresource

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/cenkalti/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/operatorkit/controller"
)

// Test_RetryCRUDResourceOps_ProcessDelete_ResourceOrder_RetryOnError ensures the
// resource's methods are executed as expected when retrying the deletion
// process.
func Test_RetryCRUDResourceOps_ProcessDelete_ResourceOrder_RetryOnError(t *testing.T) {
	testCases := []struct {
		Resource            *testCRUDResource
		ExpectedMethodCalls []string
	}{
		{
			Resource: newTestCRUDResource("r0"),
			ExpectedMethodCalls: []string{
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewDeletePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",
			},
		},
		{
			Resource: newTestCRUDResource("r0").ErrorAt("GetCurrentState", 1),
			ExpectedMethodCalls: []string{
				"r0.GetCurrentState",
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewDeletePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",
			},
		},
		{
			Resource: newTestCRUDResource("r0").ErrorAt("GetCurrentState", 2),
			ExpectedMethodCalls: []string{
				"r0.GetCurrentState",
				"r0.GetCurrentState",
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewDeletePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",
			},
		},
		{
			Resource: newTestCRUDResource("r0").ErrorAt("GetDesiredState", 2),
			ExpectedMethodCalls: []string{
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.GetDesiredState",
				"r0.GetDesiredState",
				"r0.NewDeletePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",
			},
		},
		{
			Resource: newTestCRUDResource("r0").ErrorAt("ApplyDeleteChange", 2),
			ExpectedMethodCalls: []string{
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewDeletePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",
			},
		},
	}

	for i, tc := range testCases {
		rs := []controller.Resource{
			tc.Resource,
		}
		bf := func() backoff.BackOff {
			return &backoff.ZeroBackOff{}
		}

		c := WrapConfig{
			Logger:         microloggertest.New(),
			BackOffFactory: bf,
		}
		wrapped, err := Wrap(rs, c)
		if err != nil {
			t.Fatal("test", i, "expected", nil, "got", err)
		}

		err = controller.ProcessDelete(context.TODO(), nil, wrapped)
		if err != nil {
			t.Fatal("test", i, "expected", nil, "got", err)
		}

		if !reflect.DeepEqual(tc.ExpectedMethodCalls, tc.Resource.MethodCalls()) {
			t.Fatal("test", i, "expected", tc.ExpectedMethodCalls, "got", tc.Resource.MethodCalls())
		}
	}
}

// Test_RetryCRUDResourceOps_ProcessUpdate_ResourceOrder_RetryOnError ensures the
// resource's methods are executed as expected when retrying the update
// process.
func Test_RetryCRUDResourceOps_ProcessUpdate_ResourceOrder_RetryOnError(t *testing.T) {
	testCases := []struct {
		Resource            *testCRUDResource
		ExpectedMethodCalls []string
	}{
		{
			Resource: newTestCRUDResource("r0"),
			ExpectedMethodCalls: []string{
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewUpdatePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",
			},
		},
		{
			Resource: newTestCRUDResource("r0").ErrorAt("GetCurrentState", 1),
			ExpectedMethodCalls: []string{
				"r0.GetCurrentState",
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewUpdatePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",
			},
		},
		{
			Resource: newTestCRUDResource("r0").ErrorAt("GetCurrentState", 2),
			ExpectedMethodCalls: []string{
				"r0.GetCurrentState",
				"r0.GetCurrentState",
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewUpdatePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",
			},
		},
		{
			Resource: newTestCRUDResource("r0").ErrorAt("GetDesiredState", 2),
			ExpectedMethodCalls: []string{
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.GetDesiredState",
				"r0.GetDesiredState",
				"r0.NewUpdatePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",
			},
		},
		{
			Resource: newTestCRUDResource("r0").ErrorAt("ApplyUpdateChange", 2),
			ExpectedMethodCalls: []string{
				"r0.GetCurrentState",
				"r0.GetDesiredState",
				"r0.NewUpdatePatch",
				"r0.ApplyCreateChange",
				"r0.ApplyDeleteChange",
				"r0.ApplyUpdateChange",
				"r0.ApplyUpdateChange",
				"r0.ApplyUpdateChange",
			},
		},
	}

	for i, tc := range testCases {
		rs := []controller.Resource{
			tc.Resource,
		}
		bf := func() backoff.BackOff {
			return &backoff.ZeroBackOff{}
		}

		c := WrapConfig{
			Logger:         microloggertest.New(),
			BackOffFactory: bf,
		}
		wrapped, err := Wrap(rs, c)
		if err != nil {
			t.Fatal("test", i, "expected", nil, "got", err)
		}

		err = controller.ProcessUpdate(context.TODO(), nil, wrapped)
		if err != nil {
			t.Fatal("test", i, "expected", nil, "got", err)
		}

		if !reflect.DeepEqual(tc.ExpectedMethodCalls, tc.Resource.MethodCalls()) {
			t.Fatal("test", i, "expected", tc.ExpectedMethodCalls, "got", tc.Resource.MethodCalls())
		}
	}
}

type testCRUDResource struct {
	*controller.CRUDResource
	ops *testCRUDResourceOps
}

func newTestCRUDResource(name string) *testCRUDResource {
	var err error

	ops := newTestResourceOps(name)

	var crudResource *controller.CRUDResource
	{
		c := controller.CRUDResourceConfig{
			Logger: microloggertest.New(),
			Ops:    ops,
		}

		crudResource, err = controller.NewCRUDResource(c)
		if err != nil {
			panic(fmt.Sprintf("%#v", microerror.Mask(err)))
		}
	}

	return &testCRUDResource{
		CRUDResource: crudResource,
		ops:          ops,
	}
}

func (r *testCRUDResource) ErrorAt(method string, errorCnt int) *testCRUDResource {
	r.ops.ErrorAt(method, errorCnt)
	return r
}

func (r *testCRUDResource) MethodCalls() []string {
	return r.ops.MethodCalls
}

type testCRUDResourceOps struct {
	name        string
	errorMethod string
	errorCnt    int

	MethodCalls []string
}

func newTestResourceOps(name string) *testCRUDResourceOps {
	return &testCRUDResourceOps{
		name: name,
	}
}

func (o *testCRUDResourceOps) Name() string {
	return o.name
}

func (o *testCRUDResourceOps) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	err := o.executeMethod(ctx, "GetCurrentState")
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return nil, nil
}

func (o *testCRUDResourceOps) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	err := o.executeMethod(ctx, "GetDesiredState")
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return nil, nil
}

func (o *testCRUDResourceOps) NewUpdatePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*controller.Patch, error) {
	err := o.executeMethod(ctx, "NewUpdatePatch")
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return newFullPatch(), nil
}

func (o *testCRUDResourceOps) NewDeletePatch(ctx context.Context, obj, currentState, desiredState interface{}) (*controller.Patch, error) {
	err := o.executeMethod(ctx, "NewDeletePatch")
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return newFullPatch(), nil
}

func (o *testCRUDResourceOps) ApplyCreateChange(ctx context.Context, obj, createChange interface{}) error {
	err := o.executeMethod(ctx, "ApplyCreateChange")
	if err != nil {
		return microerror.Mask(err)
	}
	return nil
}

func (o *testCRUDResourceOps) ApplyDeleteChange(ctx context.Context, obj, deleteChange interface{}) error {
	err := o.executeMethod(ctx, "ApplyDeleteChange")
	if err != nil {
		return microerror.Mask(err)
	}
	return nil
}

func (o *testCRUDResourceOps) ApplyUpdateChange(ctx context.Context, obj, updateChange interface{}) error {
	err := o.executeMethod(ctx, "ApplyUpdateChange")
	if err != nil {
		return microerror.Mask(err)
	}
	return nil
}

func (o *testCRUDResourceOps) ErrorAt(method string, errorCnt int) {
	o.errorMethod = method
	o.errorCnt = errorCnt
}

func (o *testCRUDResourceOps) executeMethod(ctx context.Context, method string) error {
	o.MethodCalls = append(o.MethodCalls, o.name+"."+method)

	if o.errorMethod == method && o.errorCnt > 0 {
		o.errorCnt--
		return microerror.Mask(fmt.Errorf("test error from method %s", method))
	}

	return nil
}

// newFullPatch returns Patch filled so all Apply*Change methods are executed.
func newFullPatch() *controller.Patch {
	p := controller.NewPatch()
	p.SetCreateChange("test create data")
	p.SetUpdateChange("test update data")
	p.SetDeleteChange("test delete data")

	return p
}
