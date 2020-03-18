package v1alpha2

import (
	"github.com/ghodss/yaml"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	crDocsAnnotation            = "giantswarm.io/docs"
	kindAWSCluster              = "AWSCluster"
	awsClusterDocumentationLink = "https://pkg.go.dev/github.com/giantswarm/apiextensions/pkg/apis/infrastructure/v1alpha2?tab=doc#AWSCluster"
)

const awsClusterCRDYAML = `
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: awsclusters.infrastructure.giantswarm.io
spec:
  conversion:
    strategy: None
  group: infrastructure.giantswarm.io
  names:
    kind: AWSCluster
    listKind: AWSClusterList
    plural: awsclusters
    singular: awscluster
  preserveUnknownFields: true
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: false
    storage: false
    schema:
      openAPIV3Schema:
        type: object
        properties: {}
  - name: v1alpha2
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        description: |
          Defines a tenant cluster in a Giant Swarm AWS installation.
          Introduced with release v10.x.x, reconciled by aws-operator.
        type: object
        properties:
          spec:
            type: object
            properties:
              cluster:
                description: |
                  Provides cluster specification details.
                type: object
                properties:
                  description:
                    description: |
                      User-friendly description that should explain the purpose of the
                      cluster.
                    maxLength: 100
                    type: string
                  dns:
                    description: |
                      DNS configuration details.
                    type: object
                    properties:
                      domain:
                        description: |
                          Base domain for several endpoints of this cluster.
                        type: string
                  oidc:
                    description: |
                      Configuration for OpenID Connect (OIDC) authentication.
                    type: object
              provider:
                description: |
                  AWS-specific configuration details.
                type: object
                properties:
                  master:
                    description: |
                      Master node configuration details.
                    type: object
                    properties:
                      availabilityZone:
                        description: |
                          Name of the AWS Availability Zone to place the master node in.
                        type: string
                      instanceType:
                        description: |
                          EC2 instance type to use for the master node.
                        type: string
                  region:
                    description: |
                      AWS region the cluster is to be running in.
                    type: string
    subresources:
      status: {}
`

var awsClusterCRD *apiextensionsv1beta1.CustomResourceDefinition

func init() {
	err := yaml.Unmarshal([]byte(awsClusterCRDYAML), &awsClusterCRD)
	if err != nil {
		panic(err)
	}
}

func NewAWSClusterCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return awsClusterCRD.DeepCopy()
}

func NewAWSClusterTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: SchemeGroupVersion.String(),
		Kind:       kindAWSCluster,
	}
}

