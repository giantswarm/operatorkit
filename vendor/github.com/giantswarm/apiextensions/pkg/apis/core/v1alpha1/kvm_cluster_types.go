package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewKVMClusterConfigCRD returns a new custom resource definition for
// KVMClusterConfig. This might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: kvmclusterconfigs.core.giantswarm.io
//     spec:
//       group: core.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: KVMClusterConfig
//         plural: kvmclusterconfigs
//         singular: kvmclusterconfig
//
func NewKVMClusterConfigCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "kvmclusterconfigs.core.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "core.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "KVMClusterConfig",
				Plural:   "kvmclusterconfigs",
				Singular: "kvmclusterconfig",
			},
		},
	}
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type KVMClusterConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              KVMClusterConfigSpec `json:"spec"`
}

type KVMClusterConfigSpec struct {
	Guest         KVMClusterConfigSpecGuest         `json:"guest" yaml:"guest"`
	VersionBundle KVMClusterConfigSpecVersionBundle `json:"versionBundle" yaml:"versionBundle"`
}

type KVMClusterConfigSpecGuest struct {
	ClusterGuestConfig `json:",inline" yaml:",inline"`
	Masters            []KVMClusterConfigSpecGuestMaster `json:"masters,omitempty" yaml:"masters,omitempty"`
	Workers            []KVMClusterConfigSpecGuestWorker `json:"workers,omitempty" yaml:"workers,omitempty"`
}

type KVMClusterConfigSpecGuestMaster struct {
	KVMClusterConfigSpecGuestNode `json:",inline" yaml:",inline"`
}

type KVMClusterConfigSpecGuestWorker struct {
	KVMClusterConfigSpecGuestNode `json:",inline" yaml:",inline"`
	Labels                        map[string]string `json:"labels" yaml:"labels"`
}

type KVMClusterConfigSpecGuestNode struct {
	ID            string  `json:"id" yaml:"id"`
	CPUCores      int     `json:"cpuCores,omitempty" yaml:"cpuCores,omitempty"`
	MemorySizeGB  float64 `json:"memorySizeGB,omitempty" yaml:"memorySizeGB,omitempty"`
	StorageSizeGB float64 `json:"storageSizeGB,omitempty" yaml:"storageSizeGB,omitempty"`
}

type KVMClusterConfigSpecVersionBundle struct {
	Version string `json:"version" yaml:"version"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type KVMClusterConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []KVMClusterConfig `json:"items"`
}
