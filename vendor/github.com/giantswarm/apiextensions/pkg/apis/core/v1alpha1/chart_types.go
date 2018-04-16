package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewChartConfigCRD returns a new custom resource definition for ChartConfig.
// This might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: chartconfigs.core.giantswarm.io
//     spec:
//       group: core.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: ChartConfig
//         plural: chartconfigs
//         singular: chartconfig
//
func NewChartConfigCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "chartconfigs.core.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "core.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "ChartConfig",
				Plural:   "chartconfigs",
				Singular: "chartconfig",
			},
		},
	}
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ChartConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ChartConfigSpec `json:"spec"`
}

type ChartConfigSpec struct {
	Chart         ChartConfigSpecChart         `json:"chart" yaml:"chart"`
	VersionBundle ChartConfigSpecVersionBundle `json:"versionBundle" yaml:"versionBundle"`
}

type ChartConfigSpecChart struct {
	// Channel is the name of the Appr channel to reconcile against.
	// e.g. 1.0-stable
	Channel string `json:"channel" yaml:"channel"`
	// Name is the fully qualified name of the Helm chart to deploy.
	// e.g. quay.io/giantswarm/chart-operator-chart
	Name string `json:"name" yaml:"name"`
	// Namespace is the namespace where the Helm chart is to be deployed.
	// e.g. giantswarm
	Namespace string `json:"namespace" yaml:"namespace"`
	// Release is the name of the Helm release when the chart is deployed.
	// e.g. chart-operator
	Release string `json:"release" yaml:"release"`
}

type ChartConfigSpecVersionBundle struct {
	Version string `json:"version" yaml:"version"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ChartConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ChartConfig `json:"items"`
}
