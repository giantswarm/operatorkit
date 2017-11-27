package retryresource

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/cenkalti/backoff"
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
				"NewUpdatePatch",
				"ApplyCreatePatch",
				"ApplyDeletePatch",
				"ApplyUpdatePatch",
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
				"NewUpdatePatch",
				"ApplyCreatePatch",
				"ApplyDeletePatch",
				"ApplyUpdatePatch",
			},
		},
		{
			ErrorCount:  2,
			ErrorMethod: "ApplyCreatePatch",
			ExpectedMethodOrder: []string{
				"GetCurrentState",
				"GetDesiredState",
				"NewUpdatePatch",
				"ApplyCreatePatch",
				"ApplyCreatePatch",
				"ApplyCreatePatch",
				"ApplyDeletePatch",
				"ApplyUpdatePatch",
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

		err = framework.ProcessCreate(context.TODO(), nil, wrapped)
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

	err = framework.ProcessCreate(context.TODO(), nil, wrapped)
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
				"NewDeletePatch",
				"ApplyCreatePatch",
				"ApplyDeletePatch",
				"ApplyUpdatePatch",
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
				"NewDeletePatch",
				"ApplyCreatePatch",
				"ApplyDeletePatch",
				"ApplyUpdatePatch",
			},
		},
		{
			ErrorCount:  2,
			ErrorMethod: "ApplyDeletePatch",
			ExpectedMethodOrder: []string{
				"GetCurrentState",
				"GetDesiredState",
				"NewDeletePatch",
				"ApplyCreatePatch",
				"ApplyDeletePatch",
				"ApplyDeletePatch",
				"ApplyDeletePatch",
				"ApplyUpdatePatch",
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

		err = framework.ProcessDelete(context.TODO(), nil, wrapped)
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
				"NewUpdatePatch",
				"ApplyCreatePatch",
				"ApplyDeletePatch",
				"ApplyUpdatePatch",
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
				"NewUpdatePatch",
				"ApplyCreatePatch",
				"ApplyDeletePatch",
				"ApplyUpdatePatch",
			},
		},
		{
			ErrorCount:  2,
			ErrorMethod: "ApplyUpdatePatch",
			ExpectedMethodOrder: []string{
				"GetCurrentState",
				"GetDesiredState",
				"NewUpdatePatch",
				"ApplyCreatePatch",
				"ApplyDeletePatch",
				"ApplyUpdatePatch",
				"ApplyUpdatePatch",
				"ApplyUpdatePatch",
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

		err = framework.ProcessUpdate(context.TODO(), nil, wrapped)
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

	if r.returnErrorFor(m) {
		return r.Error
	}

	return nil
}

func (r *testResource) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
	m := "ApplyDeletePatch"
	r.Order = append(r.Order, m)

	if r.returnErrorFor(m) {
		return r.Error
	}

	return nil
}

func (r *testResource) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	m := "ApplyUpdatePatch"
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
