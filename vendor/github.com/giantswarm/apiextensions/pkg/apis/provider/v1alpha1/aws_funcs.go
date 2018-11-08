package v1alpha1

func (a AWSConfig) AvailabilityZones() int {
	return a.Spec.AWS.AvailabilityZones
}

func (a AWSConfig) ClusterStatus() StatusCluster {
	return a.Status.Cluster
}
