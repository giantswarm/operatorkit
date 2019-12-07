package v1alpha1

func (a AzureConfig) AvailabilityZones() []int {
	return a.Spec.Azure.AvailabilityZones
}

func (a AzureConfig) ClusterStatus() StatusCluster {
	return a.Status.Cluster
}
