package resource

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newConfigMapFromFilled(modifyFunc func(*corev1.ConfigMap)) *corev1.ConfigMap {
	c := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-name",
			Namespace: "test-namespace",
			Annotations: map[string]string{
				"test-annotation-key": "test-annotation-value",
			},
			Labels: map[string]string{
				"test-label-key": "test-label-value",
			},
		},
		BinaryData: map[string][]byte{
			"test-binarydata-key": []byte("test-bindarydata-value"),
		},
		Data: map[string]string{
			"test-data-key": "test-data-value",
		},
	}

	modifyFunc(c)
	return c
}

func Test_newConfigMapToUpdate(t *testing.T) {
	testCases := []struct {
		name             string
		current          *corev1.ConfigMap
		desired          *corev1.ConfigMap
		expectedToUpdate *corev1.ConfigMap
	}{
		{
			name: "case 0: returns updated object",
			current: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Annotations = map[string]string{
					"current-annotation-key-0": "current-annotation-value-0",
				}
				v.Labels = nil
				v.BinaryData = nil
				v.Data = map[string]string{
					"current-data-key-0": "current-data-value-0",
				}
			}),
			desired: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Annotations = nil
				v.Labels = map[string]string{
					"desired-label-key-0": "desired-label-value-0",
				}
				v.BinaryData = map[string][]byte{
					"desired-binarydata-key-0": []byte("desired-binarydata-value-0"),
				}
				v.Data = map[string]string{
					"desired-data-key-0": "desired-data-value-0",
				}
			}),
			expectedToUpdate: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Annotations = nil
				v.Labels = map[string]string{
					"desired-label-key-0": "desired-label-value-0",
				}
				v.BinaryData = map[string][]byte{
					"desired-binarydata-key-0": []byte("desired-binarydata-value-0"),
				}
				v.Data = map[string]string{
					"desired-data-key-0": "desired-data-value-0",
				}
			}),
		},
		{
			name: "case 1: returns nil for the same objects",
			current: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
			}),
			desired: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
			}),
			expectedToUpdate: nil,
		},
		{
			name: "case 2: updates only Annotations, Labels and ConfigMap unique fields",
			current: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Name = "current-name"
				v.Finalizers = []string{
					"current-finalizer-0",
				}
				v.SelfLink = "current-selflink"

				v.Annotations = map[string]string{
					"current-annotation-key-0": "current-annotation-value-0",
				}
			}),
			desired: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Name = "desired-name"
				v.Finalizers = nil
				v.SelfLink = "desired-selflink"

				v.Annotations = map[string]string{
					"desired-annotation-key-0": "desired-annotation-value-0",
				}
			}),
			expectedToUpdate: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Name = "current-name"
				v.Finalizers = []string{
					"current-finalizer-0",
				}
				v.SelfLink = "current-selflink"

				v.Annotations = map[string]string{
					"desired-annotation-key-0": "desired-annotation-value-0",
				}
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			toUpdate := newConfigMapToUpdate(tc.current, tc.desired)

			if !reflect.DeepEqual(toUpdate, tc.expectedToUpdate) {
				t.Fatalf("toUpdate == %v, want %v", toUpdate, tc.expectedToUpdate)
			}
		})
	}
}

func Test_newConfigMapToUpdate_Annotations(t *testing.T) {
	testCases := []struct {
		name             string
		current          *corev1.ConfigMap
		desired          *corev1.ConfigMap
		expectedToUpdate *corev1.ConfigMap
	}{
		{
			name: "case 0: returns object with Annotations updated",
			current: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Annotations = map[string]string{
					"current-annotation-key-0": "current-annotation-value-0",
				}
			}),
			desired: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Annotations = map[string]string{
					"desired-annotation-key-0": "desired-annotation-value-0",
					"desired-annotation-key-1": "desired-annotation-value-1",
				}
			}),
			expectedToUpdate: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Annotations = map[string]string{
					"desired-annotation-key-0": "desired-annotation-value-0",
					"desired-annotation-key-1": "desired-annotation-value-1",
				}
			}),
		},
		{
			name: "case 1: returns object with Annotations set",
			current: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Annotations = nil
			}),
			desired: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Annotations = map[string]string{
					"desired-annotation-key-0": "desired-annotation-value-0",
					"desired-annotation-key-1": "desired-annotation-value-1",
				}
			}),
			expectedToUpdate: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Annotations = map[string]string{
					"desired-annotation-key-0": "desired-annotation-value-0",
					"desired-annotation-key-1": "desired-annotation-value-1",
				}
			}),
		},
		{
			name: "case 2: returns object with Annotations unset",
			current: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Annotations = map[string]string{
					"current-annotation-key-0": "current-annotation-value-0",
				}
			}),
			desired: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Annotations = nil
			}),
			expectedToUpdate: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Annotations = nil
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			toUpdate := newConfigMapToUpdate(tc.current, tc.desired)

			if !reflect.DeepEqual(toUpdate, tc.expectedToUpdate) {
				t.Fatalf("toUpdate == %v, want %v", toUpdate, tc.expectedToUpdate)
			}
		})
	}
}

