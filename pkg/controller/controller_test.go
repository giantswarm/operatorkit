package controller

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	"github.com/giantswarm/k8sclient/v5/pkg/k8sclienttest"
	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/giantswarm/operatorkit/v2/pkg/resource"
)

func Test_Controller_Collector_Register(t *testing.T) {
	prometheus.MustRegister(mustNewTestController("c-1").collector)
	prometheus.MustRegister(mustNewTestController("c-2").collector)
}

func Test_Controller_Collector_Register_Error(t *testing.T) {
	prometheus.MustRegister(mustNewTestController("same").collector)

	err := prometheus.Register(mustNewTestController("same").collector)
	_, ok := err.(prometheus.AlreadyRegisteredError)
	if !ok {
		panic("registering the same controller collector twice must not be possible")
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

func Test_hasPauseAnnotation(t *testing.T) {
	testCases := []struct {
		annotations    map[string]string
		expectedResult bool
	}{
		{
			annotations:    nil,
			expectedResult: false,
		},
		{
			annotations:    map[string]string{"other": ""},
			expectedResult: false,
		},
		{
			annotations:    map[string]string{"": "true"},
			expectedResult: false,
		},
		{
			annotations:    map[string]string{"": ""},
			expectedResult: false,
		},
		{
			annotations:    map[string]string{"operatorkit.giantswarm.io/paused": ""},
			expectedResult: false,
		},
		{
			annotations:    map[string]string{"cluster.x-k8s.io/paused": ""},
			expectedResult: false,
		},
		{
			annotations:    map[string]string{"operatorkit.giantswarm.io/paused": "false"},
			expectedResult: false,
		},
		{
			annotations:    map[string]string{"cluster.x-k8s.io/paused": "false"},
			expectedResult: false,
		},
		{
			annotations:    map[string]string{"operatorkit.giantswarm.io/paused": "true"},
			expectedResult: true,
		},
		{
			annotations:    map[string]string{"cluster.x-k8s.io/paused": "true"},
			expectedResult: true,
		},
	}
	controller := mustNewTestController("test")

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.annotations, "should be", tc.expectedResult)
			actualResult := controller.hasPauseAnnotation(tc.annotations)
			if actualResult != tc.expectedResult {
				t.Errorf("expected hasPauseAnnotation to return %#v, got %#v", tc.expectedResult, actualResult)
			}
		})
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

func mustNewTestController(n string) *Controller {
	var err error

	var controller *Controller
	{
		c := Config{
			K8sClient: k8sclienttest.NewEmpty(),
			Logger:    microloggertest.New(),
			NewRuntimeObjectFunc: func() runtime.Object {
				return new(corev1.Service)
			},
			Resources: []resource.Interface{
				&testResource{},
			},

			Name: n,
		}

		controller, err = New(c)
		if err != nil {
			panic(err)
		}
	}

	return controller
}

type testResource struct {
}

func (r *testResource) Name() string {
	return ""
}

func (r *testResource) EnsureCreated(ctx context.Context, obj interface{}) error {
	return nil
}

func (r *testResource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	return nil
}
