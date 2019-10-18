package controller

import (
	"context"
	"reflect"
	"testing"

	"github.com/giantswarm/micrologger/microloggertest"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func Test_hashEvent(t *testing.T) {
	testCases := []struct {
		Event        watch.Event
		Concurrency  int
		ExpectedHash int
	}{
		// Test 0 ensures the hash for a generic object.
		{
			Event: watch.Event{
				Type: watch.Added,
				Object: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foobar",
						Namespace: "blah",
					},
				},
			},
			Concurrency:  32,
			ExpectedHash: 30,
		},

		// Test 1 ensures the hash for another generic object.
		{
			Event: watch.Event{
				Type: watch.Added,
				Object: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "qweqweqwe",
						Namespace: "bloop",
					},
				},
			},
			Concurrency:  32,
			ExpectedHash: 9,
		},

		// Test 2 ensures the hash for yet another generic object.
		{
			Event: watch.Event{
				Type: watch.Added,
				Object: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "dajboubfob",
						Namespace: "foiniof",
					},
				},
			},
			Concurrency:  32,
			ExpectedHash: 4,
		},
	}

	for i, tc := range testCases {
		controller := Controller{
			logger: microloggertest.New(),

			concurrency: tc.Concurrency,
		}

		hash := controller.hashEvent(context.TODO(), tc.Event)
		if tc.ExpectedHash != hash {
			t.Fatal("test", i, "expected", tc.ExpectedHash, "got", hash)
		}
	}
}

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
