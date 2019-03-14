package secret

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

type StateGetter interface {
	// GetCurrentState returns a current state of the system for the given
	// carnation of the observed Kubernetes object.
	GetCurrentState(ctx context.Context, obj interface{}) ([]*v1.Secret, error)
	// GetDesiredState returns a desired state of the system for the given
	// carnation of the observed Kubernetes object.
	//
	// NOTE: This state may be different if the observed object is
	// created/updated or deleted. Deletion timestamp can be checked to
	// figure if the object is scheduled for deletion.
	GetDesiredState(ctx context.Context, obj interface{}) ([]*v1.Secret, error)
}
