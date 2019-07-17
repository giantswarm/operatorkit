package retryresource

// TODO Uncomment tests when new framework.Resource interface is created.
//
//import (
//	"context"
//	"fmt"
//	"reflect"
//	"testing"
//
//	"github.com/giantswarm/backoff"
//	"github.com/giantswarm/micrologger/microloggertest"
//
//	"github.com/giantswarm/operatorkit/controller"
//)
//
//// Test_RetryCRUDResourceOps_ProcessDelete_ResourceOrder_RetryOnError ensures the
//// resource's methods are executed as expected when retrying the deletion
//// process.
//func Test_RetryCRUDResourceOps_ProcessDelete_ResourceOrder_RetryOnError(t *testing.T) {
//	testCases := []struct {
//		ErrorCount          int
//		ErrorMethod         string
//		ExpectedMethodOrder []string
//	}{
//		{
//			ErrorCount:  1,
//			ErrorMethod: "GetCurrentState",
//			ExpectedMethodOrder: []string{
//				"GetCurrentState",
//				"GetCurrentState",
//				"GetDesiredState",
//				"NewDeletePatch",
//				"ApplyCreatePatch",
//				"ApplyDeletePatch",
//				"ApplyUpdatePatch",
//			},
//		},
//		{
//			ErrorCount:  2,
//			ErrorMethod: "GetCurrentState",
//			ExpectedMethodOrder: []string{
//				"GetCurrentState",
//				"GetCurrentState",
//				"GetCurrentState",
//				"GetDesiredState",
//				"NewDeletePatch",
//				"ApplyCreatePatch",
//				"ApplyDeletePatch",
//				"ApplyUpdatePatch",
//			},
//		},
//		{
//			ErrorCount:  2,
//			ErrorMethod: "ApplyDeletePatch",
//			ExpectedMethodOrder: []string{
//				"GetCurrentState",
//				"GetDesiredState",
//				"NewDeletePatch",
//				"ApplyCreatePatch",
//				"ApplyDeletePatch",
//				"ApplyDeletePatch",
//				"ApplyDeletePatch",
//				"ApplyUpdatePatch",
//			},
//		},
//	}
//
//	for i, tc := range testCases {
//		tr := &testCRUDResource{
//			Error:       fmt.Errorf("test error"),
//			ErrorCount:  tc.ErrorCount,
//			ErrorMethod: tc.ErrorMethod,
//		}
//		rs := []framework.Resource{
//			tr,
//		}
//		bf := func() backoff.Interface {
//			return &backoff.ZeroBackOff{}
//		}
//
//		c := WrapConfig{
//			Logger:         microloggertest.New(),
//			BackOffFactory: bf,
//		}
//		wrapped, err := Wrap(rs, c)
//		if err != nil {
//			t.Fatal("test", i, "expected", nil, "got", err)
//		}
//
//		err = framework.ProcessDelete(context.TODO(), nil, wrapped)
//		if err != nil {
//			t.Fatal("test", i, "expected", nil, "got", err)
//		}
//
//		if !reflect.DeepEqual(tc.ExpectedMethodOrder, tr.Order) {
//			t.Fatal("test", i, "expected", tc.ExpectedMethodOrder, "got", tr.Order)
//		}
//	}
//}
//
//// Test_RetryCRUDResourceOps_ProcessDelete_ResourceOrder ensures the resource's methods
//// are executed as expected when deleting resources using the wrapping retry
//// resource.
//func Test_RetryCRUDResourceOps_ProcessDelete_ResourceOrder(t *testing.T) {
//	tr := &testCRUDResource{}
//	rs := []framework.Resource{
//		tr,
//	}
//	bf := func() backoff.Interface {
//		return &backoff.ZeroBackOff{}
//	}
//
//	c := WrapConfig{
//		Logger:         microloggertest.New(),
//		BackOffFactory: bf,
//	}
//	wrapped, err := Wrap(rs, c)
//	if err != nil {
//		t.Fatal("expected", nil, "got", err)
//	}
//
//	err = framework.ProcessDelete(context.TODO(), nil, wrapped)
//	if err != nil {
//		t.Fatal("expected", nil, "got", err)
//	}
//
//	e := []string{
//		"GetCurrentState",
//		"GetDesiredState",
//		"NewDeletePatch",
//		"ApplyCreatePatch",
//		"ApplyDeletePatch",
//		"ApplyUpdatePatch",
//	}
//	if !reflect.DeepEqual(e, tr.Order) {
//		t.Fatal("expected", e, "got", tr.Order)
//	}
//}
//
//// Test_RetryCRUDResourceOps_ProcessUpdate_ResourceOrder_RetryOnError ensures the
//// resource's methods are executed as expected when retrying the update
//// process.
//func Test_RetryCRUDResourceOps_ProcessUpdate_ResourceOrder_RetryOnError(t *testing.T) {
//	testCases := []struct {
//		ErrorCount          int
//		ErrorMethod         string
//		ExpectedMethodOrder []string
//	}{
//		{
//			ErrorCount:  1,
//			ErrorMethod: "GetCurrentState",
//			ExpectedMethodOrder: []string{
//				"GetCurrentState",
//				"GetCurrentState",
//				"GetDesiredState",
//				"NewUpdatePatch",
//				"ApplyCreatePatch",
//				"ApplyDeletePatch",
//				"ApplyUpdatePatch",
//			},
//		},
//		{
//			ErrorCount:  2,
//			ErrorMethod: "GetCurrentState",
//			ExpectedMethodOrder: []string{
//				"GetCurrentState",
//				"GetCurrentState",
//				"GetCurrentState",
//				"GetDesiredState",
//				"NewUpdatePatch",
//				"ApplyCreatePatch",
//				"ApplyDeletePatch",
//				"ApplyUpdatePatch",
//			},
//		},
//		{
//			ErrorCount:  2,
//			ErrorMethod: "ApplyUpdatePatch",
//			ExpectedMethodOrder: []string{
//				"GetCurrentState",
//				"GetDesiredState",
//				"NewUpdatePatch",
//				"ApplyCreatePatch",
//				"ApplyDeletePatch",
//				"ApplyUpdatePatch",
//				"ApplyUpdatePatch",
//				"ApplyUpdatePatch",
//			},
//		},
//	}
//
//	for i, tc := range testCases {
//		tr := &testCRUDResource{
//			Error:       fmt.Errorf("test error"),
//			ErrorCount:  tc.ErrorCount,
//			ErrorMethod: tc.ErrorMethod,
//		}
//		rs := []framework.Resource{
//			tr,
//		}
//		bf := func() backoff.Interface {
//			return &backoff.ZeroBackOff{}
//		}
//
//		c := WrapConfig{
//			Logger:         microloggertest.New(),
//			BackOffFactory: bf,
//		}
//		wrapped, err := Wrap(rs, c)
//		if err != nil {
//			t.Fatal("test", i, "expected", nil, "got", err)
//		}
//
//		err = framework.ProcessUpdate(context.TODO(), nil, wrapped)
//		if err != nil {
//			t.Fatal("test", i, "expected", nil, "got", err)
//		}
//
//		if !reflect.DeepEqual(tc.ExpectedMethodOrder, tr.Order) {
//			t.Fatal("test", i, "expected", tc.ExpectedMethodOrder, "got", tr.Order)
//		}
//	}
//}
//
//// Test_RetryCRUDResourceOps_ProcessUpdate_ResourceOrder ensures the resource's methods
//// are executed as expected when updating resources using the wrapping retry
//// resource.
//func Test_RetryCRUDResourceOps_ProcessUpdate_ResourceOrder(t *testing.T) {
//	tr := &testCRUDResource{}
//	rs := []framework.Resource{
//		tr,
//	}
//	bf := func() backoff.Interface {
//		return &backoff.ZeroBackOff{}
//	}
//
//	c := WrapConfig{
//		Logger:         microloggertest.New(),
//		BackOffFactory: bf,
//	}
//	wrapped, err := Wrap(rs, c)
//	if err != nil {
//		t.Fatal("expected", nil, "got", err)
//	}
//
//	err = framework.ProcessUpdate(context.TODO(), nil, wrapped)
//	if err != nil {
//		t.Fatal("expected", nil, "got", err)
//	}
//
//	e := []string{
//		"GetCurrentState",
//		"GetDesiredState",
//		"NewUpdatePatch",
//		"ApplyCreatePatch",
//		"ApplyDeletePatch",
//		"ApplyUpdatePatch",
//	}
//	if !reflect.DeepEqual(e, tr.Order) {
//		t.Fatal("expected", e, "got", tr.Order)
//	}
//}
//
//type testCRUDResource struct {
//	Error       error
//	ErrorCount  int
//	ErrorMethod string
//	Order       []string
//
//	errorCount int
//}
//
//func (r *testCRUDResource) Name() string {
//	return "testCRUDResource"
//}
//
//func (r *testCRUDResource) GetCurrentState(ctx context.Context, obj interface{}) (interface{}, error) {
//	m := "GetCurrentState"
//	r.Order = append(r.Order, m)
//
//	if r.returnErrorFor(m) {
//		return nil, r.Error
//	}
//
//	return nil, nil
//}
//
//func (r *testCRUDResource) GetDesiredState(ctx context.Context, obj interface{}) (interface{}, error) {
//	m := "GetDesiredState"
//	r.Order = append(r.Order, m)
//
//	if r.returnErrorFor(m) {
//		return nil, r.Error
//	}
//
//	return nil, nil
//}
//
//func (r *testCRUDResource) NewUpdatePatch(ctx context.Context, obj, cur, des interface{}) (*framework.Patch, error) {
//	m := "NewUpdatePatch"
//	r.Order = append(r.Order, m)
//
//	p := framework.NewPatch()
//	p.SetCreateChange("test create data")
//	p.SetUpdateChange("test update data")
//	p.SetDeleteChange("test delete data")
//	return p, nil
//}
//
//func (r *testCRUDResource) NewDeletePatch(ctx context.Context, obj, cur, des interface{}) (*framework.Patch, error) {
//	m := "NewDeletePatch"
//	r.Order = append(r.Order, m)
//
//	p := framework.NewPatch()
//	p.SetCreateChange("test create data")
//	p.SetUpdateChange("test update data")
//	p.SetDeleteChange("test delete data")
//	return p, nil
//}
//
//func (r *testCRUDResource) ApplyCreateChange(ctx context.Context, obj, createState interface{}) error {
//	m := "ApplyCreatePatch"
//	r.Order = append(r.Order, m)
//
//	if r.returnErrorFor(m) {
//		return r.Error
//	}
//
//	return nil
//}
//
//func (r *testCRUDResource) ApplyDeleteChange(ctx context.Context, obj, deleteState interface{}) error {
//	m := "ApplyDeletePatch"
//	r.Order = append(r.Order, m)
//
//	if r.returnErrorFor(m) {
//		return r.Error
//	}
//
//	return nil
//}
//
//func (r *testCRUDResource) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
//	m := "ApplyUpdatePatch"
//	r.Order = append(r.Order, m)
//
//	if r.returnErrorFor(m) {
//		return r.Error
//	}
//
//	return nil
//}
//
//func (r *testCRUDResource) returnErrorFor(errorMethod string) bool {
//	ok := r.Error != nil && r.ErrorCount > r.errorCount && r.ErrorMethod == errorMethod
//
//	if ok {
//		r.errorCount++
//		return true
//	}
//
//	return false
//}
