package framework

import (
	"context"
	"testing"

	"github.com/giantswarm/micrologger/loggermeta"
	"github.com/giantswarm/micrologger/microloggertest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_ResourceSet_InitCtx(t *testing.T) {
	testCases := []struct {
		Object                   interface{}
		ExpectedLoggerMetaSet    bool
		ExpectedLoggerMetaObject string
	}{
		{
			Object: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					SelfLink: "/api/v1/namespace/default/pods/test-pod",
				},
			},
			ExpectedLoggerMetaSet:    true,
			ExpectedLoggerMetaObject: "/api/v1/namespace/default/pods/test-pod",
		},
		{
			Object:                nil,
			ExpectedLoggerMetaSet: false,
		},
		{
			Object:                "non-runtime-object",
			ExpectedLoggerMetaSet: false,
		},
	}
	for i, tc := range testCases {
		var err error

		var r *ResourceSet
		{
			c := ResourceSetConfig{}

			c.Handles = func(obj interface{}) bool { return false }
			c.Logger = microloggertest.New()
			c.Resources = []Resource{
				&testResource{},
			}

			r, err = NewResourceSet(c)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
		}

		ctx, err := r.InitCtx(context.Background(), tc.Object)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		meta, ok := loggermeta.FromContext(ctx)
		if tc.ExpectedLoggerMetaSet != ok {
			t.Fatal("test", i+1, "expected", tc.ExpectedLoggerMetaSet, "got", ok)
		}

		if ok && (tc.ExpectedLoggerMetaObject != meta.KeyVals["object"]) {
			t.Fatal("test", i+1, "expected", tc.ExpectedLoggerMetaObject, "got", meta.KeyVals["object"])
		}
	}
}
