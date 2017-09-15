package logresource

import (
	"bufio"
	"bytes"
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

		err = framework.ProcessCreate(nil, wrapped)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// Ensure the operations are properly executed in order.
	{
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

	// Ensure the operation names are properly logged in order.
	{
		fields := []string{
			"GetCurrentState",
			"GetCurrentState",
			"GetDesiredState",
			"GetDesiredState",
			"GetCreateState",
			"GetCreateState",
			"ProcessCreateState",
			"ProcessCreateState",
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

			val, ok := got["operation"]
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

		err = framework.ProcessDelete(nil, wrapped)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// Ensure the operations are properly executed in order.
	{
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

	// Ensure the operation names are properly logged in order.
	{
		fields := []string{
			"GetCurrentState",
			"GetCurrentState",
			"GetDesiredState",
			"GetDesiredState",
			"GetDeleteState",
			"GetDeleteState",
			"ProcessDeleteState",
			"ProcessDeleteState",
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

			val, ok := got["operation"]
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

		err = framework.ProcessUpdate(nil, wrapped)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	// Ensure the operations are properly executed in order.
	{
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

	// Ensure the operation names are properly logged in order.
	{
		fields := []string{
			"GetCurrentState",
			"GetCurrentState",
			"GetDesiredState",
			"GetDesiredState",
			"GetUpdateState",
			"GetUpdateState",
			"ProcessCreateState",
			"ProcessCreateState",
			"ProcessDeleteState",
			"ProcessDeleteState",
			"ProcessUpdateState",
			"ProcessUpdateState",
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

			val, ok := got["operation"]
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

func (r *testResource) GetCreateState(obj, cur, des interface{}) (interface{}, error) {
	m := "GetCreateState"
	r.Order = append(r.Order, m)

	return nil, nil
}

func (r *testResource) GetDeleteState(obj, cur, des interface{}) (interface{}, error) {
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

func (r *testResource) ProcessCreateState(obj, cre interface{}) error {
	m := "ProcessCreateState"
	r.Order = append(r.Order, m)

	return nil
}

func (r *testResource) ProcessDeleteState(obj, del interface{}) error {
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
