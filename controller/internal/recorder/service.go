package recorder

import (
	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	"github.com/giantswarm/k8sclient/v3/pkg/k8sclienttest"
	clientv1 "k8s.io/api/core/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
)

type Config struct {
	Component string
	K8sClient k8sclient.Interface
}

// New creates an event recorder to send custom events to Kubernetes to be recorded for targeted Kubernetes objects.
func New(c Config) Interface {
	eventBroadcaster := record.NewBroadcaster()
	if _, isfake := c.K8sClient.(*k8sclienttest.Clients); !isfake {
		eventBroadcaster.StartRecordingToSink(
			&typedcorev1.EventSinkImpl{
				Interface: c.K8sClient.K8sClient().CoreV1().Events("")})
	}
	return &Recorder{
		eventBroadcaster.NewRecorder(c.K8sClient.Scheme(), clientv1.EventSource{Component: c.Component}),
	}
}
