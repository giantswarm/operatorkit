package certificatetpr

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// List represents a list of CustomObject resources.
type List struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`

	Items []CustomObject `json:"items" yaml:"items"`
}
