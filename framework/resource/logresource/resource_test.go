package logresource

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/framework"
)

// Test_LogResource_ProcessCreate_ResourceOrder ensures the resource's
// methods are executed as expected when creating resources using the wrapping
// prometheus resource.
func Test_LogResource_ProcessCreate_ResourceOrder(t *testing.T) {
	// Setup and execute the test functionality.
	tr := &testResource{}
	out := new(bytes.Buffer)
	{
		rs := []framework.Resource{
			tr,
		}

		loggerConfig := micrologger.DefaultConfig()
		loggerConfig.IOWriter = out
		logger, err := micrologger.New(loggerConfig)
		if err != nil {
			t.Fatalf("setting up logger: %#v", err)
		}

		config := DefaultWrapConfig()
		config.Logger = logger
		wrapped, err := Wrap(rs, config)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		err = framework.ProcessCreate(context.TODO(), nil, wrapped)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// Ensure the operations are properly executed in order.
	{
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

	// Ensure the operation names are properly logged in order.
	{
		fields := []string{
			"GetCurrentState",
			"GetCurrentState",
			"GetDesiredState",
			"GetDesiredState",
			"NewUpdatePatch",
			"NewUpdatePatch",
			"ApplyCreatePatch",
			"ApplyCreatePatch",
		}
		scanner := bufio.NewScanner(out)
		for _, f := range fields {
			scanner.Scan()
			b := scanner.Bytes()

			var got map[string]string
			err := json.Unmarshal(b, &got)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}

			val, ok := got["function"]
			if !ok {
				t.Fatal("expected", true, "got", false)
			}
			if val != f {
				t.Fatal("expected", f, "got", val)
			}
		}
	}
}

// Test_LogResource_ProcessDelete_ResourceOrder ensures the resource's
// methods are executed as expected when deleting resources using the wrapping
// prometheus resource.
func Test_LogResource_ProcessDelete_ResourceOrder(t *testing.T) {
	// Setup and execute the test functionality.
	tr := &testResource{}
	out := new(bytes.Buffer)
	{
		rs := []framework.Resource{
			tr,
		}

		loggerConfig := micrologger.DefaultConfig()
		loggerConfig.IOWriter = out
		logger, err := micrologger.New(loggerConfig)
		if err != nil {
			t.Fatalf("setting up logger: %#v", err)
		}

		config := DefaultWrapConfig()
		config.Logger = logger
		wrapped, err := Wrap(rs, config)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		err = framework.ProcessDelete(context.TODO(), nil, wrapped)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// Ensure the operations are properly executed in order.
	{
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

	// Ensure the operation names are properly logged in order.
	{
		fields := []string{
			"GetCurrentState",
			"GetCurrentState",
			"GetDesiredState",
			"GetDesiredState",
			"NewDeletePatch",
			"NewDeletePatch",
			"ApplyCreatePatch",
			"ApplyCreatePatch",
			"ApplyDeletePatch",
			"ApplyDeletePatch",
			"ApplyUpdatePatch",
			"ApplyUpdatePatch",
		}
		scanner := bufio.NewScanner(out)
		for _, f := range fields {
			scanner.Scan()
			b := scanner.Bytes()

			var got map[string]string
			err := json.Unmarshal(b, &got)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}

			val, ok := got["function"]
			if !ok {
				t.Fatal("expected", true, "got", false)
			}
			if val != f {
				t.Fatal("expected", f, "got", val)
			}
		}
	}
}

// Test_LogResource_ProcessUpdate_ResourceOrder ensures the resource's methods
// are executed as expected when deleting resources using the wrapping
// prometheus resource.
func Test_LogResource_ProcessUpdate_ResourceOrder(t *testing.T) {
	// Setup and execute the test functionality.
	tr := &testResource{}
	out := new(bytes.Buffer)
	{
		rs := []framework.Resource{
			tr,
		}

		loggerConfig := micrologger.DefaultConfig()
		loggerConfig.IOWriter = out
		logger, err := micrologger.New(loggerConfig)
		if err != nil {
			t.Fatalf("setting up logger: %#v", err)
		}

		config := DefaultWrapConfig()
		config.Logger = logger
		wrapped, err := Wrap(rs, config)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}

		err = framework.ProcessUpdate(context.TODO(), nil, wrapped)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// Ensure the operations are properly executed in order.
	{
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

	// Ensure the operation names are properly logged in order.
	{
		fields := []string{
			"GetCurrentState",
			"GetCurrentState",
			"GetDesiredState",
			"GetDesiredState",
			"NewUpdatePatch",
			"NewUpdatePatch",
			"ApplyCreatePatch",
			"ApplyCreatePatch",
			"ApplyDeletePatch",
			"ApplyDeletePatch",
			"ApplyUpdatePatch",
			"ApplyUpdatePatch",
		}
		scanner := bufio.NewScanner(out)
		for _, f := range fields {
			scanner.Scan()
			b := scanner.Bytes()

			var got map[string]string
			err := json.Unmarshal(b, &got)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}

			val, ok := got["function"]
			if !ok {
				t.Fatal("expected", true, "got", false)
			}
			if val != f {
				t.Fatal("expected", f, "got", val)
			}
		}
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

func (r *testResource) ApplyCreateChange(ctx context.Context, obj, cre interface{}) error {
	m := "ApplyCreatePatch"
	r.Order = append(r.Order, m)

	return nil
}

func (r *testResource) ApplyDeleteChange(ctx context.Context, obj, del interface{}) error {
	m := "ApplyDeletePatch"
	r.Order = append(r.Order, m)

	return nil
}

func (r *testResource) ApplyUpdateChange(ctx context.Context, obj, updateState interface{}) error {
	m := "ApplyUpdatePatch"
	r.Order = append(r.Order, m)

	return nil
}

func (r *testResource) Underlying() framework.Resource {
	return r
}
