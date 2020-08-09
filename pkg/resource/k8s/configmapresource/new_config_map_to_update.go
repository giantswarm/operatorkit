package configmapresource

import (
	"reflect"

	corev1 "k8s.io/api/core/v1"
)

// newConfigMapToUpdate creates a new instance of ConfigMap ready to be used as an
// argument to Update method of generated client. It returns nil if the name or
// namespace doesn't match or if objects don't have differences in scope of
// interest.
func newConfigMapToUpdate(current, desired *corev1.ConfigMap) *corev1.ConfigMap {
	if current.Namespace != desired.Namespace {
		return nil
	}
	if current.Name != desired.Name {
		return nil
	}

	merged := current.DeepCopy()

	merged.Annotations = desired.Annotations
	merged.Labels = desired.Labels

	merged.BinaryData = desired.BinaryData
	merged.Data = desired.Data

	if reflect.DeepEqual(current, merged) {
		return nil
	}

	return merged
}
