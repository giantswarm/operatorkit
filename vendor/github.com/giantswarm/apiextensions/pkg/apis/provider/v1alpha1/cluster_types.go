package v1alpha1

import (
	"net"
)

type Cluster struct {
	Calico     ClusterCalico     `json:"calico" yaml:"calico"`
	Customer   ClusterCustomer   `json:"customer" yaml:"customer"`
	Docker     ClusterDocker     `json:"docker" yaml:"docker"`
	Etcd       ClusterEtcd       `json:"etcd" yaml:"etcd"`
	ID         string            `json:"id" yaml:"id"`
	Kubernetes ClusterKubernetes `json:"kubernetes" yaml:"kubernetes"`
	Masters    []ClusterNode     `json:"masters" yaml:"masters"`
	Scaling    ClusterScaling    `json:"scaling" yaml:"scaling"`

	// Version is DEPRECATED and should just be dropped.
	Version string `json:"version" yaml:"version"`

	Workers []ClusterNode `json:"workers" yaml:"workers"`
}

type ClusterCalico struct {
	CIDR   int    `json:"cidr" yaml:"cidr"`
	MTU    int    `json:"mtu" yaml:"mtu"`
	Subnet string `json:"subnet" yaml:"subnet"`
}

type ClusterCustomer struct {
	ID string `json:"id" yaml:"id"`
}

type ClusterDocker struct {
	Daemon ClusterDockerDaemon `json:"daemon" yaml:"daemon"`
}

type ClusterDockerDaemon struct {
	CIDR string `json:"cidr" yaml:"cidr"`
}

type ClusterEtcd struct {
	AltNames string `json:"altNames" yaml:"altNames"`
	Domain   string `json:"domain" yaml:"domain"`
	Port     int    `json:"port" yaml:"port"`
	Prefix   string `json:"prefix" yaml:"prefix"`
}

type ClusterKubernetes struct {
	API               ClusterKubernetesAPI               `json:"api" yaml:"api"`
	CloudProvider     string                             `json:"cloudProvider" yaml:"cloudProvider"`
	DNS               ClusterKubernetesDNS               `json:"dns" yaml:"dns"`
	Domain            string                             `json:"domain" yaml:"domain"`
	IngressController ClusterKubernetesIngressController `json:"ingressController" yaml:"ingressController"`
	Kubelet           ClusterKubernetesKubelet           `json:"kubelet" yaml:"kubelet"`
	NetworkSetup      ClusterKubernetesNetworkSetup      `json:"networkSetup" yaml:"networkSetup"`
	SSH               ClusterKubernetesSSH               `json:"ssh" yaml:"ssh"`
}

type ClusterKubernetesAPI struct {
	ClusterIPRange string `json:"clusterIPRange" yaml:"clusterIPRange"`
	Domain         string `json:"domain" yaml:"domain"`
	SecurePort     int    `json:"securePort" yaml:"securePort"`
}

type ClusterKubernetesDNS struct {
	IP net.IP `json:"ip" yaml:"ip"`
}

type ClusterKubernetesIngressController struct {
	Docker         ClusterKubernetesIngressControllerDocker `json:"docker" yaml:"docker"`
	Domain         string                                   `json:"domain" yaml:"domain"`
	WildcardDomain string                                   `json:"wildcardDomain" yaml:"wildcardDomain"`
	InsecurePort   int                                      `json:"insecurePort" yaml:"insecurePort"`
	SecurePort     int                                      `json:"securePort" yaml:"securePort"`
}

type ClusterKubernetesIngressControllerDocker struct {
	Image string `json:"image" yaml:"image"`
}

type ClusterKubernetesKubelet struct {
	AltNames string `json:"altNames" yaml:"altNames"`
	Domain   string `json:"domain" yaml:"domain"`
	Labels   string `json:"labels" yaml:"labels"`
	Port     int    `json:"port" yaml:"port"`
}

type ClusterKubernetesNetworkSetup struct {
	Docker ClusterKubernetesNetworkSetupDocker `json:"docker" yaml:"docker"`
}

type ClusterKubernetesNetworkSetupDocker struct {
	Image string `json:"image" yaml:"image"`
}

type ClusterKubernetesSSH struct {
	UserList []ClusterKubernetesSSHUser `json:"userList" yaml:"userList"`
}

type ClusterKubernetesSSHUser struct {
	Name      string `json:"name" yaml:"name"`
	PublicKey string `json:"publicKey" yaml:"publicKey"`
}

type ClusterNode struct {
	ID string `json:"id" yaml:"id"`
}

type ClusterScaling struct {
	// Max defines maximum number of worker nodes guest cluster is allowed to have.
	Max int `json:"max" yaml:"max"`
	// Min defines minimum number of worker nodes required to be present in guest cluster.
	Min int `json:"min" yaml:"min"`
}
