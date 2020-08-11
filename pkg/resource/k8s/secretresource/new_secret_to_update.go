package secretresource

import (
	"reflect"

	corev1 "k8s.io/api/core/v1"
)

// newSecretToUpdate creates a new instance of Secret ready to be used as an
// argument to Update method of generated client. It returns nil if the name or
// namespace doesn't match or if objects don't have differences in scope of
// interest.
func newSecretToUpdate(current, desired *corev1.Secret) *corev1.Secret {
	if current.Namespace != desired.Namespace {
		return nil
	}
	if current.Name != desired.Name {
		return nil
	}

	merged := current.DeepCopy()

	merged.Annotations = desired.Annotations
	merged.Labels = desired.Labels

	merged.Data = desired.Data
	merged.StringData = desired.StringData

	if reflect.DeepEqual(current, merged) {
		return nil
	}

	return merged
}
