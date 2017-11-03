package metricsresource

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/context/canceledcontext"
)

// Test_MetricsResource_ProcessCreate_ResourceOrder ensures the resource's
// methods are executed as expected when creating resources using the wrapping
// prometheus resource.
func Test_MetricsResource_ProcessCreate_ResourceOrder(t *testing.T) {
	tr := &testResource{}
	rs := []framework.Resource{
		tr,
	}

	config := DefaultWrapConfig()
	config.Name = t.Name()
	wrapped, err := Wrap(rs, config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	err = framework.ProcessCreate(context.TODO(), nil, wrapped)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	e := []string{
		"GetCurrentState(deleted=false)",
		"GetDesiredState(deleted=false)",
		"NewPatch",
		"ApplyCreatePatch",
		"ApplyDeletePatch",
		"ApplyUpdatePatch",
	}
	if !reflect.DeepEqual(e, tr.Order) {
		t.Fatal("expected", e, "got", tr.Order)
	}
}

// Test_MetricsResource_ProcessDelete_ResourceOrder ensures the resource's
// methods are executed as expected when deleting resources using the wrapping
// prometheus resource.
func Test_MetricsResource_ProcessDelete_ResourceOrder(t *testing.T) {
	tr := &testResource{}
	rs := []framework.Resource{
		tr,
	}

	config := DefaultWrapConfig()
	config.Name = t.Name()
	wrapped, err := Wrap(rs, config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	err = framework.ProcessDelete(context.TODO(), nil, wrapped)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	e := []string{
		"GetCurrentState(deleted=true)",
		"GetDesiredState(deleted=true)",
		"NewPatch",
		"ApplyCreatePatch",
		"ApplyDeletePatch",
		"ApplyUpdatePatch",
	}
	if !reflect.DeepEqual(e, tr.Order) {
		t.Fatal("expected", e, "got", tr.Order)
	}
}

// Test_MetricsResource_ProcessUpdate_ResourceOrder ensures the resource's
// methods are executed as expected when updating resources using the wrapping
// prometheus resource.
func Test_MetricsResource_ProcessUpdate_ResourceOrder(t *testing.T) {
	tr := &testResource{}
	rs := []framework.Resource{
		tr,
	}

	config := DefaultWrapConfig()
	config.Name = t.Name()
	wrapped, err := Wrap(rs, config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	err = framework.ProcessUpdate(context.TODO(), nil, wrapped)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	e := []string{
		"GetCurrentState(deleted=false)",
		"GetDesiredState(deleted=false)",
		"NewPatch",
		"ApplyCreatePatch",
		"ApplyDeletePatch",
		"ApplyUpdatePatch",
	}
	if !reflect.DeepEqual(e, tr.Order) {
		t.Fatal("expected", e, "got", tr.Order)
	}
}

type testResource struct {
	CancelingStep  string
	Error          error
	ErrorCount     int
	ErrorMethod    string
	Order          []string
	SetupPatchFunc func(p *framework.Patch)

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

func (r *testResource) NewPatch(ctx context.Context, obj, currentState, desiredState interface{}) (*framework.Patch, error) {
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

	p := framework.NewPatch()
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

func (r *testResource) Underlying() framework.Resource {
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
