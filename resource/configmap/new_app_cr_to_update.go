package app

import (
	"reflect"

	"github.com/giantswarm/apiextensions/pkg/apis/application/v1alpha1"
)

// newAppCRToUpdate creates a new instance of App CR ready to be used as an
// argument to Update method of generated client. It returns nil if the name or
// namespace doesn't match or if objects don't have differences in scope of
// interest.
func newAppCRToUpdate(current, desired *v1alpha1.App) *v1alpha1.App {
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
