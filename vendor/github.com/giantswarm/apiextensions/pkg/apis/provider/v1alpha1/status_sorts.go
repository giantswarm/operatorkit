package v1alpha1

type sortClusterStatusVersionsByDate []StatusClusterVersion

func (s sortClusterStatusVersionsByDate) Len() int      { return len(s) }
func (s sortClusterStatusVersionsByDate) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortClusterStatusVersionsByDate) Less(i, j int) bool {
	return s[i].Date.UnixNano() < s[j].Date.UnixNano()
}
