package v1alpha1

func (k KVMConfig) ClusterStatus() StatusCluster {
	return k.Status.Cluster
}
