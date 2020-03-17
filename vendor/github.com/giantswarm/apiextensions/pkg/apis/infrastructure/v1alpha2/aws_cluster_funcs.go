package v1alpha2

func (c *AWSCluster) GetCommonClusterStatus() CommonClusterStatus {
	return c.Status.Cluster
}

func (c *AWSCluster) SetCommonClusterStatus(s CommonClusterStatus) {
	c.Status.Cluster = s
}
