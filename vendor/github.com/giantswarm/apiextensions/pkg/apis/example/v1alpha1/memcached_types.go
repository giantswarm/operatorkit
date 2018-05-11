package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewMemcachedConfigCRD returns a new custom resource definition for
// MemcachedConfig. This might look something like the following.
//
//	apiVersion: apiextensions.k8s.io/v1beta1
//	kind: CustomResourceDefinition
//	metadata:
//	  name: memcachedconfigs.provider.giantswarm.io
//	spec:
//	  group: example.giantswarm.io
//	  scope: Namespaced
//	  version: v1alpha1
//	  names:
//	    kind: MemcachedConfig
//	    plural: memcachedconfigs
//	    singular: memcachedconfig
//
func NewMemcachedConfigCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "memcachedconfigs.example.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "example.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "MemcachedConfig",
				Plural:   "memcachedconfigs",
				Singular: "memcachedconfig",
			},
		},
	}
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MemcachedConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              MemcachedConfigSpec `json:"spec"`
}

type MemcachedConfigSpec struct {
	// Replicas is the number of instances of Memcache.
	Replicas int `json:"replicas" yaml:"replicas"`
	// e.g. 3
	// Memory is how much RAM to use for item storage.
	// e.g. 4G
	Memory string `json:"memory" yaml:"memory"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MemcachedConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []MemcachedConfig `json:"items"`
}
