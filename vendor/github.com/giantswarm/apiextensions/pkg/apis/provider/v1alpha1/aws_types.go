package v1alpha1

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewAWSConfigCRD returns a new custom resource definition for AWSConfig. This
// might look something like the following.
//
//     apiVersion: apiextensions.k8s.io/v1beta1
//     kind: CustomResourceDefinition
//     metadata:
//       name: awsconfigs.provider.giantswarm.io
//     spec:
//       group: provider.giantswarm.io
//       scope: Namespaced
//       version: v1alpha1
//       names:
//         kind: AWSConfig
//         plural: awsconfigs
//         singular: awsconfig
//       subresources:
//         status: {}
//
func NewAWSConfigCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apiextensionsv1beta1.SchemeGroupVersion.String(),
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "awsconfigs.provider.giantswarm.io",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "provider.giantswarm.io",
			Scope:   "Namespaced",
			Version: "v1alpha1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:     "AWSConfig",
				Plural:   "awsconfigs",
				Singular: "awsconfig",
			},
			Subresources: &apiextensionsv1beta1.CustomResourceSubresources{
				Status: &apiextensionsv1beta1.CustomResourceSubresourceStatus{},
			},
		},
	}
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AWSConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              AWSConfigSpec   `json:"spec"`
	Status            AWSConfigStatus `json:"status" yaml:"status"`
}

type AWSConfigSpec struct {
	Cluster       Cluster                    `json:"cluster" yaml:"cluster"`
	AWS           AWSConfigSpecAWS           `json:"aws" yaml:"aws"`
	VersionBundle AWSConfigSpecVersionBundle `json:"versionBundle" yaml:"versionBundle"`
}

type AWSConfigSpecAWS struct {
	API AWSConfigSpecAWSAPI `json:"api" yaml:"api"`
	// TODO remove the deprecated AZ field due to AvailabilityZones.
	//
	//     https://github.com/giantswarm/giantswarm/issues/4507
	//
	AZ string `json:"az" yaml:"az"`
	// AvailabilityZones is the number of AWS availability zones used to spread
	// the tenant cluster's worker nodes across. There are limitations on
	// availability zone settings due to binary IP range splitting and provider
	// specific region capabilities. When for instance choosing 3 availability
	// zones, the configured IP range will be split into 4 ranges and thus one of
	// it will not be able to be utilized. Such limitations have to be considered
	// when designing the network topology and configuring tenant cluster HA via
	// AvailabilityZones.
	//
	// The selection and usage of the actual availability zones for the created
	// tenant cluster is randomized. In case there are 4 availability zones
	// provided in the used region and the user selects 2 availability zones, the
	// actually used availability zones in which tenant cluster workload is put
	// into will tend to be different across tenant cluster creations. This is
	// done in order to provide more HA during single availability zone failures.
	// In case a specific availability zone fails, not all tenant clusters will be
	// affected due to the described selection process.
	AvailabilityZones int                  `json:"availabilityZones" yaml:"availabilityZones"`
	CredentialSecret  CredentialSecret     `json:"credentialSecret" yaml:"credentialSecret"`
	Etcd              AWSConfigSpecAWSEtcd `json:"etcd" yaml:"etcd"`

	// HostedZones is AWS hosted zones names in the host cluster account.
	// For each zone there will be "CLUSTER_ID.k8s" NS record created in
	// the host cluster account. Then for each created NS record there will
	// be a zone created in the guest account. After that component
	// specific records under those zones:
	//	- api.CLUSTER_ID.k8s.{{ .Spec.AWS.HostedZones.API.Name }}
	//	- etcd.CLUSTER_ID.k8s.{{ .Spec.AWS.HostedZones.Etcd.Name }}
	//	- ingress.CLUSTER_ID.k8s.{{ .Spec.AWS.HostedZones.Ingress.Name }}
	//	- *.CLUSTER_ID.k8s.{{ .Spec.AWS.HostedZones.Ingress.Name }}
	HostedZones AWSConfigSpecAWSHostedZones `json:"hostedZones" yaml:"hostedZones"`

	Ingress AWSConfigSpecAWSIngress `json:"ingress" yaml:"ingress"`
	Masters []AWSConfigSpecAWSNode  `json:"masters" yaml:"masters"`
	Region  string                  `json:"region" yaml:"region"`
	VPC     AWSConfigSpecAWSVPC     `json:"vpc" yaml:"vpc"`
	Workers []AWSConfigSpecAWSNode  `json:"workers" yaml:"workers"`
}

