package v1alpha2

type sortClusterStatusVersionsByDate []CommonClusterStatusVersion

func (s sortClusterStatusVersionsByDate) Len() int      { return len(s) }
func (s sortClusterStatusVersionsByDate) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortClusterStatusVersionsByDate) Less(i, j int) bool {
	return s[i].LastTransitionTime.UnixNano() < s[j].LastTransitionTime.UnixNano()
}

type sortClusterStatusConditionsByDate []CommonClusterStatusCondition

func (s sortClusterStatusConditionsByDate) Len() int      { return len(s) }
func (s sortClusterStatusConditionsByDate) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortClusterStatusConditionsByDate) Less(i, j int) bool {
	return s[i].LastTransitionTime.UnixNano() < s[j].LastTransitionTime.UnixNano()
}
