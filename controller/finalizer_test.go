package controller

import (
	"reflect"
	"testing"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_createAddFinalizerPatch(t *testing.T) {
	testCases := []struct {
		name                         string
		object                       *apiv1.Pod
		operatorName                 string
		expectedCancelReconciliation bool
		expectedPatch                []patchSpec
		errorMatcher                 func(error) bool
	}{
		{
			name: "case 0: No finalizer is set yet",
			object: &apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "TestPod",
					Namespace:       "TestNamespace",
					ResourceVersion: "123",
					SelfLink:        "/some/path",
				},
			},
			operatorName:                 "test-operator",
			expectedCancelReconciliation: true,
			expectedPatch: []patchSpec{
				{
					Op:    "add",
					Path:  "/metadata/finalizers",
					Value: []string{},
				},
				{
					Op:    "add",
					Path:  "/metadata/finalizers/-",
					Value: "operatorkit.giantswarm.io/test-operator",
				},
				{
					Op:    "test",
					Path:  "/metadata/resourceVersion",
					Value: "123",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 1: Finalizer is already set",
			object: &apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "TestPod",
					Namespace: "TestNamespace",
					SelfLink:  "/some/path",
					Finalizers: []string{
						"operatorkit.giantswarm.io/test-operator",
					},
				},
			},
			operatorName:                 "test-operator",
			expectedCancelReconciliation: false,
			expectedPatch:                nil,
			errorMatcher:                 nil,
		},
		{
			name: "case 2: DeletionTimestamp is already set",
			object: &apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					DeletionTimestamp: getTime(),
					Name:              "TestPod",
					Namespace:         "TestNamespace",
					SelfLink:          "/some/path",
				},
			},
			operatorName:                 "test-operator",
			expectedCancelReconciliation: true,
			expectedPatch:                nil,
			errorMatcher:                 nil,
		},
		{
			name: "case 3: Other finalizers are already set",
			object: &apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "TestPod",
					Namespace:       "TestNamespace",
					ResourceVersion: "123",
					SelfLink:        "/some/path",
					Finalizers: []string{
						"operatorkit.giantswarm.io/other-operator",
						"operatorkit.giantswarm.io/123-operator",
					},
				},
			},
			operatorName:                 "test-operator",
			expectedCancelReconciliation: true,
			expectedPatch: []patchSpec{
				{
					Op:    "add",
					Path:  "/metadata/finalizers/-",
					Value: "operatorkit.giantswarm.io/test-operator",
				},
				{
					Op:    "test",
					Path:  "/metadata/resourceVersion",
					Value: "123",
				},
			},
			errorMatcher: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			patch, cancelReconciliation, err := createAddFinalizerPatch(tc.object, tc.operatorName)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}
			if !reflect.DeepEqual(patch, tc.expectedPatch) {
				t.Fatalf("patch == %v, want %v", patch, tc.expectedPatch)
			}
			if cancelReconciliation != tc.expectedCancelReconciliation {
				t.Fatalf("cancelReconciliation == %v, want %v", cancelReconciliation, tc.expectedCancelReconciliation)
			}
		})
	}
}

func Test_createRemoveFinalizerPatch(t *testing.T) {
	testCases := []struct {
		name          string
		object        *apiv1.Pod
		operatorName  string
		expectedPatch []patchSpec
		expectedPath  string
		errorMatcher  func(error) bool
	}{
		{
			name: "case 0: No finalizer is set",
			object: &apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "TestPod",
					Namespace: "TestNamespace",
					SelfLink:  "/some/path",
				},
			},
			operatorName:  "test-operator",
			expectedPatch: nil,
			errorMatcher:  nil,
		},
		{
			name: "case 1: Finalizer is set correctly",
			object: &apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					DeletionTimestamp: getTime(),
					Finalizers: []string{
						"operatorkit.giantswarm.io/test-operator",
					},
					Name:      "TestPod",
					Namespace: "TestNamespace",
					SelfLink:  "/some/path",
				},
			},
			operatorName: "test-operator",
			expectedPatch: []patchSpec{
				{
					Op:    "replace",
					Path:  "/metadata/finalizers",
					Value: []string{},
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 2: Mulitple finalizers are set",
			object: &apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					DeletionTimestamp: getTime(),
					Finalizers: []string{
						"operatorkit.giantswarm.io/123-operator",
						"operatorkit.giantswarm.io/test-operator",
						"operatorkit.giantswarm.io/other-operator",
					},
					Name:      "TestPod",
					Namespace: "TestNamespace",
					SelfLink:  "/some/path",
				},
			},
			operatorName: "test-operator",
			expectedPatch: []patchSpec{
				{
					Op:   "replace",
					Path: "/metadata/finalizers",
					Value: []string{
						"operatorkit.giantswarm.io/123-operator",
						"operatorkit.giantswarm.io/other-operator",
					},
				},
			},
			errorMatcher: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			patch, err := createRemoveFinalizerPatch(tc.object, tc.operatorName)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}
			if !reflect.DeepEqual(patch, tc.expectedPatch) {
				t.Fatalf("patch == %v, want %v", patch, tc.expectedPatch)
			}
		})
	}
}

func getTime() *metav1.Time {
	time := metav1.Now()
	return &time
}
