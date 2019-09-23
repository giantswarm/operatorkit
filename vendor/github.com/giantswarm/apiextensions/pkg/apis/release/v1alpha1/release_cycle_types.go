package v1alpha1

import (
	"github.com/ghodss/yaml"
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

const releaseCycleCRDYAML = `
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: releasecycles.release.giantswarm.io
spec:
  group: release.giantswarm.io
  scope: Cluster
  version: v1alpha1
  names:
    kind: ReleaseCycle
    plural: releasecycles
    singular: releasecycle
  subresources:
    status: {}
  validation:
    # Note: When changing ReleaseCycle CRD schema please make sure to update
    # ReleaseCycle details in Release CRD schema accordingly.
    openAPIV3Schema:
      properties:
        spec:
          type: object
          properties:
            disabledDate:
              type: string
              format: date
            enabledDate:
              type: string
              format: date
            phase:
              enum:
              - upcoming
              - enabled
              - disabled
              - eol
          required:
          - phase
`

type ReleaseCyclePhase string

func (r ReleaseCyclePhase) String() string {
	return string(r)
}

var releaseCycleCRD *apiextensionsv1beta1.CustomResourceDefinition

func init() {
	err := yaml.Unmarshal([]byte(releaseCycleCRDYAML), &releaseCycleCRD)
	if err != nil {
		panic(err)
	}
}

func NewReleaseCycleCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return releaseCycleCRD.DeepCopy()
}

func NewReleaseCycleTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: SchemeGroupVersion.String(),
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