// AWSConfigSpecAWSAPI deprecated since aws-operator v12 resources.
type AWSConfigSpecAWSAPI struct {
	HostedZones string                 `json:"hostedZones" yaml:"hostedZones"`
	ELB         AWSConfigSpecAWSAPIELB `json:"elb" yaml:"elb"`
}

// AWSConfigSpecAWSAPIELB deprecated since aws-operator v12 resources.
type AWSConfigSpecAWSAPIELB struct {
	IdleTimeoutSeconds int `json:"idleTimeoutSeconds" yaml:"idleTimeoutSeconds"`
}

// AWSConfigSpecAWSEtcd deprecated since aws-operator v12 resources.
type AWSConfigSpecAWSEtcd struct {
	HostedZones string                  `json:"hostedZones" yaml:"hostedZones"`
	ELB         AWSConfigSpecAWSEtcdELB `json:"elb" yaml:"elb"`
}

// AWSConfigSpecAWSEtcdELB deprecated since aws-operator v12 resources.
type AWSConfigSpecAWSEtcdELB struct {
	IdleTimeoutSeconds int `json:"idleTimeoutSeconds" yaml:"idleTimeoutSeconds"`
}

type AWSConfigSpecAWSHostedZones struct {
	API     AWSConfigSpecAWSHostedZonesZone `json:"api" yaml:"api"`
	Etcd    AWSConfigSpecAWSHostedZonesZone `json:"etcd" yaml:"etcd"`
	Ingress AWSConfigSpecAWSHostedZonesZone `json:"ingress" yaml:"ingress"`
}

type AWSConfigSpecAWSHostedZonesZone struct {
	Name string `json:"name" yaml:"name"`
}

// AWSConfigSpecAWSIngress deprecated since aws-operator v12 resources.
type AWSConfigSpecAWSIngress struct {
	HostedZones string                     `json:"hostedZones" yaml:"hostedZones"`
	ELB         AWSConfigSpecAWSIngressELB `json:"elb" yaml:"elb"`
}

// AWSConfigSpecAWSIngressELB deprecated since aws-operator v12 resources.
type AWSConfigSpecAWSIngressELB struct {
	IdleTimeoutSeconds int `json:"idleTimeoutSeconds" yaml:"idleTimeoutSeconds"`
}

type AWSConfigSpecAWSNode struct {
	ImageID            string `json:"imageID" yaml:"imageID"`
	InstanceType       string `json:"instanceType" yaml:"instanceType"`
	DockerVolumeSizeGB int    `json:"dockerVolumeSizeGB" yaml:"dockerVolumeSizeGB"`
}

type AWSConfigSpecAWSVPC struct {
	CIDR              string   `json:"cidr" yaml:"cidr"`
	PrivateSubnetCIDR string   `json:"privateSubnetCidr" yaml:"privateSubnetCidr"`
	PublicSubnetCIDR  string   `json:"publicSubnetCidr" yaml:"publicSubnetCidr"`
	RouteTableNames   []string `json:"routeTableNames" yaml:"routeTableNames"`
	PeerID            string   `json:"peerId" yaml:"peerId"`
}

type AWSConfigSpecVersionBundle struct {
	Version string `json:"version" yaml:"version"`
}

type AWSConfigStatus struct {
	AWS     AWSConfigStatusAWS `json:"aws" yaml:"aws"`
	Cluster StatusCluster      `json:"cluster" yaml:"cluster"`
}

type AWSConfigStatusAWS struct {
	AvailabilityZones []AWSConfigStatusAWSAvailabilityZone `json:"availabilityZones" yaml:"availabilityZones"`
	AutoScalingGroup  AWSConfigStatusAWSAutoScalingGroup   `json:"autoScalingGroup" yaml:"autoScalingGroup"`
}

type AWSConfigStatusAWSAutoScalingGroup struct {
	Name string `json:"name"`
}

type AWSConfigStatusAWSAvailabilityZone struct {
	Name   string                                   `json:"name" yaml:"name"`
	Subnet AWSConfigStatusAWSAvailabilityZoneSubnet `json:"subnet" yaml:"subnet"`
}

type AWSConfigStatusAWSAvailabilityZoneSubnet struct {
	Private AWSConfigStatusAWSAvailabilityZoneSubnetPrivate `json:"private" yaml:"private"`
	Public  AWSConfigStatusAWSAvailabilityZoneSubnetPublic  `json:"public" yaml:"public"`
}

type AWSConfigStatusAWSAvailabilityZoneSubnetPrivate struct {
	CIDR string `json:"cidr" yaml:"cidr"`
}

type AWSConfigStatusAWSAvailabilityZoneSubnetPublic struct {
	CIDR string `json:"cidr" yaml:"cidr"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AWSConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AWSConfig `json:"items"`
}
