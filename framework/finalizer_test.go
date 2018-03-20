package framework

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
		expectedPatch                *patchSpec
		expectedPath                 string
		errorMatcher                 func(error) bool
	}{
		{
			name: "case 0: No finalizer is set yet",
			object: &apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "TestPod",
					Namespace: "TestNamespace",
					SelfLink:  "/some/path",
				},
			},
			operatorName:                 "test-operator",
			expectedCancelReconciliation: true,
			expectedPatch: &patchSpec{
				Op:    "add",
				Path:  "/metadata/finalizers/-",
				Value: "operatorkit.giantswarm.io/test-operator",
			},
			expectedPath: "/some/path",
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
			expectedPath:                 "",
			errorMatcher:                 nil,
		},
		{
			name: "case 2: DeletionTimestamp is already set",
			object: &apiv1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "TestPod",
					Namespace:         "TestNamespace",
					SelfLink:          "/some/path",
					DeletionTimestamp: getTime(),
				},
			},
			operatorName:                 "test-operator",
			expectedCancelReconciliation: true,
			expectedPatch:                nil,
			expectedPath:                 "",
			errorMatcher:                 nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			patch, path, cancelReconciliation, err := createAddFinalizerPatch(tc.object, tc.operatorName)

			switch {
			case err == nil && tc.errorMatcher == nil: // correct; carry on
			case err != nil && tc.errorMatcher != nil:
				if !tc.errorMatcher(err) {
					t.Fatalf("error == %#v, want matching", err)
				}
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			}
			if !reflect.DeepEqual(patch, tc.expectedPatch) {
				t.Fatalf("patch == %v, want %v", patch, tc.expectedPatch)
			}
			if path != tc.expectedPath {
				t.Fatalf("path == %v, want %v", path, tc.expectedPath)
			}
			if cancelReconciliation != tc.expectedCancelReconciliation {
				t.Fatalf("cancelReconciliation == %v, want %v", cancelReconciliation, tc.expectedCancelReconciliation)
			}
		})
	}
}

func Test_removeFinalizer(t *testing.T) {
}

func getTime() *metav1.Time {
	time := metav1.Now()
	return &time
}
