package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewCertConfigCRD returns a new custom resource definition for CertConfig.
// This might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: certconfigs.core.giantswarm.io
//     spec:
//       group: core.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: CertConfig
//         plural: certconfigs
//         singular: certconfig
//
func NewCertConfigCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "certconfigs.core.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "core.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "CertConfig",
				Plural:   "certconfigs",
				Singular: "certconfig",
			},
		},
	}
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CertConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              CertConfigSpec `json:"spec"`
}

type CertConfigSpec struct {
	Cert          CertConfigSpecCert          `json:"cert" yaml:"cert"`
	VersionBundle CertConfigSpecVersionBundle `json:"versionBundle" yaml:"versionBundle"`
}

type CertConfigSpecCert struct {
	AllowBareDomains    bool     `json:"allowBareDomains" yaml:"allowBareDomains"`
	AltNames            []string `json:"altNames" yaml:"altNames"`
	ClusterComponent    string   `json:"clusterComponent" yaml:"clusterComponent"`
	ClusterID           string   `json:"clusterID" yaml:"clusterID"`
	CommonName          string   `json:"commonName" yaml:"commonName"`
	DisableRegeneration bool     `json:"disableRegeneration" yaml:"disableRegeneration"`
	IPSANs              []string `json:"ipSans" yaml:"ipSans"`
	Organizations       []string `json:"organizations" yaml:"organizations"`
	TTL                 string   `json:"ttl" yaml:"ttl"`
}

type CertConfigSpecVersionBundle struct {
	Version string `json:"version" yaml:"version"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CertConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []CertConfig `json:"items"`
}
