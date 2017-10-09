package fake

import (
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// List represents a list of custom objects.
type List struct {
	apismetav1.TypeMeta `json:",inline"`
	apismetav1.ListMeta `json:"metadata,omitempty"`

	Items []*CustomObject `json:"items" yaml:"items"`
}