func Test_newConfigMapToUpdate_Labels(t *testing.T) {
	testCases := []struct {
		name             string
		current          *corev1.ConfigMap
		desired          *corev1.ConfigMap
		expectedToUpdate *corev1.ConfigMap
	}{
		{
			name: "case 0: returns object with Labels updated",
			current: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Labels = map[string]string{
					"current-label-key-0": "current-label-value-0",
				}
			}),
			desired: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Labels = map[string]string{
					"desired-label-key-0": "desired-label-value-0",
					"desired-label-key-1": "desired-label-value-1",
				}
			}),
			expectedToUpdate: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Labels = map[string]string{
					"desired-label-key-0": "desired-label-value-0",
					"desired-label-key-1": "desired-label-value-1",
				}
			}),
		},
		{
			name: "case 1: returns object with Labels set",
			current: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Labels = nil
			}),
			desired: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Labels = map[string]string{
					"desired-label-key-0": "desired-label-value-0",
					"desired-label-key-1": "desired-label-value-1",
				}
			}),
			expectedToUpdate: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Labels = map[string]string{
					"desired-label-key-0": "desired-label-value-0",
					"desired-label-key-1": "desired-label-value-1",
				}
			}),
		},
		{
			name: "case 2: returns object with Labels unset",
			current: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Labels = map[string]string{
					"current-label-key-0": "current-label-value-0",
				}
			}),
			desired: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Labels = nil
			}),
			expectedToUpdate: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Labels = nil
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			toUpdate := newConfigMapToUpdate(tc.current, tc.desired)

			if !reflect.DeepEqual(toUpdate, tc.expectedToUpdate) {
				t.Fatalf("toUpdate == %v, want %v", toUpdate, tc.expectedToUpdate)
			}
		})
	}
}

func Test_newConfigMapToUpdate_BinaryData(t *testing.T) {
	testCases := []struct {
		name             string
		current          *corev1.ConfigMap
		desired          *corev1.ConfigMap
		expectedToUpdate *corev1.ConfigMap
	}{
		{
			name: "case 0: returns object with BinaryData updated",
			current: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.BinaryData = map[string][]byte{
					"current-binarydata-key-0": []byte("current-binarydata-value-0"),
				}
			}),
			desired: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.BinaryData = map[string][]byte{
					"desired-binarydata-key-0": []byte("desired-binarydata-value-0"),
					"desired-binarydata-key-1": []byte("desired-binarydata-value-1"),
				}
			}),
			expectedToUpdate: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.BinaryData = map[string][]byte{
					"desired-binarydata-key-0": []byte("desired-binarydata-value-0"),
					"desired-binarydata-key-1": []byte("desired-binarydata-value-1"),
				}
			}),
		},
		{
			name: "case 1: returns object with BinaryData set",
			current: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.BinaryData = nil
			}),
			desired: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.BinaryData = map[string][]byte{
					"desired-binarydata-key-0": []byte("desired-binarydata-value-0"),
					"desired-binarydata-key-1": []byte("desired-binarydata-value-1"),
				}
			}),
			expectedToUpdate: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.BinaryData = map[string][]byte{
					"desired-binarydata-key-0": []byte("desired-binarydata-value-0"),
					"desired-binarydata-key-1": []byte("desired-binarydata-value-1"),
				}
			}),
		},
		{
			name: "case 2: returns object with BinaryData unset",
			current: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.BinaryData = map[string][]byte{
					"current-binarydata-key-0": []byte("current-binarydata-value-0"),
				}
			}),
			desired: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.BinaryData = nil
			}),
			expectedToUpdate: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.BinaryData = nil
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			toUpdate := newConfigMapToUpdate(tc.current, tc.desired)

			if !reflect.DeepEqual(toUpdate, tc.expectedToUpdate) {
				t.Fatalf("toUpdate == %v, want %v", toUpdate, tc.expectedToUpdate)
			}
		})
	}
}

func Test_newConfigMapToUpdate_Data(t *testing.T) {
	testCases := []struct {
		name             string
		current          *corev1.ConfigMap
		desired          *corev1.ConfigMap
		expectedToUpdate *corev1.ConfigMap
	}{
		{
			name: "case 0: returns object with Data updated",
			current: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Data = map[string]string{
					"current-data-key-0": "current-data-value-0",
				}
			}),
			desired: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Data = map[string]string{
					"desired-data-key-0": "desired-data-value-0",
					"desired-data-key-1": "desired-data-value-1",
				}
			}),
			expectedToUpdate: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Data = map[string]string{
					"desired-data-key-0": "desired-data-value-0",
					"desired-data-key-1": "desired-data-value-1",
				}
			}),
		},
		{
			name: "case 1: returns object with Data set",
			current: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Data = nil
			}),
			desired: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Data = map[string]string{
					"desired-data-key-0": "desired-data-value-0",
					"desired-data-key-1": "desired-data-value-1",
				}
			}),
			expectedToUpdate: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Data = map[string]string{
					"desired-data-key-0": "desired-data-value-0",
					"desired-data-key-1": "desired-data-value-1",
				}
			}),
		},
		{
			name: "case 2: returns object with Data unset",
			current: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Data = map[string]string{
					"current-data-key-0": "current-data-value-0",
				}
			}),
			desired: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Data = nil
			}),
			expectedToUpdate: newConfigMapFromFilled(func(v *corev1.ConfigMap) {
				v.Data = nil
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			toUpdate := newConfigMapToUpdate(tc.current, tc.desired)

			if !reflect.DeepEqual(toUpdate, tc.expectedToUpdate) {
				t.Fatalf("toUpdate == %v, want %v", toUpdate, tc.expectedToUpdate)
			}
		})
	}
}
