package v1alpha1

import "time"

const (
	ClusterVersionLimit = 5
)

const (
	StatusClusterStatusFalse = "False"
	StatusClusterStatusTrue  = "True"
)

const (
	StatusClusterTypeCreated  = "Created"
	StatusClusterTypeCreating = "Creating"
)

const (
	StatusClusterTypeUpdated  = "Updated"
	StatusClusterTypeUpdating = "Updating"
)

type StatusCluster struct {
	// Conditions is a list of status information expressing the current
	// conditional state of a guest cluster. This may reflect the status of the
	// guest cluster being updating or being up to date.
	Conditions []StatusClusterCondition `json:"conditions" yaml:"conditions"`
	Network    StatusClusterNetwork     `json:"network" yaml:"network"`
	// Nodes is a list of guest cluster node information reflecting the current
	// state of the guest cluster nodes.
	Nodes []StatusClusterNode `json:"nodes" yaml:"nodes"`
	// Versions is a list that acts like a historical track record of versions a
	// guest cluster went through. A version is only added to the list as soon as
	// the guest cluster successfully migrated to the version added here.
	Versions []StatusClusterVersion `json:"versions" yaml:"versions"`
}

// StatusClusterCondition expresses the conditions in which a guest cluster may
// is.
type StatusClusterCondition struct {
	// Status may be True, False or Unknown.
	Status string `json:"status" yaml:"status"`
	// Type may be Creating, Created, Scaling, Scaled, Draining, Drained,
	// Deleting, Deleted.
	Type string `json:"type" yaml:"type"`
}

// StatusClusterNetwork expresses the network segment that is allocated for a
// guest cluster.
type StatusClusterNetwork struct {
	CIDR string `json:"cidr" yaml:"cidr"`
}

// StatusClusterNode holds information about a guest cluster node.
type StatusClusterNode struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

// StatusClusterVersion expresses the versions in which a guest cluster was and
// may still be.
type StatusClusterVersion struct {
	// Date is the time of the given guest cluster version being updated.
	Date time.Time `json:"date" yaml:"date"`
	// Semver is some semver version, e.g. 1.0.0.
	Semver string `json:"semver" yaml:"semver"`
}

// DeepCopyInto implements the deep copy magic the k8s codegen is not able to
// generate out of the box.
func (in *StatusClusterVersion) DeepCopyInto(out *StatusClusterVersion) {
	*out = *in
}
