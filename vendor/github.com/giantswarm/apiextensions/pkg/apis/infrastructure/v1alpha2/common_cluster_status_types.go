package v1alpha2

const (
	// ClusterConditionLimit is the maximum amount of conditions tracked in the
	// condition list of a tenant cluster's status. The limit here is applied to
	// equal condition pairs. For instance a cluster having transitioned through 6
	// cluster upgrades throughout its lifetime will only track 5 Updating/Updated
	// condition pairs in its condition list.
	//
	//     conditions:
	//     - lastTransitionTime: "2019-08-23T13:15:19.830177296Z"
	//       condition: Updated
	//     - lastTransitionTime: "2019-08-23T12:12:25.942680489Z"
	//       condition: Updating
	//     - lastTransitionTime: "2019-08-15T14:27:12.813903533Z"
	//       condition: Updated
	//     - lastTransitionTime: "2019-08-15T13:20:16.955248597Z"
	//       condition: Updating
	//     - lastTransitionTime: "2019-07-23T09:31:28.761118959Z"
	//       condition: Updated
	//     - lastTransitionTime: "2019-07-23T08:15:07.523067044Z"
	//       condition: Updating
	//     - lastTransitionTime: "2019-06-17T18:20:30.29872263Z"
	//       condition: Updated
	//     - lastTransitionTime: "2019-06-17T17:14:12.707323902Z"
	//       condition: Updating
	//     - lastTransitionTime: "2019-06-04T13:14:03.523010234Z"
	//       condition: Updated
	//     - lastTransitionTime: "2019-06-04T12:18:09.334829389Z"
	//       condition: Updating
	//     - lastTransitionTime: "2019-05-17T11:25:37.495980406Z"
	//       condition: Created
	//     - lastTransitionTime: "2019-05-17T10:16:25.736159078Z"
	//       condition: Creating
	//
	ClusterConditionLimit = 5
	// ClusterVersionLimit is the maximum amount of versions tracked in the
	// version list of a tenant cluster's status. The limit here is applied to the
	// total amount of the list's number of entries. For instance a cluster having
	// transitioned through 6 cluster upgrades throughout its lifetime will only
	// track 5 versions in its version list.
	//
	//     versions:
	//     - lastTransitionTime: "2019-02-14T11:18:25.212331926Z"
	//       version: 4.6.0
	//     - lastTransitionTime: "2018-12-05T16:57:58.21652461Z"
	//       version: 4.4.1
	//     - lastTransitionTime: "2018-12-05T15:42:22.443182449Z"
	//       version: 4.2.1
	//     - lastTransitionTime: "2018-10-29T03:31:08.874296621Z"
	//       version: 4.2.0
	//     - lastTransitionTime: "2018-10-29T02:09:20.393986006Z"
	//       version: 3.3.3
	//
	ClusterVersionLimit = 15
)

const (
	ClusterStatusConditionCreated  = "Created"
	ClusterStatusConditionCreating = "Creating"
)

const (
	ClusterStatusConditionDeleted  = "Deleted"
	ClusterStatusConditionDeleting = "Deleting"
)

const (
	ClusterStatusConditionUpdated  = "Updated"
	ClusterStatusConditionUpdating = "Updating"
)

var (
	conditionPairs = [][]string{
		[]string{
			ClusterStatusConditionCreated,
			ClusterStatusConditionCreating,
		},
		[]string{
			ClusterStatusConditionDeleted,
			ClusterStatusConditionDeleting,
		},
		[]string{
			ClusterStatusConditionUpdated,
			ClusterStatusConditionUpdating,
		},
	}
)

// CommonClusterStatus is shared type to contain provider independent cluster
// status information.
type CommonClusterStatus struct {
	// One or several conditions that are currently applicable to the cluster.
	Conditions []CommonClusterStatusCondition `json:"conditions,omitempty" yaml:"conditions,omitempty"`
	// Identifier of the cluster.
	ID string `json:"id" yaml:"id"`
	// Release versions the cluster used so far.
	Versions []CommonClusterStatusVersion `json:"versions,omitempty" yaml:"versions,omitempty"`
}

// CommonClusterStatusCondition explains the current condition(s) of the cluster.
type CommonClusterStatusCondition struct {
	// Time the condition occurred.
	LastTransitionTime DeepCopyTime `json:"lastTransitionTime" yaml:"lastTransitionTime"`
	// Condition string, e. g. `Creating`, `Created`, `Upgraded`
	Condition string `json:"condition" yaml:"condition"`
}

// CommonClusterStatusVersion informs which aws-operator version was/responsible for this cluster.
type CommonClusterStatusVersion struct {
	// Time the cluster assumed the given version.
	LastTransitionTime DeepCopyTime `json:"lastTransitionTime" yaml:"lastTransitionTime"`
	// The aws-operator version responsible for handling the cluster.
	Version string `json:"version" yaml:"version"`
}
