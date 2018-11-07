package collector

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// Watcher provides Watch method compatible with Kubernetes clients.
type Watcher interface {
	Watch(metav1.ListOptions) (watch.Interface, error)
}
