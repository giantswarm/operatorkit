package v1alpha1

import (
	"github.com/ghodss/yaml"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	kindApp              = "App"
	appDocumentationLink = "https://pkg.go.dev/github.com/giantswarm/apiextensions/pkg/apis/application/v1alpha1?tab=doc#App"
)

const appCRDYAML = `
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: apps.application.giantswarm.io
spec:
  group: application.giantswarm.io
  scope: Namespaced
  version: v1alpha1
  names:
    kind: App
    plural: apps
    singular: app
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: |
        Defines an App resource, which represents an application to be running in a Kubernetes cluster.
        Reconciled by app-operator.
      properties:
        spec:
          type: object
          properties:
            catalog:
              description: |
                Name of the AppCatalog to install this app from. Find more information in the AppCatalog
                CRD documentation.
              type: string
            name:
              description: |
                Name of this App.
              type: string
            namespace:
              description: |
                Kubernetes namespace in which to install the workloads defined by this App.
              type: string
            version:
              description: Version of the app to be deployed.
              type: string
            config:
              description: |
                Configuration details for the app.
              type: object
              properties:
                configMap:
                  description: |
                    If present, points to a ConfigMap resource that holds configuration data
                    used by the app.
                  type: object
                  properties:
                    name:
                      description: |
                        Name of the ConfigMap.
                      type: string
                    namespace:
                      description: |
                        Namespace to find the ConfigMap in.
                      type: string
                  required: ["name", "namespace"]
                secret:
                  description: |
                    If present, points to a Secret resoure that can be used by the app.
                  type: object
                  properties:
                    name:
                      description: |
                        Name of the Secret.
                      type: string
                    namespace:
                      description: |
                        Namespace to find the Secret in.
                      type: string
                  required: ["name", "namespace"]
            kubeConfig:
              description: |
                The kubeconfig to use to connect to the tenant cluster when deploying the app.
              type: object
              properties:
                inCluster:
                  description: |
                    Defines whether to use inCluster credentials. If true, the context and secret
                    properties must not be set.
                  type: boolean
                context:
                  description: |
                    Kubeconfig context part to use when not using inCluster credentials.
                  type: object
                  properties:
                    name:
                      description: |
                        Context name.
                      type: string
                secret:
                  description: |
                    References a Secret resource holding the kubeconfig details, if not using inCluster credentials.
                  type: object
                  properties:
                    name:
                      description: |
                        Name of the Secret resource.
                      type: string
                    namespace:
                      description: |
                        Namespace holding the Secret resource.
                      type: string
                  required: ["name", "namespace"]
            userConfig:
              description: |
                Additional and optional user-provided configuration for the app.
              type: object
              properties:
                configMap:
                  description: |
                    Reference to an optional ConfigMap.
                  type: object
                  properties:
                    name:
                      description: |
                        Name of the ConfigMap resource.
                      type: string
                    namespace:
                      description: |
                        Namespace holding the ConfigMap resource.
                      type: string
                  required: ["name", "namespace"]
                secret:
                  description: |
                    Reference to an optional Secret resource.
                  type: object
                  properties:
                    name:
                      description: |
                        Name of the Secret resource.
                      type: string
                    namespace:
                      description: |
                        Namespace holding the Secret resource.
                      type: string
                  required: ["name", "namespace"]
          required: ["catalog", "name", "namespace", "version"]
`

var appCRD *apiextensionsv1beta1.CustomResourceDefinition

func init() {
	err := yaml.Unmarshal([]byte(appCRDYAML), &appCRD)
	if err != nil {
		panic(err)
	}
}

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
	return appCRD.DeepCopy()
}

func NewAppTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: SchemeGroupVersion.String(),
		Kind:       kindApp,
	}
}

// NewAppCR returns an App Custom Resource.
func NewAppCR() *App {
	return &App{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				crDocsAnnotation: appDocumentationLink,
			},
		},
		TypeMeta: NewAppTypeMeta(),
	}
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// App CRs might look something like the following.
//
//     apiVersion: application.giantswarm.io/v1alpha1
//     kind: App
//     metadata:
//       name: "prometheus"
//       labels:
//         app-operator.giantswarm.io/version: "1.0.0"
//
//     spec:
//       catalog: "giantswarm"
//       name: "prometheus"
//       namespace: "monitoring"
//       version: "1.0.0"
//       config:
//         configMap:
//           name: "prometheus-values"
//           namespace: "monitoring"
//         secret:
//           name: "prometheus-secrets"
//           namespace: "monitoring"
//       kubeConfig:
//         inCluster: false
//         context:
//           name: "giantswarm-12345"
//         secret:
//           name: "giantswarm-12345"
//           namespace: "giantswarm"
//         userConfig:
//           configMap:
//             name: "prometheus-user-values"
//             namespace: "monitoring"
//
//     status:
//       appVersion: "2.4.3" # Optional value from Chart.yaml with the version of the deployed app.
//       release:
//         lastDeployed: "2018-11-30T21:06:20Z"
//         status: "DEPLOYED"
//       version: "1.1.0" # Required value from Chart.yaml with the version of the chart.
//
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
	// UserConfig is the user config to be applied when the app is deployed.
	UserConfig AppSpecUserConfig `json:"userConfig" yaml:"userConfig"`
	// Version is the version of the app that should be deployed.
	// e.g. 1.0.0
	Version string `json:"version" yaml:"version"`
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
	// InCluster is a flag for whether to use InCluster credentials. When true the
	// context name and secret should not be set.
	InCluster bool `json:"inCluster" yaml:"inCluster"`
	// Context is the kubeconfig context.
	Context AppSpecKubeConfigContext `json:"context" yaml:"context"`
	// Secret references a secret containing the kubconfig.
	Secret AppSpecKubeConfigSecret `json:"secret" yaml:"secret"`
}

type AppSpecKubeConfigContext struct {
	// Name is the name of the kubeconfig context.
	// e.g. giantswarm-12345.
	Name string `json:"name" yaml:"name"`
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
	// Namespace is the namespace of the user values config map on the control plane,
	// e.g. 123ab.
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
	// AppVersion is the value of the AppVersion field in the Chart.yaml of the
	// deployed app. This is an optional field with the version of the
	// component being deployed.
	// e.g. 0.21.0.
	// https://docs.helm.sh/developing_charts/#the-chart-yaml-file
	AppVersion string `json:"appVersion" yaml:"appVersion"`
	// Release is the status of the Helm release for the deployed app.
	Release AppStatusRelease `json:"release" yaml:"release"`
	// Version is the value of the Version field in the Chart.yaml of the
	// deployed app.
	// e.g. 1.0.0.
	Version string `json:"version" yaml:"version"`
}

type AppStatusRelease struct {
	// LastDeployed is the time when the app was last deployed.
	LastDeployed DeepCopyTime `json:"lastDeployed" yaml:"lastDeployed"`
	// Reason is the description of the last status of helm release when the app is
	// not installed successfully, e.g. deploy resource already exists.
	Reason string `json:"reason,omitempty" yaml:"reason,omitempty"`
	// Status is the status of the deployed app,
	// e.g. DEPLOYED.
	Status string `json:"status" yaml:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []App `json:"items"`
}
