package recorder

import (
	"context"
	"errors"

	"github.com/giantswarm/k8sclient/v5/pkg/k8sclient"
	"github.com/giantswarm/k8sclient/v5/pkg/k8sclienttest"
	"github.com/giantswarm/microerror"
	clientv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	pkgruntime "k8s.io/apimachinery/pkg/runtime"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
)

type Config struct {
	Component string
	K8sClient k8sclient.Interface
}

type Recorder struct {
	record.EventRecorder
}

// New creates an event recorder to send custom events to Kubernetes to be recorded for targeted Kubernetes objects.
func New(c Config) Interface {
	eventBroadcaster := record.NewBroadcaster()
	_, isfake := c.K8sClient.(*k8sclienttest.Clients)
	if !isfake {
		eventBroadcaster.StartRecordingToSink(
			&typedcorev1.EventSinkImpl{
				Interface: c.K8sClient.K8sClient().CoreV1().Events(""),
			},
		)
	}
	return &Recorder{
		eventBroadcaster.NewRecorder(c.K8sClient.Scheme(), clientv1.EventSource{Component: c.Component}),
	}
}

func (r *Recorder) Emit(ctx context.Context, obj pkgruntime.Object, err error) {
	var merr *microerror.Error
	if errors.As(err, &merr) {
		if merr.Kind != "" && merr.Desc != "" {
			r.Event(obj, corev1.EventTypeWarning, merr.Kind, merr.Desc)
		}
	}
}
