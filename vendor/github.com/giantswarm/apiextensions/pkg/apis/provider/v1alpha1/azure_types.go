package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewAzureConfigCRD returns a new custom resource definition for AzureConfig.
// This might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: azureconfigs.provider.giantswarm.io
//     spec:
//       group: provider.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: AzureConfig
//         plural: azureconfigs
//         singular: azureconfig
//       subresources:
//         status: {}
//
func NewAzureConfigCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "azureconfigs.provider.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "provider.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "AzureConfig",
				Plural:   "azureconfigs",
				Singular: "azureconfig",
			},
			Subresources: &apiextensionsv1beta1.CustomResourceSubresources{
				Status: &apiextensionsv1beta1.CustomResourceSubresourceStatus{},
			},
		},
	}
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AzureConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              AzureConfigSpec   `json:"spec"`
	Status            AzureConfigStatus `json:"status" yaml:"status"`
}

type AzureConfigSpec struct {
	Cluster       Cluster                      `json:"cluster" yaml:"cluster"`
	Azure         AzureConfigSpecAzure         `json:"azure" yaml:"azure"`
	VersionBundle AzureConfigSpecVersionBundle `json:"versionBundle" yaml:"versionBundle"`
}

type AzureConfigSpecAzure struct {
	CredentialSecret CredentialSecret                   `json:"credentialSecret" yaml:"credentialSecret"`
	DNSZones         AzureConfigSpecAzureDNSZones       `json:"dnsZones" yaml:"dnsZones"`
	VirtualNetwork   AzureConfigSpecAzureVirtualNetwork `json:"virtualNetwork" yaml:"virtualNetwork"`

	Masters []AzureConfigSpecAzureNode `json:"masters" yaml:"masters"`
	Workers []AzureConfigSpecAzureNode `json:"workers" yaml:"workers"`
}

// AzureConfigSpecAzureDNSZones contains the DNS Zones of the cluster.
type AzureConfigSpecAzureDNSZones struct {
	// API is the DNS Zone for the Kubernetes API.
	API AzureConfigSpecAzureDNSZonesDNSZone `json:"api" yaml:"api"`
	// Etcd is the DNS Zone for the etcd cluster.
	Etcd AzureConfigSpecAzureDNSZonesDNSZone `json:"etcd" yaml:"etcd"`
	// Ingress is the DNS Zone for the Ingress resource, used for customer traffic.
	Ingress AzureConfigSpecAzureDNSZonesDNSZone `json:"ingress" yaml:"ingress"`
}

// AzureConfigSpecAzureDNSZonesDNSZone points to a DNS Zone in Azure.
type AzureConfigSpecAzureDNSZonesDNSZone struct {
	// ResourceGroup is the resource group of the zone.
	ResourceGroup string `json:"resourceGroup" yaml:"resourceGroup"`
	// Name is the name of the zone.
	Name string `json:"name" yaml:"name"`
}

type AzureConfigSpecAzureVirtualNetwork struct {
	// CIDR is the CIDR for the Virtual Network.
	CIDR string `json:"cidr" yaml:"cidr"`

	// TODO: remove Master, Worker and Calico subnet cidr after azure-operator v2
	// is deleted. MasterSubnetCIDR is the CIDR for the master subnet.
	//
	//     https://github.com/giantswarm/giantswarm/issues/4358
	//
	MasterSubnetCIDR string `json:"masterSubnetCIDR" yaml:"masterSubnetCIDR"`
	// WorkerSubnetCIDR is the CIDR for the worker subnet.
	WorkerSubnetCIDR string `json:"workerSubnetCIDR" yaml:"workerSubnetCIDR"`

	// CalicoSubnetCIDR is the CIDR for the calico subnet. It has to be
	// also a worker subnet (Azure limitation).
	CalicoSubnetCIDR string `json:"calicoSubnetCIDR" yaml:"calicoSubnetCIDR"`
}

type AzureConfigSpecAzureNode struct {
	// VMSize is the master vm size (e.g. Standard_A1)
	VMSize string `json:"vmSize" yaml:"vmSize"`
	// Size of a volume mounted to /var/lib/docker.
	DockerVolumeSizeGB int `json:"dockerVolumeSizeGB" yaml:"dockerVolumeSizeGB"`
}

type AzureConfigSpecVersionBundle struct {
	Version string `json:"version" yaml:"version"`
}

type AzureConfigStatus struct {
	Cluster StatusCluster `json:"cluster" yaml:"cluster"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AzureConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AzureConfig `json:"items"`
}
