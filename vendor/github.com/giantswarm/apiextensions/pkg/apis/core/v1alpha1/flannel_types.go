package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewFlannelConfigCRD returns a new custom resource definition for
// FlannelConfig. This might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: flannelconfigs.core.giantswarm.io
//     spec:
//       group: core.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: FlannelConfig
//         plural: flannelconfigs
//         singular: flannelconfig
//
func NewFlannelConfigCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "flannelconfigs.core.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "core.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "FlannelConfig",
				Plural:   "flannelconfigs",
				Singular: "flannelconfig",
			},
		},
	}
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FlannelConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              FlannelConfigSpec `json:"spec"`
}

type FlannelConfigSpec struct {
	Bridge        FlannelConfigSpecBridge        `json:"bridge" yaml:"bridge"`
	Cluster       FlannelConfigSpecCluster       `json:"cluster" yaml:"cluster"`
	Flannel       FlannelConfigSpecFlannel       `json:"flannel" yaml:"flannel"`
	Health        FlannelConfigSpecHealth        `json:"health" yaml:"health"`
	VersionBundle FlannelConfigSpecVersionBundle `json:"versionBundle" yaml:"versionBundle"`
}

type FlannelConfigSpecBridge struct {
	Docker FlannelConfigSpecBridgeDocker `json:"docker" yaml:"docker"`
	Spec   FlannelConfigSpecBridgeSpec   `json:"spec" yaml:"spec"`
}

type FlannelConfigSpecBridgeDocker struct {
	Image string `json:"image" yaml:"image"`
}

type FlannelConfigSpecBridgeSpec struct {
	Interface      string                         `json:"interface" yaml:"interface"`
	PrivateNetwork string                         `json:"privateNetwork" yaml:"privateNetwork"`
	DNS            FlannelConfigSpecBridgeSpecDNS `json:"dns" yaml:"dns"`
	NTP            FlannelConfigSpecBridgeSpecNTP `json:"ntp" yaml:"ntp"`
}

type FlannelConfigSpecBridgeSpecDNS struct {
	Servers []string `json:"servers" yaml:"servers"`
}

type FlannelConfigSpecBridgeSpecNTP struct {
	Servers []string `json:"servers" yaml:"servers"`
}

type FlannelConfigSpecCluster struct {
	ID        string `json:"id" yaml:"id"`
	Customer  string `json:"customer" yaml:"customer"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

type FlannelConfigSpecFlannel struct {
	Spec FlannelConfigSpecFlannelSpec `json:"spec" yaml:"spec"`
}

type FlannelConfigSpecFlannelSpec struct {
	Network   string `json:"network" yaml:"network"`
	SubnetLen int    `json:"subnetLen" yaml:"subnetLen"`
	RunDir    string `json:"runDir" yaml:"runDir"`
	VNI       int    `json:"vni" yaml:"vni"`
}

type FlannelConfigSpecHealth struct {
	Docker FlannelConfigSpecHealthDocker `json:"docker" yaml:"docker"`
}

type FlannelConfigSpecHealthDocker struct {
	Image string `json:"image" yaml:"image"`
}

type FlannelConfigSpecVersionBundle struct {
	Version string `json:"version" yaml:"version"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FlannelConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []FlannelConfig `json:"items"`
}
