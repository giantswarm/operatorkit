package v1alpha1

import (
	"time"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DrainerConfigStatusStatusTrue = "True"
)

const (
	DrainerConfigStatusTypeDrained = "Drained"
)

const (
	DrainerConfigStatusTypeTimeout = "Timeout"
)

// NewDrainerConfigCRD returns a new custom resource definition for
// DrainerConfig. This might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: drainerconfigs.core.giantswarm.io
//     spec:
//       group: core.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: DrainerConfig
//         plural: drainerconfigs
//         singular: drainerconfig
//
func NewDrainerConfigCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "drainerconfigs.core.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "core.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "DrainerConfig",
				Plural:   "drainerconfigs",
				Singular: "drainerconfig",
			},
		},
	}
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DrainerConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              DrainerConfigSpec   `json:"spec"`
	Status            DrainerConfigStatus `json:"status"`
}

type DrainerConfigSpec struct {
	Guest         DrainerConfigSpecGuest         `json:"guest" yaml:"guest"`
	VersionBundle DrainerConfigSpecVersionBundle `json:"versionBundle" yaml:"versionBundle"`
}

type DrainerConfigSpecGuest struct {
	Cluster DrainerConfigSpecGuestCluster `json:"cluster" yaml:"cluster"`
	Node    DrainerConfigSpecGuestNode    `json:"node" yaml:"node"`
}

type DrainerConfigSpecGuestCluster struct {
	API DrainerConfigSpecGuestClusterAPI `json:"api" yaml:"api"`
	// ID is the guest cluster ID of which a node should be drained.
	ID string `json:"id" yaml:"id"`
}

type DrainerConfigSpecGuestClusterAPI struct {
	// Endpoint is the guest cluster API endpoint.
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type DrainerConfigSpecGuestNode struct {
	// Name is the identifier of the guest cluster's master and worker nodes. In
	// Kubernetes/Kubectl they are represented as node names. The names are manage
	// in an abstracted way because of provider specific differences.
	//
	//     AWS: EC2 instance DNS.
	//     Azure: VM name.
	//     KVM: host cluster pod name.
	//
	Name string `json:"name" yaml:"name"`
}

type DrainerConfigSpecVersionBundle struct {
	Version string `json:"version" yaml:"version"`
}

type DrainerConfigStatus struct {
	Conditions []DrainerConfigStatusCondition `json:"conditions" yaml:"conditions"`
}

// DrainerConfigStatusCondition expresses a condition in which a node may is.
type DrainerConfigStatusCondition struct {
	// Date is the time of the given condition being updated.
	Date time.Time `json:"date" yaml:"date"`
	// Status may be True, False or Unknown.
	Status string `json:"status" yaml:"status"`
	// Type may be Pending, Ready, Draining, Drained.
	Type string `json:"type" yaml:"type"`
}

// DeepCopyInto implements the deep copy magic the k8s codegen is not able to
// generate out of the box.
func (in *DrainerConfigStatusCondition) DeepCopyInto(out *DrainerConfigStatusCondition) {
	*out = *in
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DrainerConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []DrainerConfig `json:"items"`
}
