package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewAWSClusterConfigCRD returns a new custom resource definition for
// AWSClusterConfig. This might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: awsclusterconfigs.core.giantswarm.io
//     spec:
//       group: core.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: AWSClusterConfig
//         plural: awsclusterconfigs
//         singular: awsclusterconfig
//
func NewAWSClusterConfigCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "awsclusterconfigs.core.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "core.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "AWSClusterConfig",
				Plural:   "awsclusterconfigs",
				Singular: "awsclusterconfig",
			},
		},
	}
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AWSClusterConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              AWSClusterConfigSpec `json:"spec"`
}

type AWSClusterConfigSpec struct {
	Guest         AWSClusterConfigSpecGuest         `json:"guest" yaml:"guest"`
	VersionBundle AWSClusterConfigSpecVersionBundle `json:"versionBundle" yaml:"versionBundle"`
}

type AWSClusterConfigSpecGuest struct {
	ClusterGuestConfig `json:",inline" yaml:",inline"`
	CredentialSecret   AWSClusterConfigSpecGuestCredentialSecret `json:"credentialSecret" yaml:"credentialSecret"`
	Masters            []AWSClusterConfigSpecGuestMaster         `json:"masters,omitempty" yaml:"masters,omitempty"`
	Workers            []AWSClusterConfigSpecGuestWorker         `json:"workers,omitempty" yaml:"workers,omitempty"`
}

// AWSClusterConfigSpecGuestCredentialSecret points to the K8s Secret
// containing credentials for an AWS account in which the guest cluster should
// be created.
type AWSClusterConfigSpecGuestCredentialSecret struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

type AWSClusterConfigSpecGuestMaster struct {
	AWSClusterConfigSpecGuestNode `json:",inline" yaml:",inline"`
}

type AWSClusterConfigSpecGuestWorker struct {
	AWSClusterConfigSpecGuestNode `json:",inline" yaml:",inline"`
	Labels                        map[string]string `json:"labels" yaml:"labels"`
}
type AWSClusterConfigSpecGuestNode struct {
	ID           string `json:"id" yaml:"id"`
	InstanceType string `json:"instanceType,omitempty" yaml:"instanceType,omitempty"`
}

type AWSClusterConfigSpecVersionBundle struct {
	Version string `json:"version" yaml:"version"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AWSClusterConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AWSClusterConfig `json:"items"`
}
