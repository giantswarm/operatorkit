package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewNodeConfigCRD returns a new custom resource definition for NodeConfig.
// This might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: nodeconfigs.core.giantswarm.io
//     spec:
//       group: core.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: NodeConfig
//         plural: nodeconfigs
//         singular: nodeconfig
//
func NewNodeConfigCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "nodeconfigs.core.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "core.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "NodeConfig",
				Plural:   "nodeconfigs",
				Singular: "nodeconfig",
			},
		},
	}
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NodeConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              NodeConfigSpec   `json:"spec"`
	Status            NodeConfigStatus `json:"status"`
}

type NodeConfigSpec struct {
	Guest         NodeConfigSpecGuest         `json:"guest" yaml:"guest"`
	VersionBundle NodeConfigSpecVersionBundle `json:"versionBundle" yaml:"versionBundle"`
}

type NodeConfigSpecGuest struct {
	Cluster NodeConfigSpecGuestCluster `json:"cluster" yaml:"cluster"`
	Node    NodeConfigSpecGuestNode    `json:"node" yaml:"node"`
}

type NodeConfigSpecGuestCluster struct {
	API NodeConfigSpecGuestClusterAPI `json:"api" yaml:"api"`
	// ID is the guest cluster ID of which a node should be drained.
	ID string `json:"id" yaml:"id"`
}

type NodeConfigSpecGuestClusterAPI struct {
	// Endpoint is the guest cluster API endpoint.
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type NodeConfigSpecGuestNode struct {
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

type NodeConfigSpecVersionBundle struct {
	Version string `json:"version" yaml:"version"`
}

type NodeConfigStatus struct {
	Conditions []NodeConfigStatusCondition `json:"conditions" yaml:"conditions"`
}

// NodeConfigStatusCondition expresses a condition in which a node may is.
type NodeConfigStatusCondition struct {
	// LastHeartbeatTime is the last time we got an update on a given condition.
	LastHeartbeatTime DeepCopyTime `json:"lastHeartbeatTime" yaml:"lastHeartbeatTime"`
	// LastTransitionTime is the last time the condition transitioned from one
	// status to another.
	LastTransitionTime DeepCopyTime `json:"lastTransitionTime" yaml:"lastTransitionTime"`
	// Status may be True, False or Unknown.
	Status string `json:"status" yaml:"status"`
	// Type may be Pending, Ready, Draining, Drained.
	Type string `json:"type" yaml:"type"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NodeConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []NodeConfig `json:"items"`
}
