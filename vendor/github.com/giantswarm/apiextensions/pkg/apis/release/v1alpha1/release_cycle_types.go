package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	kindReleaseCycle = "ReleaseCycle"

	CyclePhaseUpcoming ReleaseCyclePhase = "upcoming"
	CyclePhaseEnabled  ReleaseCyclePhase = "enabled"
	CyclePhaseDisabled ReleaseCyclePhase = "disabled"
	CyclePhaseEOL      ReleaseCyclePhase = "eol"
)

type ReleaseCyclePhase string

func (r ReleaseCyclePhase) String() string {
	return string(r)
}

var releaseCycleValidation = &apiextensionsv1beta1.CustomResourceValidation{
	// See http://json-schema.org/learn.
	OpenAPIV3Schema: &apiextensionsv1beta1.JSONSchemaProps{
		Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
			"spec": {
				Type: "object",
				Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
					"disabledDate": {
						Type:   "string",
						Format: "date",
					},
					"enabledDate": {
						Type:   "string",
						Format: "date",
					},
					"phase": {
						Enum: []apiextensionsv1beta1.JSON{
							{Raw: []byte(`"upcoming"`)},
							{Raw: []byte(`"enabled"`)},
							{Raw: []byte(`"disabled"`)},
							{Raw: []byte(`"eol"`)},
						},
					},
				},
				Required: []string{
					"phase",
				},
			},
			"status": {
				Type: "object",
			},
		},
	},
}

// NewReleaseCycleCRD looks like following.
//
//	kind: CustomResourceDefinition
//	apiVersion: apiextensions.k8s.io/v1beta1
//	metadata:
//	  name: releasecycles.release.giantswarm.io
//	spec:
//	  group: release.giantswarm.io
//	  version: v1alpha1
//	  names:
//	    plural: releasecycles
//	    singular: releasecycle
//	    kind: ReleaseCycle
//	  scope: Cluster
//	  subresources:
//	    status: {}
//	  validation:
//	    openAPIV3Schema:
//	      properties:
//	        spec:
//	          type: object
//	          required:
//	            - phase
//	          properties:
//	            disabledDate:
//	              type: string
//	              format: date
//	            enabledDate:
//	              type: string
//	              format: date
//	            phase:
//	              enum:
//	                - upcoming
//	                - enabled
//	                - disabled
//	                - eol
//	        status:
//	          type: object
//
func NewReleaseCycleCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "releasecycles.release.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "release.giantswarm.io",
			Scope:   "Cluster",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "ReleaseCycle",
				Plural:   "releasecycles",
				Singular: "releasecycle",
			},
			Subresources: &apiextensionsv1beta1.CustomResourceSubresources{
				Status: &apiextensionsv1beta1.CustomResourceSubresourceStatus{},
			},
			Validation: releaseCycleValidation,
		},
	}
}

func NewReleaseCycleTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: version,
		Kind:       kindReleaseCycle,
	}
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ReleaseCycle CRs might look something like the following.
//
//	apiVersion: "release.giantswarm.io/v1alpha1"
//	kind: "ReleaseCycle"
//	metadata:
//	  name: "aws.v6.1.0"
//	  labels:
//	    giantswarm.io/managed-by: "opsctl"
//	    giantswarm.io/provider: "aws"
//	spec:
//	  disabledDate: 2019-01-12
//	  enabledDate: 2019-01-08
//	  phase: "enabled"
//
type ReleaseCycle struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              ReleaseCycleSpec   `json:"spec" yaml:"spec"`
	Status            ReleaseCycleStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type ReleaseCycleSpec struct {
	// DisabledDate is the date of the cycle phase being changed to "disabled".
	DisabledDate DeepCopyDate `json:"disabledDate,omitempty" yaml:"disabledDate,omitempty"`
	// EnabledDate is the date of the cycle phase being changed to "enabled".
	EnabledDate DeepCopyDate `json:"enabledDate,omitempty" yaml:"enabledDate,omitempty"`
	// Phase is the release phase. It can be one of: "upcoming", "enabled",
	// "disabled", "eol".
	Phase ReleaseCyclePhase `json:"phase" yaml:"phase"`
}

type ReleaseCycleStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ReleaseCycleList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []ReleaseCycle `json:"items" yaml:"items"`
}
