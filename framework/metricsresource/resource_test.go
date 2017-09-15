package metricsresource

import (
	"reflect"
	"testing"

	"github.com/giantswarm/operatorkit/framework"
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
	config.Namespace = t.Name()
	wrapped, err := Wrap(rs, config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	err = framework.ProcessCreate(nil, wrapped)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	e := []string{
		"GetCurrentState",
		"GetDesiredState",
		"GetCreateState",
		"ProcessCreateState",
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
	config.Namespace = t.Name()
	wrapped, err := Wrap(rs, config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	err = framework.ProcessDelete(nil, wrapped)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	e := []string{
		"GetCurrentState",
		"GetDesiredState",
		"GetDeleteState",
		"ProcessDeleteState",
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
	config.Namespace = t.Name()
	wrapped, err := Wrap(rs, config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	err = framework.ProcessUpdate(nil, wrapped)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	e := []string{
		"GetCurrentState",
		"GetDesiredState",
		"GetUpdateState",
		"ProcessCreateState",
		"ProcessDeleteState",
		"ProcessUpdateState",
	}
	if !reflect.DeepEqual(e, tr.Order) {
		t.Fatal("expected", e, "got", tr.Order)
	}
}

type testResource struct {
	Order []string
}

func (r *testResource) GetCurrentState(obj interface{}) (interface{}, error) {
	m := "GetCurrentState"
	r.Order = append(r.Order, m)

	return nil, nil
}

func (r *testResource) GetDesiredState(obj interface{}) (interface{}, error) {
	m := "GetDesiredState"
	r.Order = append(r.Order, m)

	return nil, nil
}

func (r *testResource) GetCreateState(obj, currentState, desiredState interface{}) (interface{}, error) {
	m := "GetCreateState"
	r.Order = append(r.Order, m)

	return nil, nil
}

func (r *testResource) GetDeleteState(obj, currentState, desiredState interface{}) (interface{}, error) {
	m := "GetDeleteState"
	r.Order = append(r.Order, m)

	return nil, nil
}

func (r *testResource) GetUpdateState(obj, currentState, desiredState interface{}) (interface{}, interface{}, interface{}, error) {
	m := "GetUpdateState"
	r.Order = append(r.Order, m)

	return nil, nil, nil, nil
}

func (r *testResource) Name() string {
	return "testResource"
}

func (r *testResource) ProcessCreateState(obj, createState interface{}) error {
	m := "ProcessCreateState"
	r.Order = append(r.Order, m)

	return nil
}

func (r *testResource) ProcessDeleteState(obj, deleteState interface{}) error {
	m := "ProcessDeleteState"
	r.Order = append(r.Order, m)

	return nil
}

func (r *testResource) ProcessUpdateState(obj, updateState interface{}) error {
	m := "ProcessUpdateState"
	r.Order = append(r.Order, m)

	return nil
}

func (r *testResource) Underlying() framework.Resource {
	return r
}
