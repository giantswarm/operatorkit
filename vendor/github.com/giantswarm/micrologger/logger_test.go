package micrologger

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/giantswarm/micrologger/loggercontext"
)

func Test_Logger_LogWithCtx(t *testing.T) {
	var err error

	out := new(bytes.Buffer)

	var log Logger
	{
		c := DefaultConfig()
		{
			c.IOWriter = out
		}
		log, err = New(c)
		if err != nil {
			t.Fatalf("setting up logger: %#v", err)
		}
	}

	{
		log.LogWithCtx(context.TODO(), "foo", "bar")

		var got map[string]string
		err := json.Unmarshal(out.Bytes(), &got)
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}

		v1, ok := got["foo"]
		if !ok {
			t.Fatalf("expected %s got %s", "foo key", "nothing")
		}
		if v1 != "bar" {
			t.Fatalf("expected %s got %s", "bar", v1)
		}
		v2, ok := got["baz"]
		if ok {
			t.Fatalf("expected %s got %s", "nothing", "baz key")
		}
		if v2 == "zap" {
			t.Fatalf("expected %s got %s", "nothing", v2)
		}
	}

	var ctx context.Context
	{
		container := loggercontext.NewContainer()
		container.KeyVals["baz"] = "zap"

		ctx = loggercontext.NewContext(context.Background(), container)
	}

	{
		out.Reset()
		log.LogWithCtx(ctx, "foo", "bar")

		var got map[string]string
		err := json.Unmarshal(out.Bytes(), &got)
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}

		v1, ok := got["foo"]
		if !ok {
			t.Fatalf("expected %s got %s", "foo key", "nothing")
		}
		if v1 != "bar" {
			t.Fatalf("expected %s got %s", "bar", v1)
		}
		v2, ok := got["baz"]
		if !ok {
			t.Fatalf("expected %s got %s", "baz key", "nothing")
		}
		if v2 != "zap" {
			t.Fatalf("expected %s got %s", "zap", v2)
		}
	}
}

func Test_Logger_With(t *testing.T) {
	var err error

	out := new(bytes.Buffer)

	var log Logger
	{
		c := DefaultConfig()
		{
			c.IOWriter = out
		}
		log, err = New(c)
		if err != nil {
			t.Fatalf("setting up logger: %#v", err)
		}
	}

	var (
		field       = "ctxField"
		wfieldValue = "test ctx field value"

		parentLog = log
		childLog  = log.With(field, wfieldValue)
	)

	// Make sure caller (old field) and added contextual field are logged.
	{
		wfieldValue := "test ctx field value"

		out.Reset()
		childLog.Log("msg", "whats up bro?")

		var got map[string]string
		json.Unmarshal(out.Bytes(), &got)

		// NOTE this tests a line number which may change if lines are modified in
		// this file.
		wcaller := "github.com/giantswarm/micrologger/logger_test.go:119"
		caller, ok := got["caller"]
		if !ok {
			t.Errorf("expected caller key")
		}
		if caller != wcaller {
			t.Errorf("want caller %s, got %s", wcaller, caller)
		}

		fieldValue, ok := got[field]
		if !ok {
			t.Errorf("want set field %s", field)
		}
		if fieldValue != wfieldValue {
			t.Errorf("want fieldValue %s, got %s", wfieldValue, fieldValue)
		}
	}

	// Make sure parent logger remained unchanged.
	{
		out.Reset()
		parentLog.Log("msg", "how are you?")

		var got map[string]string
		json.Unmarshal(out.Bytes(), &got)

		// NOTE this tests a line number which may change if lines are modified in
		// this file.
		wcaller := "github.com/giantswarm/micrologger/logger_test.go:147"
		caller, ok := got["caller"]
		if !ok {
			t.Errorf("expected caller key")
		}
		if caller != wcaller {
			t.Errorf("want caller %s, got %s", wcaller, caller)
		}

		fieldValue, ok := got[field]
		if ok {
			t.Errorf("want unset field %s, got %s", field, fieldValue)
		}
	}
}
