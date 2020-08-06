package recorder

import (
	"context"

	pkgruntime "k8s.io/apimachinery/pkg/runtime"
)

type Interface interface {
	// Emit is used to create Kubernetes events.
	Emit(ctx context.Context, obj pkgruntime.Object, err error)
}
