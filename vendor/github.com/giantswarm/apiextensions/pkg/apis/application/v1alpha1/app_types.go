package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	kindApp = "App"
)

// NewAppCRD returns a new custom resource definition for App.
// This might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: apps.application.giantswarm.io
//     spec:
//       group: application.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: App
//         plural: apps
//         singular: app
//
func NewAppCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "apps.application.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "application.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "App",
				Plural:   "apps",
				Singular: "app",
			},
			Subresources: &apiextensionsv1beta1.CustomResourceSubresources{
				Status: &apiextensionsv1beta1.CustomResourceSubresourceStatus{},
			},
		},
	}
}

func NewAppTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: version,
		Kind:       kindApp,
	}
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type App struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              AppSpec   `json:"spec"`
	Status            AppStatus `json:"status" yaml:"status"`
}

type AppSpec struct {
	// Catalog is the name of the app catalog this app belongs to.
	// e.g. giantswarm
	Catalog string `json:"catalog" yaml:"catalog"`
	// Config is the config to be applied when the app is deployed.
	Config AppSpecConfig `json:"config" yaml:"config"`
	// KubeConfig is the kubeconfig to connect to the cluster when deploying
	// the app.
	KubeConfig AppSpecKubeConfig `json:"kubeConfig" yaml:"kubeConfig"`
	// Name is the name of the app to be deployed.
	// e.g. kubernetes-prometheus
	Name string `json:"name" yaml:"name"`
	// Namespace is the namespace where the app should be deployed.
	// e.g. monitoring
	Namespace string `json:"namespace" yaml:"namespace"`
	// Release is the version of the app that should be deployed.
	// e.g. 1.0.0
	Release string `json:"release" yaml:"release"`
	// UserConfig is the user config to be applied when the app is deployed.
	UserConfig AppSpecConfig `json:"userConfig" yaml:"userConfig"`
}

type AppSpecConfig struct {
	// ConfigMap references a config map containing values that should be
	// applied to the app.
	ConfigMap AppSpecConfigConfigMap `json:"configMap" yaml:"configMap"`
	// Secret references a secret containing secret values that should be
	// applied to the app.
	Secret AppSpecConfigSecret `json:"secret" yaml:"secret"`
}

type AppSpecConfigConfigMap struct {
	// Name is the name of the config map containing app values to apply,
	// e.g. prometheus-values.
	Name string `json:"name" yaml:"name"`
	// Namespace is the namespace of the values config map,
	// e.g. monitoring.
	Namespace string `json:"namespace" yaml:"namespace"`
}

type AppSpecConfigSecret struct {
	// Name is the name of the secret containing app values to apply,
	// e.g. prometheus-secret.
	Name string `json:"name" yaml:"name"`
	// Namespace is the namespace of the secret,
	// e.g. kube-system.
	Namespace string `json:"namespace" yaml:"namespace"`
}

type AppSpecKubeConfig struct {
	// Secret references a secret containing the kubconfig.
	Secret AppSpecKubeConfigSecret `json:"secret" yaml:"secret"`
}

type AppSpecKubeConfigSecret struct {
	// Name is the name of the secret containing the kubeconfig,
	// e.g. app-operator-kubeconfig.
	Name string `json:"name" yaml:"name"`
	// Namespace is the namespace of the secret containing the kubeconfig,
	// e.g. giantswarm.
	Namespace string `json:"namespace" yaml:"namespace"`
}

type AppSpecUserConfig struct {
	// ConfigMap references a config map containing user values that should be
	// applied to the app.
	ConfigMap AppSpecUserConfigConfigMap `json:"configMap" yaml:"configMap"`
	// Secret references a secret containing user secret values that should be
	// applied to the app.
	Secret AppSpecUserConfigSecret `json:"secret" yaml:"secret"`
}

type AppSpecUserConfigConfigMap struct {
	// Name is the name of the config map containing user values to apply,
	// e.g. prometheus-user-values.
	Name string `json:"name" yaml:"name"`
	// Namespace is the namespace of the user values config map,
	// e.g. monitoring.
	Namespace string `json:"namespace" yaml:"namespace"`
}

type AppSpecUserConfigSecret struct {
	// Name is the name of the secret containing user values to apply,
	// e.g. prometheus-user-secret.
	Name string `json:"name" yaml:"name"`
	// Namespace is the namespace of the secret,
	// e.g. kube-system.
	Namespace string `json:"namespace" yaml:"namespace"`
}

type AppStatus struct {
	// Status is the status of the deployed app.
	Status string `json:"status" yaml:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []App `json:"items"`
}
