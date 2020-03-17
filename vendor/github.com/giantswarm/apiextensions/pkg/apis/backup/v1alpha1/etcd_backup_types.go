package v1alpha1

import (
	"github.com/ghodss/yaml"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const etcdBackupCRDYAML = `
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: etcdbackups.backup.giantswarm.io
spec:
  group: backup.giantswarm.io
  scope: Cluster
  version: v1alpha1
  names:
    kind: ETCDBackup
    plural: etcdbackups
    singular: etcdbackup
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        spec:
          type: object
          properties:
            guestBackup:
              type: boolean
          required:
          - guestBackup
        status:
          type: object
          properties:
            status:
              enum:
              - Pending
              - V2BackupRunning
              - V2BackupCompleted
              - V3BackupRunning
              - V3BackupCompleted
              - Failed
              - Completed
            startedTimestamp:
              type: string
              format: date-time
            finishedTimestamp:
              type: string
              format: date-time
            instances:
              type: object
          required:
          - status
  additionalPrinterColumns:
  - name: guestBackup
    type: boolean
    description: Wether guest backups are backed up or not.
    JSONPath: .spec.guestBackup
  - name: Status
    type: string
    description: The status this backup is in
    JSONPath: .status.status
  - name: Started
    type: date
    description: The date the backup has been first attempted
    JSONPath: .status.startedTimestamp
  - name: Finished
    type: date
    description: The date the backup has been last attempted
    JSONPath: .status.finishedTimestamp
`

var etcdBackupCRD *apiextensionsv1beta1.CustomResourceDefinition

func init() {
	err := yaml.Unmarshal([]byte(etcdBackupCRDYAML), &etcdBackupCRD)
	if err != nil {
		panic(err)
	}
}

func NewETCDBackupCRD() *apiextensionsv1beta1.CustomResourceDefinition {
	return etcdBackupCRD.DeepCopy()
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ETCDBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ETCDBackupSpec   `json:"spec"`
	Status            ETCDBackupStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type ETCDBackupSpec struct {
	// GuestBackup is a boolean indicating if the tenant clusters have to be backupped
	GuestBackup bool `json:"guestBackup" yaml:"guestBackup"`
}

type ETCDBackupStatus struct {
	// map containing the state of the backup for all instances
	Instances map[string]ETCDInstanceBackupStatusIndex `json:"instances,omitempty" yaml:"instances"`
	// Status of the whole backup job (can be 'Pending', 'Running'. 'Completed', 'Failed')
	Status string `json:"status" yaml:"status"`
	// Timestamp when the first attempt was made
	StartedTimestamp DeepCopyTime `json:"startedTimestamp,omitempty" yaml:"startedTimestamp"`
	// Timestamp when the last (final) attempt was made (when the Phase became either 'Completed' or 'Failed'
	FinishedTimestamp DeepCopyTime `json:"finishedTimestamp,omitempty" yaml:"finishedTimestamp"`
}

type ETCDInstanceBackupStatusIndex struct {
	// Name of the tenant cluster or 'Control Plane'
	Name string `json:"name" yaml:"name"`
	// Status of the V2 backup for this instance
	V2 ETCDInstanceBackupStatus `json:"v2" yaml:"v2"`
	// Status of the V3 backup for this instance
	V3 ETCDInstanceBackupStatus `json:"v3" yaml:"v3"`
}

type ETCDInstanceBackupStatus struct {
	// Status of this isntance's backup job (can be 'Pending', 'Running'. 'Completed', 'Failed')
	Status string `json:"status" yaml:"status"`
	// Timestamp when the first attempt was made
	StartedTimestamp DeepCopyTime `json:"startedTimestamp,omitempty" yaml:"startedTimestamp"`
	// Timestamp when the last (final) attempt was made (when the Phase became either 'Completed' or 'Failed'
	FinishedTimestamp DeepCopyTime `json:"finishedTimestamp,omitempty" yaml:"finishedTimestamp"`
	// Latest backup error message
	LatestError string `json:"latestError,omitempty" yaml:"latestError,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ETCDBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ETCDBackup `json:"items"`
}
