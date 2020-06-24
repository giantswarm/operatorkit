package recorder

import (
	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	clientv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
)

type Config struct {
	Component string
	K8sClient k8sclient.Interface
}

// New creates an event recorder to send custom events to Kubernetes to be recorded for targeted Kubernetes objects
func New(c Config) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	return eventBroadcaster.NewRecorder(c.K8sClient.Scheme(), clientv1.EventSource{Component: c.Component})
}
