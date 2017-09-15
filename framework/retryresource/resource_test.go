package retryresource

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/cenk/backoff"
	"github.com/giantswarm/operatorkit/framework"
)

// Test_RetryResource_ProcessCreate_ResourceOrder_RetryOnError ensures the
// resource's methods are executed as expected when retrying the creation
// process.
func Test_RetryResource_ProcessCreate_ResourceOrder_RetryOnError(t *testing.T) {
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
			Error:       fmt.Errorf("test error"),
			ErrorCount:  tc.ErrorCount,
			ErrorMethod: tc.ErrorMethod,
		}
		rs := []framework.Resource{
			tr,
		}
		bf := func() backoff.BackOff {
			return &backoff.ZeroBackOff{}
		}

		config := DefaultWrapConfig()
		config.BackOffFactory = bf
		wrapped, err := Wrap(rs, config)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		err = framework.ProcessCreate(nil, wrapped)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		if !reflect.DeepEqual(tc.ExpectedMethodOrder, tr.Order) {
			t.Fatal("test", i+1, "expected", tc.ExpectedMethodOrder, "got", tr.Order)
		}
	}
}

// Test_RetryResource_ProcessCreate_ResourceOrder ensures the resource's methods
// are executed as expected when creating resources using the wrapping retry
// resource.
func Test_RetryResource_ProcessCreate_ResourceOrder(t *testing.T) {
	tr := &testResource{}
	rs := []framework.Resource{
		tr,
	}
	bf := func() backoff.BackOff {
		return &backoff.ZeroBackOff{}
	}

	config := DefaultWrapConfig()
	config.BackOffFactory = bf
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

// Test_RetryResource_ProcessDelete_ResourceOrder_RetryOnError ensures the
// resource's methods are executed as expected when retrying the deletion
// process.
func Test_RetryResource_ProcessDelete_ResourceOrder_RetryOnError(t *testing.T) {
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
			Error:       fmt.Errorf("test error"),
			ErrorCount:  tc.ErrorCount,
			ErrorMethod: tc.ErrorMethod,
		}
		rs := []framework.Resource{
			tr,
		}
		bf := func() backoff.BackOff {
			return &backoff.ZeroBackOff{}
		}

		config := DefaultWrapConfig()
		config.BackOffFactory = bf
		wrapped, err := Wrap(rs, config)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		err = framework.ProcessDelete(nil, wrapped)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		if !reflect.DeepEqual(tc.ExpectedMethodOrder, tr.Order) {
			t.Fatal("test", i+1, "expected", tc.ExpectedMethodOrder, "got", tr.Order)
		}
	}
}

// Test_RetryResource_ProcessDelete_ResourceOrder ensures the resource's methods
// are executed as expected when deleting resources using the wrapping retry
// resource.
func Test_RetryResource_ProcessDelete_ResourceOrder(t *testing.T) {
	tr := &testResource{}
	rs := []framework.Resource{
		tr,
	}
	bf := func() backoff.BackOff {
		return &backoff.ZeroBackOff{}
	}

	config := DefaultWrapConfig()
	config.BackOffFactory = bf
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

// Test_RetryResource_ProcessUpdate_ResourceOrder_RetryOnError ensures the
// resource's methods are executed as expected when retrying the update
// process.
func Test_RetryResource_ProcessUpdate_ResourceOrder_RetryOnError(t *testing.T) {
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
				"GetUpdateState",
				"ProcessCreateState",
				"ProcessDeleteState",
				"ProcessUpdateState",
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
				"GetUpdateState",
				"ProcessCreateState",
				"ProcessDeleteState",
				"ProcessUpdateState",
			},
		},
		{
			ErrorCount:  2,
			ErrorMethod: "ProcessUpdateState",
			ExpectedMethodOrder: []string{
				"GetCurrentState",
				"GetDesiredState",
				"GetUpdateState",
				"ProcessCreateState",
				"ProcessDeleteState",
				"ProcessUpdateState",
				"ProcessUpdateState",
				"ProcessUpdateState",
			},
		},
	}

	for i, tc := range testCases {
		tr := &testResource{
			Error:       fmt.Errorf("test error"),
			ErrorCount:  tc.ErrorCount,
			ErrorMethod: tc.ErrorMethod,
		}
		rs := []framework.Resource{
			tr,
		}
		bf := func() backoff.BackOff {
			return &backoff.ZeroBackOff{}
		}

		config := DefaultWrapConfig()
		config.BackOffFactory = bf
		wrapped, err := Wrap(rs, config)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		err = framework.ProcessUpdate(nil, wrapped)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		if !reflect.DeepEqual(tc.ExpectedMethodOrder, tr.Order) {
			t.Fatal("test", i+1, "expected", tc.ExpectedMethodOrder, "got", tr.Order)
		}
	}
}

// Test_RetryResource_ProcessUpdate_ResourceOrder ensures the resource's methods
// are executed as expected when updating resources using the wrapping retry
// resource.
func Test_RetryResource_ProcessUpdate_ResourceOrder(t *testing.T) {
	tr := &testResource{}
	rs := []framework.Resource{
		tr,
	}
	bf := func() backoff.BackOff {
		return &backoff.ZeroBackOff{}
	}

	config := DefaultWrapConfig()
	config.BackOffFactory = bf
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

func (r *testResource) GetUpdateState(obj, currentState, desiredState interface{}) (interface{}, interface{}, interface{}, error) {
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

func (r *testResource) ProcessUpdateState(obj, updateState interface{}) error {
	m := "ProcessUpdateState"
	r.Order = append(r.Order, m)

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
