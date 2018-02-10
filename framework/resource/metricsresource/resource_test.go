package metricsresource

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/operatorkit/framework"
	"github.com/giantswarm/operatorkit/framework/resource/internal"
)

func Test_Wrapper(t *testing.T) {
	// This won't compile if the *Resource doesn't implement Wrapper
	// interface.
	var _ internal.Wrapper = &Resource{}
}

// Test_MetricsResource_ProcessDelete_ResourceOrder ensures the resource's
// methods are executed as expected when deleting resources using the wrapping
// prometheus resource.
func Test_MetricsResource_ProcessDelete_ResourceOrder(t *testing.T) {
	tr := &testResource{}
	rs := []framework.Resource{
		tr,
	}

	c := WrapConfig{
		Name: t.Name(),
	}
	wrapped, err := Wrap(rs, c)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	err = framework.ProcessDelete(context.TODO(), nil, wrapped)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	e := []string{
		"GetCurrentState",
		"GetDesiredState",
		"NewDeletePatch",
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

	c := WrapConfig{
		Name: t.Name(),
	}
	wrapped, err := Wrap(rs, c)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	err = framework.ProcessUpdate(context.TODO(), nil, wrapped)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	e := []string{
		"GetCurrentState",
		"GetDesiredState",
		"NewUpdatePatch",
		"ApplyCreatePatch",
		"ApplyDeletePatch",
		"ApplyUpdatePatch",
	}
	if !reflect.DeepEqual(e, tr.Order) {
		t.Fatal("expected", e, "got", tr.Order)
	}
}

type testResource struct {
	Order []string
}

func (r *testResource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	m := "GetCurrentState"
	r.Order = append(r.Order, m)

	return nil, nil
}

func (r *testResource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	m := "GetDesiredState"
	r.Order = append(r.Order, m)

	return nil, nil
}

func (r *testResource) NewUpdatePatch(ctx context.Context, obj, cur, des interface{}) (*framework.Patch, error) {
	m := "NewUpdatePatch"
	r.Order = append(r.Order, m)

	p := framework.NewPatch()
	p.SetCreateChange("test create data")
	p.SetUpdateChange("test update data")
	p.SetDeleteChange("test delete data")
	return p, nil
}

func (r *testResource) NewDeletePatch(ctx context.Context, obj, cur, des interface{}) (*framework.Patch, error) {
	m := "NewDeletePatch"
	r.Order = append(r.Order, m)

	p := framework.NewPatch()
	p.SetCreateChange("test create data")
	p.SetUpdateChange("test update data")
	p.SetDeleteChange("test delete data")
	return p, nil
}

func (r *testResource) Name() string {
	return "testResource"
}

func (r *testResource) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
	m := "ApplyCreatePatch"
	r.Order = append(r.Order, m)

	return nil
}

func (r *testResource) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	m := "ApplyDeletePatch"
	r.Order = append(r.Order, m)

	return nil
}

func (r *testResource) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	m := "ApplyUpdatePatch"
	r.Order = append(r.Order, m)

	return nil
}
