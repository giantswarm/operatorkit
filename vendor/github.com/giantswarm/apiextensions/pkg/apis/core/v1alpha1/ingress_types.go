package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewIngressConfigCRD returns a new custom resource definition for
// IngressConfig. This might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: ingressconfigs.core.giantswarm.io
//     spec:
//       group: core.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: IngressConfig
//         plural: ingressconfigs
//         singular: ingressconfig
//
func NewIngressConfigCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "ingressconfigs.core.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "core.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "IngressConfig",
				Plural:   "ingressconfigs",
				Singular: "ingressconfig",
			},
		},
	}
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type IngressConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              IngressConfigSpec `json:"spec"`
}

type IngressConfigSpec struct {
	GuestCluster  IngressConfigSpecGuestCluster   `json:"guestCluster" yaml:"guestCluster"`
	HostCluster   IngressConfigSpecHostCluster    `json:"hostCluster" yaml:"hostCluster"`
	ProtocolPorts []IngressConfigSpecProtocolPort `json:"protocolPorts" yaml:"protocolPorts"`
	VersionBundle IngressConfigSpecVersionBundle  `json:"versionBundle" yaml:"versionBundle"`
}

type IngressConfigSpecGuestCluster struct {
	ID        string `json:"id" yaml:"id"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Service   string `json:"service" yaml:"service"`
}

type IngressConfigSpecHostCluster struct {
	IngressController IngressConfigSpecHostClusterIngressController `json:"ingressController" yaml:"ingressController"`
}

type IngressConfigSpecHostClusterIngressController struct {
	ConfigMap string `json:"configMap" yaml:"configMap"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Service   string `json:"service" yaml:"service"`
}

type IngressConfigSpecProtocolPort struct {
	IngressPort int    `json:"ingressPort" yaml:"ingressPort"`
	LBPort      int    `json:"lbPort" yaml:"lbPort"`
	Protocol    string `json:"protocol" yaml:"protocol"`
}

type IngressConfigSpecVersionBundle struct {
	Version string `json:"version" yaml:"version"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type IngressConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []IngressConfig `json:"items"`
}
