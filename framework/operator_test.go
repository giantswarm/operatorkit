package framework

import (
	"reflect"
	"testing"
)

type testResource struct {
	Order []string
}

func (r *testResource) GetCurrentState(obj interface{}) (interface{}, error) {
	r.Order = append(r.Order, "GetCurrentState")
	return nil, nil
}

func (r *testResource) GetDesiredState(obj interface{}) (interface{}, error) {
	r.Order = append(r.Order, "GetDesiredState")
	return nil, nil
}

func (r *testResource) GetCreateState(obj, currentState, desiredState interface{}) (interface{}, error) {
	r.Order = append(r.Order, "GetCreateState")
	return nil, nil
}

func (r *testResource) GetDeleteState(obj, currentState, desiredState interface{}) (interface{}, error) {
	r.Order = append(r.Order, "GetDeleteState")
	return nil, nil
}

func (r *testResource) ProcessCreateState(obj, createState interface{}) error {
	r.Order = append(r.Order, "ProcessCreateState")
	return nil
}

func (r *testResource) ProcessDeleteState(obj, deleteState interface{}) error {
	r.Order = append(r.Order, "ProcessDeleteState")
	return nil
}

// Test_Operator_ProcessCreate_NoResource ensures there is an error thrown when
// executing ProcessCreate without having any resources provided.
func Test_Operator_ProcessCreate_NoResource(t *testing.T) {
	err := testMustNewFramework(t).ProcessCreate(nil, nil)
	if !IsExecutionFailed(err) {
		t.Fatal("expected", true, "got", false)
	}
}

// Test_Operator_ProcessDelete_NoResource ensures there is an error thrown when
// executing ProcessDelete without having any resources provided.
func Test_Operator_ProcessDelete_NoResource(t *testing.T) {
	err := testMustNewFramework(t).ProcessDelete(nil, nil)
	if !IsExecutionFailed(err) {
		t.Fatal("expected", true, "got", false)
	}
}

// Test_Operator_ProcessCreate_ResourceOrder ensures the resource's methods are
// executed as expected when creating resources.
func Test_Operator_ProcessCreate_ResourceOrder(t *testing.T) {
	tr := &testResource{}
	rs := []Resource{
		tr,
	}

	err := testMustNewFramework(t).ProcessCreate(nil, rs)
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

// Test_Operator_ProcessDelete_ResourceOrder ensures the resource's methods are
// executed as expected when deleting resources.
func Test_Operator_ProcessDelete_ResourceOrder(t *testing.T) {
	tr := &testResource{}
	rs := []Resource{
		tr,
	}

	err := testMustNewFramework(t).ProcessDelete(nil, rs)
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

func testMustNewFramework(t *testing.T) *Framework {
	frameworkConfig := DefaultConfig()
	newFramework, err := New(frameworkConfig)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	return newFramework
}
