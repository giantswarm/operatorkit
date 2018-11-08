package v1alpha1

func (a AzureConfig) ClusterStatus() StatusCluster {
	return a.Status.Cluster
}
