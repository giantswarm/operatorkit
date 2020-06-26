package recorder

import (
	"context"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	pkgruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

type Recorder struct {
	record.EventRecorder
}

func (r *Recorder) Emit(ctx context.Context, obj pkgruntime.Object, err error) {
	if merr, ok := microerror.Cause(err).(*microerror.Error); ok {
		if merr.Kind != "" && merr.Desc != "" {
			r.Event(obj, corev1.EventTypeWarning, merr.Kind, merr.Desc)
		}
	}
}
