package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	kindIgnition    = "Ignition"
	ignitionCRDYAML = `apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  name: ignitions.core.giantswarm.io
spec:
  additionalPrinterColumns:
    - jsonPath: .status.ready
      description: Indicates that the ignition secret has been successfully rendered
        and is ready to be used
      name: ready
      type: boolean
  group: core.giantswarm.io
  names:
    kind: Ignition
    plural: ignitions
    singular: ignition
    shortNames:
    - "ign"
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      type: object
      properties:
        spec:
          type: object
          properties:
            apiServerEncryptionKey:
              type: string
  versions:
  - name: v1alpha1
    served: true
    storage: true
`
)

var ignitionCRD *apiextensionsv1beta1.CustomResourceDefinition

func init() {
	err := yaml.UnmarshalStrict([]byte(ignitionCRDYAML), &ignitionCRD)
	if err != nil {
		panic(err)
	}
}

// NewIgnitionCRD returns a new custom resource definition for an Ignition resource.
// Ignitions contain a rendered ignition template specific to nodes or groups of nodes
// in a particular cluster.
func NewIgnitionCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return ignitionCRD.DeepCopy()
}

func NewIgnitionTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: SchemeGroupVersion.String(),
		Kind:       kindIgnition,
	}
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Ignition is a Kubernetes resource (CR) which is based on the Ignition CRD defined above.
//
// An example Ignition resource can be viewed here
// https://github.com/giantswarm/apiextensions/blob/master/docs/cr/core.giantswarm.io_v1alpha1_ignition.yaml
type Ignition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              IgnitionSpec   `json:"spec"`
	Status            IgnitionStatus `json:"status"`
}

// IgnitionSpec is the interface which defines the input parameters for
// a newly rendered g8s ignition template.
type IgnitionSpec struct {
	// APIServerEncryptionKey is used in EncryptionConfiguration to encrypt Kubernetes secrets at rest.
	APIServerEncryptionKey string `json:"apiServerEncryptionKey" yaml:"apiServerEncryptionKey"`
	// BaseDomain is the base domain for all cluster services.
	// For test installations, this may be in the form
	// <clusterId>.k8s.<installation>.<region>.<provider>.gigantic.io.
	BaseDomain string `json:"baseDomain" yaml:"baseDomain"`
	// Calico provides configuration for all calico-related services.
	Calico IgnitionSpecCalico `json:"calico" yaml:"calico"`
	// ClusterID is the name of the tenant cluster to be created.
	ClusterID string `json:"clusterID" yaml:"clusterID"`
	// DisableEncryptionAtRest will disable secret encryption at rest when set to true.
	DisableEncryptionAtRest bool `json:"disableEncryptionAtRest" yaml:"disableEncryptionAtRest"`
	// Docker provides configuration for all calico-related services.
	Docker IgnitionSpecDocker `json:"docker" yaml:"docker"`
	// Etcd provides configuration for all etcd-related services.
	Etcd IgnitionSpecEtcd `json:"etcd" yaml:"etcd"`
	// Extension can be used to extend an ignition with extra configuration provided by the provider operator.
	Extension IgnitionSpecExtension `json:"extension" yaml:"extension"`
	// Ingress provides configuration for all ingress-related services.
	Ingress IgnitionSpecIngress `json:"ingress" yaml:"ingress"`
	// IsMaster determines if the rendered ignition should contain master-specific configuration.
	IsMaster bool `json:"isMaster" yaml:"isMaster"`
	// Kubernetes provides configuration for all Kubernetes-related services.
	Kubernetes IgnitionSpecKubernetes `json:"kubernetes" yaml:"kubernetes"`
	// Defines the provider which should be rendered.
	Provider string `json:"provider" yaml:"provider"`
	// Registry provides configuration for the docker registry used for core component images.
	Registry IgnitionSpecRegistry `json:"registry" yaml:"registry"`
	// SSO provides configuration for all SSO-related services.
	SSO IgnitionSpecSSO `json:"sso" yaml:"sso"`
}

type IgnitionSpecCalico struct {
	// CIDR is the CIDR-component of the IPv4 overlay subnetwork. Combined with Subnet below.
	CIDR string `json:"cidr" yaml:"cidr"`
	// Disable can be set to true to disable Calico setup.
	Disable bool `json:"disable" yaml:"disable"`
	// MTU is the maximum size of packets sent over Calico in bytes.
	MTU string `json:"mtu" yaml:"mtu"`
	// Subnet is the IP-component of the IPv4 overlay subnetwork. Combined with CIDR above.
	Subnet string `json:"subnet" yaml:"subnet"`
}

type IgnitionSpecDocker struct {
	// Daemon provides information about the Docker daemon running on TC nodes.
	Daemon IgnitionSpecDockerDaemon `json:"daemon" yaml:"daemon"`
	// NetworkSetup provides the Docker image to be used for network environment setup.
	NetworkSetup IgnitionSpecDockerNetworkSetup `json:"networkSetup" yaml:"networkSetup"`
}

