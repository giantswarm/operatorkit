package controller

import (
	"context"
	"reflect"
	"testing"
)

func Test_setLoggerCtxValue_doesnt_leak(t *testing.T) {
	ctx := context.Background()

	l := valueCtxLen(ctx)
	expected := 0
	if l != expected {
		t.Fatalf("countValueContextLength(ctx) - expected %d, got %d", expected, l)
	}

	ctx = setLoggerCtxValue(ctx, "foo", "bar")
	ctx = setLoggerCtxValue(ctx, "bar", "baz")
	ctx = setLoggerCtxValue(ctx, "baz", "foo")

	l = valueCtxLen(ctx)
	expected = 1
	if l != expected {
		t.Fatalf("countValueContextLength(ctx) - expected %d, got %d", expected, l)
	}
}

func valueCtxLen(ctx context.Context) int {
	return countValueContextLength(0, ctx)
}

func countValueContextLength(i int, ctx interface{}) int {
	if !isValueCtx(ctx) {
		return i
	}

	v := reflect.ValueOf(ctx)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	v = v.FieldByName("Context")
	i++

	if v.IsValid() {
		return countValueContextLength(i, v.Interface())
	}

	return i
}

func isValueCtx(ctx interface{}) bool {
	t := reflect.TypeOf(ctx)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.PkgPath() != "context" || t.Name() != "valueCtx" {
		return false
	}

	return true
}
