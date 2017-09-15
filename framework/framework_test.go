package framework

import (
	"context"
	"reflect"
	"testing"
)

// Test_Framework_ProcessCreate_NoResource ensures there is an error thrown when
// executing ProcessCreate without having any resources provided.
func Test_Framework_ProcessCreate_NoResource(t *testing.T) {
	err := ProcessCreate(context.TODO(), nil, nil)
	if !IsExecutionFailed(err) {
		t.Fatal("expected", true, "got", false)
	}
}

// Test_Framework_ProcessDelete_NoResource ensures there is an error thrown when
// executing ProcessDelete without having any resources provided.
func Test_Framework_ProcessDelete_NoResource(t *testing.T) {
	err := ProcessDelete(context.TODO(), nil, nil)
	if !IsExecutionFailed(err) {
		t.Fatal("expected", true, "got", false)
	}
}

// Test_Framework_ProcessUpdate_NoResource ensures there is an error thrown when
// executing ProcessUpdate without having any resources provided.
func Test_Framework_ProcessUpdate_NoResource(t *testing.T) {
	err := ProcessUpdate(context.TODO(), nil, nil)
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

	err := ProcessCreate(context.TODO(), nil, rs)
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

	err := ProcessDelete(context.TODO(), nil, rs)
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

// Test_Framework_ProcessUpdate_ResourceOrder ensures the resource's methods are
// executed as expected when updating resources.
func Test_Framework_ProcessUpdate_ResourceOrder(t *testing.T) {
	tr := &testResource{}
	rs := []Resource{
		tr,
	}

	err := ProcessUpdate(context.TODO(), nil, rs)
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
	Error       error
	ErrorCount  int
	ErrorMethod string
	Order       []string

	errorCount int
}

func (r *testResource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
	m := "GetCurrentState"
	r.Order = append(r.Order, m)

	if r.returnErrorFor(m) {
		return nil, r.Error
	}

	return nil, nil
}

func (r *testResource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
	m := "GetDesiredState"
	r.Order = append(r.Order, m)

	if r.returnErrorFor(m) {
		return nil, r.Error
	}

	return nil, nil
}

func (r *testResource) GetCreateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	m := "GetCreateState"
	r.Order = append(r.Order, m)

	if r.returnErrorFor(m) {
		return nil, r.Error
	}

	return nil, nil
}

func (r *testResource) GetDeleteState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, error) {
	m := "GetDeleteState"
	r.Order = append(r.Order, m)

	if r.returnErrorFor(m) {
		return nil, r.Error
	}

	return nil, nil
}

func (r *testResource) GetUpdateState(ctx context.Context, obj, currentState, desiredState interface{}) (interface{}, interface{}, interface{}, error) {
	m := "GetUpdateState"
	r.Order = append(r.Order, m)

	if r.returnErrorFor(m) {
		return nil, nil, nil, r.Error
	}

	return nil, nil, nil, nil
}

func (r *testResource) Name() string {
	return "testResource"
}

func (r *testResource) ProcessCreateState(ctx context.Context, obj, createState interface{}) error {
	m := "ProcessCreateState"
	r.Order = append(r.Order, m)

	if r.returnErrorFor(m) {
		return r.Error
	}

	return nil
}

func (r *testResource) ProcessDeleteState(ctx context.Context, obj, deleteState interface{}) error {
	m := "ProcessDeleteState"
	r.Order = append(r.Order, m)

	if r.returnErrorFor(m) {
		return r.Error
	}

	return nil
}

func (r *testResource) ProcessUpdateState(ctx context.Context, obj, updateState interface{}) error {
	m := "ProcessUpdateState"
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
