package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewAzureClusterConfigCRD returns a new custom resource definition for
// AzureClusterConfig. This might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: azureclusterconfigs.core.giantswarm.io
//     spec:
//       group: core.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: AzureClusterConfig
//         plural: azureclusterconfigs
//         singular: azureclusterconfig
//
func NewAzureClusterConfigCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "azureclusterconfigs.core.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "core.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "AzureClusterConfig",
				Plural:   "azureclusterconfigs",
				Singular: "azureclusterconfig",
			},
		},
	}
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AzureClusterConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              AzureClusterConfigSpec `json:"spec"`
}

type AzureClusterConfigSpec struct {
	Guest         AzureClusterConfigSpecGuest         `json:"guest" yaml:"guest"`
	VersionBundle AzureClusterConfigSpecVersionBundle `json:"versionBundle" yaml:"versionBundle"`
}

type AzureClusterConfigSpecGuest struct {
	ClusterGuestConfig `json:",inline" yaml:",inline"`
	CredentialSecret   AzureClusterConfigSpecGuestCredentialSecret `json:"credentialSecret" yaml:"credentialSecret"`
	Masters            []AzureClusterConfigSpecGuestMaster         `json:"masters,omitempty" yaml:"masters,omitempty"`
	Workers            []AzureClusterConfigSpecGuestWorker         `json:"workers,omitempty" yaml:"workers,omitempty"`
}

// AzureClusterConfigSpecGuestCredentialSecret points to the K8s Secret
// containing credentials for an Azure subscription in which the tenant cluster
// should be created.
type AzureClusterConfigSpecGuestCredentialSecret struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

type AzureClusterConfigSpecGuestMaster struct {
	AzureClusterConfigSpecGuestNode `json:",inline" yaml:",inline"`
}

type AzureClusterConfigSpecGuestWorker struct {
	AzureClusterConfigSpecGuestNode `json:",inline" yaml:",inline"`
	Labels                          map[string]string `json:"labels" yaml:"labels"`
}

type AzureClusterConfigSpecGuestNode struct {
	ID     string `json:"id" yaml:"id"`
	VMSize string `json:"vmSize,omitempty" yaml:"vmSize,omitempty"`
}

type AzureClusterConfigSpecVersionBundle struct {
	Version string `json:"version" yaml:"version"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AzureClusterConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AzureClusterConfig `json:"items"`
}
