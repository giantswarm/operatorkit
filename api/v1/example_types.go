package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Example is a basic type used for integration tests.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Example struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ExampleSpec `json:"spec"`
	// +optional
	Status ExampleStatus `json:"status"`
}

type ExampleSpec struct {
	Field1 string `json:"field1"`
}

type ExampleStatus struct {
	Conditions []ExampleCondition `json:"conditions"`
}

type ConditionStatus string

const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

// ExampleConditionType is a valid value for ExampleCondition.Type
type ExampleConditionType string

const (
	Normal ExampleConditionType = "Normal"
)

type ExampleCondition struct {
	// type is the type of the condition. Types include Normal.
	Type ExampleConditionType `json:"type"`
	// status is the status of the condition.
	// Can be True, False, Unknown.
	Status ConditionStatus `json:"status"`
	// lastTransitionTime last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
type ExampleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Example `json:"items"`
}