type IgnitionSpecDockerDaemon struct {
	// CIDR is the fully specified subnet used for DOCKER_OPT_BIP.
	CIDR string `json:"cidr" yaml:"cidr"`
}

type IgnitionSpecDockerNetworkSetup struct {
	// Image provides the Docker image to be used for network environment setup.
	Image string `json:"image" yaml:"image"`
}

type IgnitionSpecEtcd struct {
	// Domain is the domain of the etcd service.
	Domain string `json:"domain" yaml:"domain"`
	// Port is the port of the etcd service, usually 2379.
	Port int `json:"port" yaml:"port"`
	// Prefix is the prefix to add to all etcd keys created by Kubernetes.
	Prefix string `json:"prefix" yaml:"prefix"`
}

type IgnitionSpecExtension struct {
	// Files is an optional array of files which will be rendered and added to the final node ignition.
	Files []IgnitionSpecExtensionFile `json:"files,omitempty" yaml:"files"`
	// Files is an optional array of systemd units which will be rendered and added to the final node ignition.
	Units []IgnitionSpecExtensionUnit `json:"units,omitempty" yaml:"units"`
	// Files is an optional array of users which will be added to the final node ignition.
	Users []IgnitionSpecExtensionUser `json:"users,omitempty" yaml:"users"`
}

type IgnitionSpecExtensionFile struct {
	// Content is the string containing a file with optional go-template-style replacements.
	Content string `json:"content" yaml:"content"`
	// Metadata is the filesystem metadata of the given file.
	Metadata IgnitionSpecExtensionFileMetadata `json:"metadata" yaml:"metadata"`
}

type IgnitionSpecExtensionFileMetadata struct {
	// Compression allows a file to be passed in as a base64-encoded compressed string.
	Compression bool `json:"compression" yaml:"compression"`
	// Owner is the owner of the file.
	Owner IgnitionSpecExtensionFileMetadataOwner `json:"owner" yaml:"owner"`
	// Path is the path of the file.
	Path string `json:"path" yaml:"path"`
	// Permissions is the numeric permissions applied to the file.
	Permissions int `json:"permissions" yaml:"permissions"`
}

type IgnitionSpecExtensionFileMetadataOwner struct {
	// Group is the group which owns the file.
	Group IgnitionSpecExtensionFileMetadataOwnerGroup `json:"group" yaml:"group"`
	// User is the user which owns the file.
	User IgnitionSpecExtensionFileMetadataOwnerUser `json:"user" yaml:"user"`
}

type IgnitionSpecExtensionFileMetadataOwnerUser struct {
	// ID is the UID of the user.
	ID string `json:"id" yaml:"id"`
	// Name is the name of the user.
	Name string `json:"name" yaml:"name"`
}

type IgnitionSpecExtensionFileMetadataOwnerGroup struct {
	// ID is the GID of the group.
	ID string `json:"id" yaml:"id"`
	// Name is the name of the group.
	Name string `json:"name" yaml:"name"`
}

type IgnitionSpecExtensionUnit struct {
	// Content is the string containing a systemd unit with optional go-template-style replacements.
	Content string `json:"content" yaml:"content"`
	// Metadata is the filesystem metadata of the given file.
	Metadata IgnitionSpecExtensionUnitMetadata `json:"metadata" yaml:"metadata"`
}

type IgnitionSpecExtensionUnitMetadata struct {
	// Enabled indicates that the unit should be enabled by default.
	Enabled bool `json:"enabled" yaml:"enabled"`
	// Name is the name of the unit on the filesystem and used in systemctl commands.
	Name string `json:"name" yaml:"name"`
}

type IgnitionSpecExtensionUser struct {
	// Name is the name of the user to be added to the node via ignition.
	Name string `json:"name" yaml:"name"`
	// PublicKey is the public key of the user to be added to the node via ignition.
	PublicKey string `json:"publicKey" yaml:"publicKey"`
}

type IgnitionSpecIngress struct {
	// Disable will disable the ingress controller in the TC when true.
	Disable bool `json:"disable" yaml:"disable"`
}

type IgnitionSpecKubernetes struct {
	// API holds information about the desired TC Kubernetes API.
	API IgnitionSpecKubernetesAPI `json:"api" yaml:"api"`
	// CloudProvider is the provider upon which the cluster is running. It is passed to API server as a flag.
	CloudProvider string `json:"cloudProvider" yaml:"cloudProvider"`
	// DNS hold information about the in-cluster DNS service.
	DNS IgnitionSpecKubernetesDNS `json:"dns" yaml:"dns"`
	// Domain is the domain used for services running in the cluster. Usually this is "cluster.local".
	Domain string `json:"domain" yaml:"domain"`
	// Kubelet holds information about the kubelet running on nodes.
	Kubelet IgnitionSpecKubernetesKubelet `json:"kubelet" yaml:"kubelet"`
	// IPRange is the range of IPs used for pod networking.
	IPRange string `json:"ipRange" yaml:"ipRange"`
	// OIDC hold configuration which will be applied to the apiserver OIDC flags.
	OIDC IgnitionSpecOIDC `json:"oidc" yaml:"oidc"`
}

