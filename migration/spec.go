package migration

import (
	"context"

	"k8s.io/apimachinery/pkg/watch"
)

// Migrator implements a migration process from one object vesion to another.
type Migrator interface {
	// Create receives the transformed object returned by Transform. Create is
	// supposed to create the new object version in some API, e.g. creating a new
	// custom object in the Kubernetes API.
	Create(obj interface{}) error
	// Delete receives the old object returned by List. Delete is supposed to
	// delete the old object version from some API, if desired, e.g. deleting the
	// old custom object from the Kubernetes API.
	Delete(obj interface{}) error
	// Init sets up the migration to e.g. create a new TPR/CRD version in the
	// Kubernetes API.
	Init() error
	// List returns the already existing old objects supposed to be migrated, e.g.
	// list all custom objects of the Kubernetes API supposed to be migrated.
	List(ctx context.Context) (chan watch.Event, error)
	// Transform receives the old object as returned by List. Transform mutates
	// the old object into a new one and essentially migrates the object. The
	// result returned by Transform is given to Create.
	Transform(obj interface{}) (interface{}, error)
}
