package v1alpha1

import (
	"sort"
	"time"
)

func NewStatusClusterNode(name, version string, labels map[string]string) StatusClusterNode {
	return StatusClusterNode{
		Labels:             labels,
		LastTransitionTime: DeepCopyTime{time.Now()},
		Name:               name,
		Version:            version,
	}
}

func (s StatusCluster) GetCreatedCondition() StatusClusterCondition {
	return getCondition(s.Conditions, StatusClusterStatusTrue, StatusClusterTypeCreated)
}

func (s StatusCluster) GetCreatingCondition() StatusClusterCondition {
	return getCondition(s.Conditions, StatusClusterStatusTrue, StatusClusterTypeCreating)
}

func (s StatusCluster) GetDeletedCondition() StatusClusterCondition {
	return getCondition(s.Conditions, StatusClusterStatusTrue, StatusClusterTypeDeleted)
}

func (s StatusCluster) GetDeletingCondition() StatusClusterCondition {
	return getCondition(s.Conditions, StatusClusterStatusTrue, StatusClusterTypeDeleting)
}

func (s StatusCluster) GetUpdatedCondition() StatusClusterCondition {
	return getCondition(s.Conditions, StatusClusterStatusTrue, StatusClusterTypeUpdated)
}

func (s StatusCluster) GetUpdatingCondition() StatusClusterCondition {
	return getCondition(s.Conditions, StatusClusterStatusTrue, StatusClusterTypeUpdating)
}

func (s StatusCluster) HasCreatedCondition() bool {
	return hasCondition(s.Conditions, StatusClusterStatusTrue, StatusClusterTypeCreated)
}

func (s StatusCluster) HasCreatingCondition() bool {
	return hasCondition(s.Conditions, StatusClusterStatusTrue, StatusClusterTypeCreating)
}

func (s StatusCluster) HasDeletedCondition() bool {
	return hasCondition(s.Conditions, StatusClusterStatusTrue, StatusClusterTypeDeleted)
}

func (s StatusCluster) HasDeletingCondition() bool {
	return hasCondition(s.Conditions, StatusClusterStatusTrue, StatusClusterTypeDeleting)
}

func (s StatusCluster) HasUpdatedCondition() bool {
	return hasCondition(s.Conditions, StatusClusterStatusTrue, StatusClusterTypeUpdated)
}

func (s StatusCluster) HasUpdatingCondition() bool {
	return hasCondition(s.Conditions, StatusClusterStatusTrue, StatusClusterTypeUpdating)
}

func (s StatusCluster) HasVersion(semver string) bool {
	return hasVersion(s.Versions, semver)
}

func (s StatusCluster) LatestVersion() string {
	if len(s.Versions) == 0 {
		return ""
	}

	latest := s.Versions[0]

	for _, v := range s.Versions {
		if latest.LastTransitionTime.Time.Before(v.LastTransitionTime.Time) || latest.Date.Before(v.Date) {
			latest = v
		}
	}

	return latest.Semver
}

func (s StatusCluster) WithCreatedCondition() []StatusClusterCondition {
	return withCondition(s.Conditions, StatusClusterTypeCreating, StatusClusterTypeCreated, StatusClusterStatusTrue, time.Now())
}

func (s StatusCluster) WithCreatingCondition() []StatusClusterCondition {
	return withCondition(s.Conditions, StatusClusterTypeCreated, StatusClusterTypeCreating, StatusClusterStatusTrue, time.Now())
}

func (s StatusCluster) WithDeletedCondition() []StatusClusterCondition {
	return withCondition(s.Conditions, StatusClusterTypeDeleting, StatusClusterTypeDeleted, StatusClusterStatusTrue, time.Now())
}

func (s StatusCluster) WithDeletingCondition() []StatusClusterCondition {
	return withCondition(s.Conditions, StatusClusterTypeDeleted, StatusClusterTypeDeleting, StatusClusterStatusTrue, time.Now())
}

func (s StatusCluster) WithNewVersion(version string) []StatusClusterVersion {
	newVersion := StatusClusterVersion{
		LastTransitionTime: DeepCopyTime{time.Now()},
		Semver:             version,
	}

	return withVersion(s.Versions, newVersion, ClusterVersionLimit)
}

func (s StatusCluster) WithUpdatedCondition() []StatusClusterCondition {
	return withCondition(s.Conditions, StatusClusterTypeUpdating, StatusClusterTypeUpdated, StatusClusterStatusTrue, time.Now())
}

func (s StatusCluster) WithUpdatingCondition() []StatusClusterCondition {
	return withCondition(s.Conditions, StatusClusterTypeUpdated, StatusClusterTypeUpdating, StatusClusterStatusTrue, time.Now())
}

func getCondition(conditions []StatusClusterCondition, s string, t string) StatusClusterCondition {
	for _, c := range conditions {
		if c.Status == s && c.Type == t {
			return c
		}
	}

	return StatusClusterCondition{}
}

func hasCondition(conditions []StatusClusterCondition, s string, t string) bool {
	for _, c := range conditions {
		if c.Status == s && c.Type == t {
			return true
		}
	}

	return false
}

func hasVersion(versions []StatusClusterVersion, search string) bool {
	for _, v := range versions {
		if v.Semver == search {
			return true
		}
	}

	return false
}

func withCondition(conditions []StatusClusterCondition, search string, replace string, status string, t time.Time) []StatusClusterCondition {
	newConditions := []StatusClusterCondition{
		{
			LastTransitionTime: DeepCopyTime{t},
			Status:             status,
			Type:               replace,
		},
	}

	for _, c := range conditions {
		if c.Type == search {
			continue
		}

		newConditions = append(newConditions, c)
	}

	return newConditions
}

// withVersion computes a list of version history using the given list and new
// version structure to append. withVersion also limits total amount of elements
// in the list by cutting off the tail with respect to the limit parameter.
func withVersion(versions []StatusClusterVersion, version StatusClusterVersion, limit int) []StatusClusterVersion {
	if hasVersion(versions, version.Semver) {
		return versions
	}

	var newVersions []StatusClusterVersion

	start := 0
	if len(versions) >= limit {
		start = len(versions) - limit + 1
	}

	sort.Sort(sortClusterStatusVersionsByDate(versions))

	for i := start; i < len(versions); i++ {
		newVersions = append(newVersions, versions[i])
	}

	newVersions = append(newVersions, version)

	return newVersions
}
