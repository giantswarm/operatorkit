package app

import (
	"context"

	"github.com/giantswarm/apiextensions/pkg/apis/application/v1alpha1"
)

type StateGetter interface {
	// GetCurrentState returns a current state of the system for the given
	// carnation of the observed Kubernetes object.
	GetCurrentState(ctx context.Context, obj interface{}) ([]*v1alpha1.App, error)
	// GetDesiredState returns a desired state of the system for the given
	// carnation of the observed Kubernetes object.
	//
	// NOTE: This state may be different if the observed object is
	// created/updated or deleted. Deletion timestamp can be checked to
	// figure if the object is scheduled for deletion.
	GetDesiredState(ctx context.Context, obj interface{}) ([]*v1alpha1.App, error)
}
