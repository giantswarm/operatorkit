package framework

import (
	"reflect"
	"testing"

	"github.com/cenk/backoff"
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

// Test_Framework_ProcessCreate_ResourceOrder_Retry ensures the resource's
// methods are executed as expected when retrying the creation process.
func Test_Framework_ProcessCreate_ResourceOrder_Retry(t *testing.T) {
	testCases := []struct {
		ErrorCount          int
		ErrorMethod         string
		ExpectedMethodOrder []string
	}{
		{
			ErrorCount:  1,
			ErrorMethod: "GetCurrentState",
			ExpectedMethodOrder: []string{
				"GetCurrentState",
				"GetCurrentState",
				"GetDesiredState",
				"GetCreateState",
				"ProcessCreateState",
			},
		},
		{
			ErrorCount:  2,
			ErrorMethod: "GetCurrentState",
			ExpectedMethodOrder: []string{
				"GetCurrentState",
				"GetCurrentState",
				"GetCurrentState",
				"GetDesiredState",
				"GetCreateState",
				"ProcessCreateState",
			},
		},
		{
			ErrorCount:  2,
			ErrorMethod: "ProcessCreateState",
			ExpectedMethodOrder: []string{
				"GetCurrentState",
				"GetDesiredState",
				"GetCreateState",
				"ProcessCreateState",
				"ProcessCreateState",
				"ProcessCreateState",
			},
		},
	}

	for i, tc := range testCases {
		tr := &testResource{
			Error:       executionFailedError,
			ErrorCount:  tc.ErrorCount,
			ErrorMethod: tc.ErrorMethod,
		}
		rs := []Resource{
			tr,
		}
		bf := func() backoff.BackOff {
			return &backoff.ZeroBackOff{}
		}

		err := testMustNewFramework(t).ProcessCreateWithBackoff(nil, rs, bf)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		if !reflect.DeepEqual(tc.ExpectedMethodOrder, tr.Order) {
			t.Fatal("test", i+1, "expected", tc.ExpectedMethodOrder, "got", tr.Order)
		}
	}
}

// Test_Framework_ProcessCreate_ResourceOrder_RetryResource ensures the
// resource's methods are executed as expected when creating resources using the
// wrapping retry resource.
func Test_Framework_ProcessCreate_ResourceOrder_RetryResource(t *testing.T) {
	tr := &testResource{}
	rs := []Resource{
		tr,
	}
	bf := func() backoff.BackOff {
		return &backoff.ZeroBackOff{}
	}

	err := testMustNewFramework(t).ProcessCreateWithBackoff(nil, rs, bf)
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

// Test_Operator_ProcessDelete_ResourceOrder_Retry ensures the resource's
// methods are executed as expected when retrying the deletion process.
func Test_Operator_ProcessDelete_ResourceOrder_Retry(t *testing.T) {
	testCases := []struct {
		ErrorCount          int
		ErrorMethod         string
		ExpectedMethodOrder []string
	}{
		{
			ErrorCount:  1,
			ErrorMethod: "GetCurrentState",
			ExpectedMethodOrder: []string{
				"GetCurrentState",
				"GetCurrentState",
				"GetDesiredState",
				"GetDeleteState",
				"ProcessDeleteState",
			},
		},
		{
			ErrorCount:  2,
			ErrorMethod: "GetCurrentState",
			ExpectedMethodOrder: []string{
				"GetCurrentState",
				"GetCurrentState",
				"GetCurrentState",
				"GetDesiredState",
				"GetDeleteState",
				"ProcessDeleteState",
			},
		},
		{
			ErrorCount:  2,
			ErrorMethod: "ProcessDeleteState",
			ExpectedMethodOrder: []string{
				"GetCurrentState",
				"GetDesiredState",
				"GetDeleteState",
				"ProcessDeleteState",
				"ProcessDeleteState",
				"ProcessDeleteState",
			},
		},
	}

	for i, tc := range testCases {
		tr := &testResource{
			Error:       executionFailedError,
			ErrorCount:  tc.ErrorCount,
			ErrorMethod: tc.ErrorMethod,
		}
		rs := []Resource{
			tr,
		}
		bf := func() backoff.BackOff {
			return &backoff.ZeroBackOff{}
		}

		err := testMustNewFramework(t).ProcessDeleteWithBackoff(nil, rs, bf)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		if !reflect.DeepEqual(tc.ExpectedMethodOrder, tr.Order) {
			t.Fatal("test", i+1, "expected", tc.ExpectedMethodOrder, "got", tr.Order)
		}
	}
}

// Test_Operator_ProcessDelete_ResourceOrder_RetryResource ensures the
// resource's methods are executed as expected when deleting resources using the
// wrapping retry resource.
func Test_Operator_ProcessDelete_ResourceOrder_RetryResource(t *testing.T) {
	tr := &testResource{}
	rs := []Resource{
		tr,
	}
	bf := func() backoff.BackOff {
		return &backoff.ZeroBackOff{}
	}

	err := testMustNewFramework(t).ProcessDeleteWithBackoff(nil, rs, bf)
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

func (r *testResource) returnErrorFor(errorMethod string) bool {
	ok := r.Error != nil && r.ErrorCount > r.errorCount && r.ErrorMethod == errorMethod

	if ok {
		r.errorCount++
		return true
	}

	return false
}