// NewAWSClusterCR returns an AWSCluster Custom Resource.
func NewAWSClusterCR() *AWSCluster {
	return &AWSCluster{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				crDocsAnnotation: awsClusterDocumentationLink,
			},
		},
		TypeMeta: NewAWSClusterTypeMeta(),
	}
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSCluster is the infrastructure provider referenced in upstream CAPI Cluster
// CRs.
//
//     apiVersion: infrastructure.giantswarm.io/v1alpha2
//     kind: AWSCluster
//     metadata:
//       labels:
//         aws-operator.giantswarm.io/version: 6.2.0
//         cluster-operator.giantswarm.io/version: 0.17.0
//         giantswarm.io/cluster: "8y5kc"
//         giantswarm.io/organization: "giantswarm"
//         release.giantswarm.io/version: 7.3.1
//       name: 8y5kc
//     spec:
//       cluster:
//         description: my fancy cluster
//         dns:
//           domain: gauss.eu-central-1.aws.gigantic.io
//         oidc:
//           claims:
//             username: email
//             groups: groups
//           clientID: foobar-dex-client
//           issuerURL: https://dex.gatekeeper.eu-central-1.aws.example.com
//       provider:
//         credentialSecret:
//           name: credential-default
//           namespace: giantswarm
//         master:
//           availabilityZone: eu-central-1a
//           instanceType: m4.large
//         region: eu-central-1
//     status:
//       cluster:
//         conditions:
//         - lastTransitionTime: "2019-03-25T17:10:09.333633991Z"
//           type: Created
//         id: 8y5kc
//         versions:
//         - lastTransitionTime: "2019-03-25T17:10:09.995948706Z"
//           version: 4.9.0
//       provider:
//         network:
//           cidr: 10.1.6.0/24
//
type AWSCluster struct {
	metav1.TypeMeta `json:",inline"`
	// metav1.ObjectMeta is standard Kubernetes resource metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AWSClusterSpec   `json:"spec" yaml:"spec"`
	Status            AWSClusterStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// AWSClusterSpec is the spec part for the AWSCluster resource.
type AWSClusterSpec struct {
	// Cluster provides cluster specification details.
	Cluster AWSClusterSpecCluster `json:"cluster" yaml:"cluster"`
	// Provider holds provider-specific configuration details.
	Provider AWSClusterSpecProvider `json:"provider" yaml:"provider"`
}

// AWSClusterSpecCluster provides cluster specification details.
type AWSClusterSpecCluster struct {
	// Description is a user-friendly description that should explain the purpose of the
	// cluster to humans.
	Description string `json:"description" yaml:"description"`
	// DNS holds DNS configuration details.
	DNS AWSClusterSpecClusterDNS `json:"dns" yaml:"dns"`
	// OIDC holds configuration for OpenID Connect (OIDC) authentication.
	OIDC AWSClusterSpecClusterOIDC `json:"oidc" yaml:"oidc"`
}

// AWSClusterSpecClusterDNS holds DNS configuration details.
type AWSClusterSpecClusterDNS struct {
	Domain string `json:"domain" yaml:"domain"`
}

// AWSClusterSpecClusterOIDC holds configuration for OpenID Connect (OIDC) authentication.
type AWSClusterSpecClusterOIDC struct {
	Claims    AWSClusterSpecClusterOIDCClaims `json:"claims" yaml:"claims"`
	ClientID  string                          `json:"clientID" yaml:"clientID"`
	IssuerURL string                          `json:"issuerURL" yaml:"issuerURL"`
}

// AWSClusterSpecClusterOIDCClaims defines OIDC claims.
type AWSClusterSpecClusterOIDCClaims struct {
	Username string `json:"username" yaml:"username"`
	Groups   string `json:"groups" yaml:"groups"`
}

// AWSClusterSpecProvider holds some AWS details.
type AWSClusterSpecProvider struct {
	// CredentialSecret specifies the location of the secret providing the ARN of AWS IAM identity
	// to use with this cluster.
	CredentialSecret AWSClusterSpecProviderCredentialSecret `json:"credentialSecret" yaml:"credentialSecret"`
	// Master holds master node configuration details.
	Master AWSClusterSpecProviderMaster `json:"master" yaml:"master"`
	// Region is the AWS region the cluster is to be running in.
	Region string `json:"region" yaml:"region"`
}

// AWSClusterSpecProviderCredentialSecret details how to chose the AWS IAM identity ARN
// to use with this cluster.
type AWSClusterSpecProviderCredentialSecret struct {
	// Name is the name of the provider credential resoure.
	Name string `json:"name" yaml:"name"`
	// Namespace is the kubernetes namespace that holds the provider credential.
	Namespace string `json:"namespace" yaml:"namespace"`
}

// AWSClusterSpecProviderMaster holds master node configuration details.
type AWSClusterSpecProviderMaster struct {
	// AvailabilityZone is the AWS availability zone to place the master node in.
	AvailabilityZone string `json:"availabilityZone" yaml:"availabilityZone"`
	// InstanceType specifies the AWS EC2 instance type to use for the master node.
	InstanceType string `json:"instanceType" yaml:"instanceType"`
}

// AWSClusterStatus holds status information about the cluster, populated once the
// cluster is in creation or created.
type AWSClusterStatus struct {
	// Cluster provides cluster-specific status details, including conditions and versions.
	Cluster CommonClusterStatus `json:"cluster,omitempty" yaml:"cluster,omitempty"`
	// Provider provides provider-specific status details.
	Provider AWSClusterStatusProvider `json:"provider,omitempty" yaml:"provider,omitempty"`
}

// AWSClusterStatusProvider holds provider-specific status details.
type AWSClusterStatusProvider struct {
	// Network provides network-specific configuration details
	Network AWSClusterStatusProviderNetwork `json:"network" yaml:"network"`
}

// AWSClusterStatusProviderNetwork holds network details.
type AWSClusterStatusProviderNetwork struct {
	// IPv4 address block used by the tenant cluster, in CIDR notation.
	CIDR string `json:"cidr" yaml:"cidr"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSClusterList is the type returned when listing AWSCLuster resources.
type AWSClusterList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []AWSCluster `json:"items" yaml:"items"`
}
