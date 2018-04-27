package informer

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type Interface interface {
	Boot(ctx context.Context) error
	ResyncPeriod() time.Duration
	Watch(ctx context.Context) (chan watch.Event, chan watch.Event, chan error)
}

// Watcher provides Watch method compatible with Kubernetes clients.
type Watcher interface {
	Watch(metav1.ListOptions) (watch.Interface, error)
}
