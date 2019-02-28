package configmap

import (
	"reflect"

	"k8s.io/api/core/v1"
)

// newConfigMapToUpdate creates a new instance of ConfigMap ready to be used as an
// argument to Update method of generated client. It returns nil if the name or
// namespace doesn't match or if objects don't have differences in scope of
// interest.
func newConfigMapToUpdate(current, desired *v1.ConfigMap) *v1.ConfigMap {
	if current.Namespace != desired.Namespace {
		return nil
	}
	if current.Name != desired.Name {
		return nil
	}

	merged := current.DeepCopy()

	merged.Annotations = desired.Annotations
	merged.Labels = desired.Labels

	merged.Spec = desired.Spec

	if reflect.DeepEqual(current, merged) {
		return nil
	}

	return merged
}
