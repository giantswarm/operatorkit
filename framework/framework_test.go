package framework

import (
	"reflect"
	"testing"
)

// Test_Framework_ProcessCreate_NoResource ensures there is an error thrown when
// executing ProcessCreate without having any resources provided.
func Test_Framework_ProcessCreate_NoResource(t *testing.T) {
	err := testMustNewFramework(t).ProcessCreate(nil, nil)
	if !IsExecutionFailed(err) {
		t.Fatal("expected", true, "got", false)
	}
}

// Test_Framework_ProcessDelete_NoResource ensures there is an error thrown when
// executing ProcessDelete without having any resources provided.
func Test_Framework_ProcessDelete_NoResource(t *testing.T) {
	err := testMustNewFramework(t).ProcessDelete(nil, nil)
	if !IsExecutionFailed(err) {
		t.Fatal("expected", true, "got", false)
	}
}

// Test_Framework_ProcessCreate_ResourceOrder ensures the resource's methods are
// executed as expected when creating resources.
func Test_Framework_ProcessCreate_ResourceOrder(t *testing.T) {
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

// Test_Framework_ProcessDelete_ResourceOrder ensures the resource's methods are
// executed as expected when deleting resources.
func Test_Framework_ProcessDelete_ResourceOrder(t *testing.T) {
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

type testResource struct {
	Error       error
	ErrorCount  int
	ErrorMethod string
	Order       []string

	errorCount int
}

func (r *testResource) GetCurrentState(obj interface{}) (interface{}, error) {
	m := "GetCurrentState"
	r.Order = append(r.Order, m)

	if r.returnErrorFor(m) {
		return nil, r.Error
	}

	return nil, nil
}

func (r *testResource) GetDesiredState(obj interface{}) (interface{}, error) {
	m := "GetDesiredState"
	r.Order = append(r.Order, m)

	if r.returnErrorFor(m) {
		return nil, r.Error
	}

	return nil, nil
}

func (r *testResource) GetCreateState(obj, currentState, desiredState interface{}) (interface{}, error) {
	m := "GetCreateState"
	r.Order = append(r.Order, m)

	if r.returnErrorFor(m) {
		return nil, r.Error
	}

	return nil, nil
}

func (r *testResource) GetDeleteState(obj, currentState, desiredState interface{}) (interface{}, error) {
	m := "GetDeleteState"
	r.Order = append(r.Order, m)

	if r.returnErrorFor(m) {
		return nil, r.Error
	}

	return nil, nil
}

func (r *testResource) Name() string {
	return "testResource"
}

func (r *testResource) ProcessCreateState(obj, createState interface{}) error {
	m := "ProcessCreateState"
	r.Order = append(r.Order, m)

	if r.returnErrorFor(m) {
		return r.Error
	}

	return nil
}

func (r *testResource) ProcessDeleteState(obj, deleteState interface{}) error {
	m := "ProcessDeleteState"
	r.Order = append(r.Order, m)

	if r.returnErrorFor(m) {
		return r.Error
	}

	return nil
}

func (r *testResource) Underlying() Resource {
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
