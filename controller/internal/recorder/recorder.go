package recorder

import (
	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	clientv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
)

type Config struct {
	Component string
	K8sclient k8sclient.Interface
}

// New creates an event recorder to send custom events to Kubernetes to be recorded for targeted Kubernetes objects
func New(c Config) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	if _, isfake := c.K8sclient.K8sClient().(*fake.Clientset); !isfake {
		eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(c.K8sclient.K8sClient().CoreV1().RESTClient()).Events("")})
	}
	return eventBroadcaster.NewRecorder(c.K8sclient.Scheme(), clientv1.EventSource{Component: c.Component})
}