type IgnitionSpecKubernetesAPI struct {
	// Domain is the domain of the API server.
	Domain string `json:"domain" yaml:"domain"`
	// Secure port is the port on which the API will listen.
	SecurePort int `json:"securePort" yaml:"securePort"`
}

type IgnitionSpecKubernetesDNS struct {
	// IP is the IP of the in-cluster DNS service. Usually this is
	// the same as the API server IP with the final component replaced with .10.
	IP string `json:"ip" yaml:"ip"`
}

type IgnitionSpecKubernetesKubelet struct {
	// Domain is the domain of the network.
	Domain string `json:"domain" yaml:"domain"`
}

type IgnitionSpecRegistry struct {
	// Domain is the domain of the registry to be used for pulling core component images.
	Domain string `json:"domain" yaml:"domain"`
	// Pull progress deadline is a string representing a duration to be used as a deadline
	// for pulling images.
	PullProgressDeadline string `json:"pullProgressDeadline" yaml:"pullProgressDeadline"`
}

type IgnitionSpecSSO struct {
	// PublicKey is the public key of the SSO service.
	PublicKey string `json:"publicKey" yaml:"publicKey"`
}

type IgnitionSpecOIDC struct {
	// Enabled indicates that the OIDC settings should be applied when true.
	Enabled bool `json:"enabled" yaml:"enabled"`
	// The client ID for the OpenID Connect client, must be set if IssuerURL is set.
	ClientID string `json:"clientID" yaml:"clientId"`
	// The URL of the OpenID issuer, only HTTPS scheme will be accepted.
	// If set, it will be used to verify the OIDC JSON Web Token (JWT).
	IssuerURL string `json:"issuerUrl" yaml:"issuerUrl"`
	// The OpenID claim to use as the user name. Note that claims other
	// than the default ('sub') is not guaranteed to be unique and immutable.
	UsernameClaim string `json:"usernameClaim" yaml:"usernameClaim"`
	// If provided, all usernames will be prefixed with this value. If not provided, username
	// claims other than 'email' are prefixed by the issuer URL to avoid clashes. To skip any
	// prefixing, provide the value '-'.
	UsernamePrefix string `json:"usernamePrefix" yaml:"usernamePrefix"`
	// If provided, the name of a custom OpenID Connect claim for specifying
	// user groups. The claim value is expected to be a string or JSON encoded array of strings.
	GroupsClaim string `json:"groupsClaim" yaml:"groupsClaim"`
	// If provided, all groups will be prefixed with this value to prevent conflicts with other
	// authentication strategies.
	GroupsPrefix string `json:"groupsPrefix" yaml:"groupsPrefix"`
}

// IgnitionStatus holds the rendering result.
type IgnitionStatus struct {
	// DataSecret is a reference to the secret containing the rendered ignition once created.
	DataSecret IgnitionStatusSecret `json:"dataSecretName" yaml:"dataSecretName"`
	// FailureReason is a short string indicating the reason rendering failed (if it did).
	FailureReason string `json:"failureReason" yaml:"failureReason"`
	// FailureMessage is a longer message indicating the reason rendering failed (if it did).
	FailureMessage string `json:"failureMessage" yaml:"failureMessage"`
	// Ready will be true when the referenced secret contains the rendered ignition and can be used for creating nodes.
	Ready bool `json:"ready" yaml:"ready"`
	// Verification is a hash of the rendered ignition to ensure that it has
	// not been changed when loaded as a remote file by the bootstrap ignition.
	// See https://coreos.com/ignition/docs/latest/configuration-v2_2.html
	Verification IgnitionStatusVerification `json:"verification" yaml:"verification"`
}

type IgnitionStatusVerification struct {
	// The content of the full rendered ignition hashed by the corresponding algorithm.
	Hash string `json:"hash" yaml:"hash"`
	// The algorithm used for hashing. Must be sha512 for now.
	Algorithm string `json:"algorithm" yaml:"algorithm"`
}

type IgnitionStatusSecret struct {
	// Name is the name of the secret containing the rendered ignition.
	Name string `json:"name" yaml:"name"`
	// Namespace is the namespace of the secret containing the rendered ignition.
	Namespace string `json:"namespace" yaml:"namespace"`
	// ResourceVersion is the Kubernetes resource version of the secret.
	// Used to detect if the secret has changed, e.g. 12345.
	ResourceVersion string `json:"resourceVersion" yaml:"resourceVersion"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type IgnitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Ignition `json:"items"`
}
