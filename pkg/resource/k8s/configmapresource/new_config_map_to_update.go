package configmapresource

import (
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
)

// newConfigMapToUpdate creates a new instance of ConfigMap ready to be used as an
// argument to Update method of generated client. It returns nil if the name or
// namespace doesn't match or if objects don't have differences in scope of
// interest.
func newConfigMapToUpdate(current, desired *corev1.ConfigMap, allowedLabels map[string]bool) *corev1.ConfigMap {
	if current.Namespace != desired.Namespace {
		return nil
	}
	if current.Name != desired.Name {
		return nil
	}

	merged := current.DeepCopy()

	merged.Annotations = desired.Annotations

	if allowedLabels != nil {
		for k, v := range current.Labels {
			if _, exist := desired.Labels[k]; exist {
				// if annotation is already in desired spec, skip it.
				continue
			}

			if _, ok := allowedLabels[k]; ok {
				merged.Labels[k] = v
			}
		}
	}

	merged.Labels = desired.Labels

	merged.BinaryData = desired.BinaryData
	merged.Data = desired.Data

	diff := cmp.Diff(current, merged)
	fmt.Printf("CONFIGMAP DIFF\n%s\n", diff)

	if reflect.DeepEqual(current, merged) {
		return nil
	}

	return merged
}
