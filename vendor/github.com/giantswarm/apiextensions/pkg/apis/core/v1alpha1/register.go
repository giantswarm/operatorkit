package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	group   = "core.giantswarm.io"
	version = "v1alpha1"
)

// knownTypes is the full list of objects to register with the scheme. It
// should contain all zero values of custom objects and custom object lists
// in the group version.
var knownTypes = []runtime.Object{
	&CertConfig{},
	&CertConfigList{},
	&ChartConfig{},
	&ChartConfigList{},
	&DrainerConfig{},
	&DrainerConfigList{},
	&AWSClusterConfig{},
	&AWSClusterConfigList{},
	&AzureClusterConfig{},
	&AzureClusterConfigList{},
	&KVMClusterConfig{},
	&KVMClusterConfigList{},
	&DraughtsmanConfig{},
	&DraughtsmanConfigList{},
	&FlannelConfig{},
	&FlannelConfigList{},
	&IngressConfig{},
	&IngressConfigList{},
	&NodeConfig{},
	&NodeConfigList{},
	&StorageConfig{},
	&StorageConfigList{},
}

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{
	Group:   group,
	Version: version,
}

var (
	schemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

	// AddToScheme is used by the generated client.
	AddToScheme = schemeBuilder.AddToScheme
)

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion, knownTypes...)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
