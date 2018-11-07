package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	kindRelease = "Release"
)

// NewReleaseCRD returns a new custom resource definition for Release. This
// might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: releases.core.giantswarm.io
//     spec:
//       group: core.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: Release
//         plural: releases
//         singular: release
//       subresources:
//         status: {}
//
func NewReleaseCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "releases.core.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "core.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "Release",
				Plural:   "releases",
				Singular: "release",
			},
			Subresources: &apiextensionsv1beta1.CustomResourceSubresources{
				Status: &apiextensionsv1beta1.CustomResourceSubresourceStatus{},
			},
		},
	}
}

func NewReleaseTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: version,
		Kind:       kindRelease,
	}
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Release represents a Giant Swarm release used to describe the managed
// versiones of tenant clusters. This might look something like the following.
//
//     typemeta:
//       apiversion: "v1alpha1"
//       kind: "Release"
//     objectmeta:
//       name: "2.0.0"
//     spec:
//       active: false
//       authorities:
//       - name: azure-operator
//         version: 2.0.0
//       - name: cert-operator
//         version: 0.1.0
//       - name: chart-operator
//         version: 0.3.0
//       - name: cluster-operator
//         version: 0.7.0
//       date: "0001-01-01T00:00:00Z"
//       provider: "azure"
//       version: "2.0.0"
//       versionBundle:
//         version: "0.1.0"
//     status:
//       conditions:
//       - status: True
//         type: ExampleCondition
//
type Release struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ReleaseSpec   `json:"spec"`
	Status            ReleaseStatus `json:"status"`
}

type ReleaseSpec struct {
	Active        bool                     `json:"active" yaml:"active"`
	Authorities   []ReleaseSpecAuthority   `json:"authorities" yaml:"authorities"`
	Date          DeepCopyTime             `json:"date" yaml:"date"`
	Provider      string                   `json:"provider" yaml:"provider"`
	Version       string                   `json:"version" yaml:"version"`
	VersionBundle ReleaseSpecVersionBundle `json:"versionBundle" yaml:"versionBundle"`
}

type ReleaseSpecAuthority struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

type ReleaseSpecVersionBundle struct {
	Version string `json:"version" yaml:"version"`
}

type ReleaseStatus struct {
	Conditions []ReleaseStatusCondition `json:"conditions" yaml:"conditions"`
}

// ReleaseStatusCondition expresses a condition in which a node may is.
type ReleaseStatusCondition struct {
	// LastHeartbeatTime is the last time we got an update on a given condition.
	LastHeartbeatTime DeepCopyTime `json:"lastHeartbeatTime" yaml:"lastHeartbeatTime"`
	// LastTransitionTime is the last time the condition transitioned from one
	// status to another.
	LastTransitionTime DeepCopyTime `json:"lastTransitionTime" yaml:"lastTransitionTime"`
	// Status may be True, False or Unknown.
	Status string `json:"status" yaml:"status"`
	// Type is not yet specified.
	Type string `json:"type" yaml:"type"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Release `json:"items"`
}
